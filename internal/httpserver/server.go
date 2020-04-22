package httpserver

import (
	"context"
	"net/http"
	"time"

	"weather/internal/log"
)

type Server interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

func ParametrisedServer(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: handler,
	}
}

func Run(ctx context.Context, server Server, logger log.Logger) (err error) {
	go start(server, logger)

	<-ctx.Done()

	return shutdown(server, logger)
}

func start(s Server, logger log.Logger) {
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("%s", err)
	}
}

func shutdown(s Server, logger log.Logger) error {
	logger.Debug("shutting down server")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.Shutdown(ctxShutDown)
}
