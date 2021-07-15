package gracefulshutdown

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gotway/gotway/pkg/log"
)

var shutdownSigs = []os.Signal{
	os.Interrupt,
	syscall.SIGINT,
	syscall.SIGTERM,
	syscall.SIGKILL,
	syscall.SIGHUP,
	syscall.SIGQUIT,
}

func WithGracefulShutdown(
	ctx context.Context,
	logger log.Logger,
) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	signals := make(chan os.Signal)
	signal.Notify(
		signals,
		shutdownSigs...,
	)
	go func() {
		s := <-signals
		logger.Info("received signal ", s.String())
		cancel()
		logger.Info("shutting down")
		<-time.After(time.Second)
	}()

	return ctx
}
