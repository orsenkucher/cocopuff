package ec

import "github.com/orsenkucher/cocopuff/pub/wait"

// Go run function and return in channel
func Go(f func() error) <-chan error {
	ch := make(chan error)

	wg := wait.Group{}
	wg.Add(func() {
		ch <- f()
	})

	go func() {
		wg.Wait()
		close(ch)
	}()

	return ch
}
