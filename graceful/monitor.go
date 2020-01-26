package graceful

import "sync"

type Monitor interface {
	StartRoutine()
	FinishRoutine()
	Wait()
}

func NewMonitor() Monitor {
	return new(graceful)
}

type graceful struct {
	wg sync.WaitGroup
}

func (g *graceful) StartRoutine() {
	g.wg.Add(1)
}

func (g *graceful) FinishRoutine() {
	g.wg.Done()
}

func (g *graceful) Wait() {
	g.wg.Wait()
}
