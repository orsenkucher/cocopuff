package wg

import "sync"

type G struct {
	wg sync.WaitGroup
}

func (g *G) Add(f func()) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		f()
	}()
}

func (g *G) Wait() {
	g.wg.Wait()
}
