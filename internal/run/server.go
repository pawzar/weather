package run

import (
	"context"

	"weather/internal/app"
	"weather/internal/httpserver"
	"weather/internal/log"
	"weather/internal/weather"
	"weather/internal/weather/current"
)

func Server(ctx context.Context, config app.Config, logger log.Logger) {
	httpService := current.NewHttpService(
		weather.NewClient(
			config.ApiBaseURL,
			config.ApiKey,
			config.Units,
			config.Lang,
			logger,
		),
	)
	inMemoryCache := make(current.InMemoryCache)

	handlerFunc := current.MultipleCitySearch(current.NewCachingWrapper(httpService, inMemoryCache), logger)

	server := httpserver.ParametrisedServer(config.ServerAddr, handlerFunc)

	if err := httpserver.Run(ctx, server, logger); err != nil {
		logger.Errorf("%s", err)
	}
}
