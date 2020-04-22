package main

import (
	"weather/internal/app"
	"weather/internal/context"
	"weather/internal/log"
	"weather/internal/run"
)

func main() {
	config := app.Configure()

	logger := &log.StdLogger{}
	logger.Debugf("%+v", config)

	ctx := context.Signalable(logger)

	run.Server(ctx, config, logger)
}
