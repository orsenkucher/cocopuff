package gs

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// With graceful shutdown
func With(c context.Context) context.Context {
	ctx, cancel := context.WithCancel(c)
	done := make(chan struct{})
	signals := []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP}
	go listen(ctx, cancel, done, signals...)
	<-done
	return ctx
}

func listen(
	ctx context.Context,
	cancel context.CancelFunc,
	done chan struct{},
	signals ...os.Signal,
) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, signals...)
	close(done)
	select {
	case <-ctx.Done():
		return
	case <-ch:
		cancel()
		<-ch
		os.Exit(1)
	}
}
