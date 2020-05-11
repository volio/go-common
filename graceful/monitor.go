package graceful

import "sync"

func NewGraceful() *Graceful {
	return new(Graceful)
}

type Graceful struct {
	wg sync.WaitGroup
}

func (g *Graceful) StartRoutine() {
	g.wg.Add(1)
}

func (g *Graceful) FinishRoutine() {
	g.wg.Done()
}

func (g *Graceful) Wait() {
	g.wg.Wait()
}
