package runner

import (
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/volio/go-common/log"
	"github.com/volio/go-common/util"
	"go.uber.org/zap"
)

type Server interface {
	Start() error
	Stop() error
}

type ServerRunner interface {
	Wait()
}

func RunServer(s Server) ServerRunner {
	r := newServiceRunner(s)
	r.run()
	return r
}

func newServiceRunner(s Server) *serverRunner {
	return &serverRunner{
		signals: make(chan os.Signal, 1),
		server:  s,
	}
}

type serverRunner struct {
	server  Server
	signals chan os.Signal
	wg      sync.WaitGroup
	stopped int32
}

func (r *serverRunner) Wait() {
	r.wg.Wait()
}

func (r *serverRunner) run() {
	r.wg.Add(1)
	go r.handleSignal()
	go r.handleStart()
}

func (r *serverRunner) handleStart() {
	func() {
		defer util.Recovery()
		err := r.server.Start()
		if err != nil {
			log.L().With(zap.Error(err)).Error("start server failed")
		}
	}()
	if atomic.LoadInt32(&r.stopped) == 0 {
		r.wg.Done()
	}
}

func (r *serverRunner) handleSignal() {
	signal.Notify(r.signals, syscall.SIGPIPE, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGABRT)
	for sig := range r.signals {
		switch sig {
		case syscall.SIGPIPE:
		case syscall.SIGINT:
			r.signalHandler()
			os.Exit(1)
		default:
			r.signalHandler()
			r.wg.Done()
		}
	}
}

func (r *serverRunner) signalHandler() {
	go func() {
		to := 10 * time.Second

		time.Sleep(to)
		log.L().Warn("server stop timeout")
		os.Exit(1)
	}()
	atomic.StoreInt32(&r.stopped, 1)
	err := r.server.Stop()
	if err != nil {
		log.L().With(zap.Error(err)).Error("stop server failed")
	}
}
