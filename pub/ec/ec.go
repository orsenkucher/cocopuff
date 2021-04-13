package ec

import "github.com/orsenkucher/cocopuff/pub/wg"

// Go <-f()
func Go(f func() error) <-chan error {
	ch := make(chan error)

	wg := wg.G{}
	wg.Add(func() {
		ch <- f()
	})

	go func() {
		wg.Wait()
		close(ch)
	}()

	return ch
}
