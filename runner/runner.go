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

func RunService(s Service) *ServiceRunner {
	r := newServiceRunner(s)
	r.Run()
	return r
}

func newServiceRunner(s Service) *ServiceRunner {
	return &ServiceRunner{
		signals: make(chan os.Signal, 1),
		service: s,
	}
}

type ServiceRunner struct {
	service Service
	signals chan os.Signal
	wg      sync.WaitGroup
	stopped int32
}

func (r *ServiceRunner) Wait() {
	r.wg.Wait()
}

func (r *ServiceRunner) Run() {
	r.wg.Add(1)
	go r.handleSignal()
	go r.handleStart()
}

func (r *ServiceRunner) handleStart() {
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

func (r *ServiceRunner) handleSignal() {
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

func (r *ServiceRunner) signalHandler() {
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
