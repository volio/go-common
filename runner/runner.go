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

type Service interface {
	Start() error
	Stop() error
}

type ServiceRunner interface {
	Wait()
}

func RunService(s Service) ServiceRunner {
	r := newServiceRunner(s)
	r.run()
	return r
}

func newServiceRunner(s Service) *serviceRunner {
	return &serviceRunner{
		signals: make(chan os.Signal, 1),
		service: s,
	}
}

type serviceRunner struct {
	service Service
	signals chan os.Signal
	wg      sync.WaitGroup
	stopped int32
}

func (r *serviceRunner) Wait() {
	r.wg.Wait()
}

func (r *serviceRunner) run() {
	r.wg.Add(1)
	go r.handleSignal()
	go r.handleStart()
}

func (r *serviceRunner) handleStart() {
	func() {
		defer util.Recovery()
		err := r.service.Start()
		if err != nil {
			log.L().With(zap.Error(err)).Error("start service failed")
		}
	}()
	if atomic.LoadInt32(&r.stopped) == 0 {
		r.wg.Done()
	}
}

func (r *serviceRunner) handleSignal() {
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

func (r *serviceRunner) signalHandler() {
	go func() {
		to := 10 * time.Second

		time.Sleep(to)
		log.L().Warn("service stop timeout")
		os.Exit(1)
	}()
	atomic.StoreInt32(&r.stopped, 1)
	err := r.service.Stop()
	if err != nil {
		log.L().With(zap.Error(err)).Error("stop service failed")
	}
}
