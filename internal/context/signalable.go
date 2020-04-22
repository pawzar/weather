package context

import (
	"context"
	"os"
	"os/signal"

	"weather/internal/log"
)

func Signalable(logger log.Logger) context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)
	signals := []os.Signal{os.Interrupt}

	signal.Notify(c, signals...)

	go func() {
		s := <-c

		logger.Debugf(`received "%s" signal`, s)

		cancel()
	}()

	return ctx
}
