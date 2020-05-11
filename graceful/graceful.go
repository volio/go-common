package graceful

import (
	"time"

	"github.com/volio/go-common/log"
	"github.com/volio/go-common/util"
)

type Monitor interface {
	StartRoutine()
	FinishRoutine()
	Wait()
}

var (
	monitor Monitor = NewGraceful()
)

func Go(f func()) {
	monitor.StartRoutine()
	go func() {
		defer func() {
			util.Recovery()
			monitor.FinishRoutine()
		}()
		f()
	}()
}

func Wait() {
	WaitTimeout(5 * time.Second)
}

// Wait until all goroutine finish
func WaitTimeout(timeout time.Duration) {
	c := make(chan struct{})
	go func() {
		defer close(c)
		monitor.Wait()
	}()
	select {
	case <-c:
		log.L().Info("Graceful.Wait ok")
	case <-time.After(timeout):
		log.L().Error("Graceful.Wait timeout")
	}
}
