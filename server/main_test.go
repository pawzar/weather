package main

import (
	"context"
	"errors"
	"os"
	"syscall"
	"testing"
	"time"

	"weather/internal/httpserver"
	"weather/internal/log"
)

func TestHttpServerErrorMessage(t *testing.T) {
	var logger log.TestLogger

	err := httpserver.Run(canceledContext(), errServer(true), &logger)

	if err == nil {
		t.Error("expecting error")
	}

	if logger != "shutting down server" {
		t.Errorf("wrong error logged in the background: %s", logger)
	}

	if err.Error() != "Shutdown" {
		t.Errorf("wrong error returned: %s", err)
	}
}

func canceledContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}

type errServer bool

func (s errServer) err(m string) error {
	if s {
		return errors.New(m)
	}

	return nil
}
func (s errServer) ListenAndServe() error          { return s.err("ListenAndServe") }
func (s errServer) Shutdown(context.Context) error { return s.err("Shutdown") }

func TestServerPanicOnMissingEnvVars(t *testing.T) {
	tcs := []struct {
		name        string
		setup       func()
		shouldPanic bool
	}{{
		name: "all env vars present",
		setup: func() {
			_ = os.Setenv("SERVER_PORT", "rand")
			_ = os.Setenv("APP_UNITS", "metric")
			_ = os.Setenv("APP_LANG", "pl")
			_ = os.Setenv("API_BASE_URL", "http://localhost")
			_ = os.Setenv("API_KEY", "key")
		},
	}, {
		name: "missing SERVER_PORT",
		setup: func() {
			_ = os.Setenv("APP_UNITS", "metric")
			_ = os.Setenv("APP_LANG", "pl")
			_ = os.Setenv("API_BASE_URL", "http://localhost")
			_ = os.Setenv("API_KEY", "key")
		},
	}, {
		name: "missing APP_UNITS",
		setup: func() {
			_ = os.Setenv("SERVER_PORT", "rand")
			_ = os.Setenv("APP_LANG", "pl")
			_ = os.Setenv("API_BASE_URL", "http://localhost")
			_ = os.Setenv("API_KEY", "key")
		},
	}, {
		name: "missing APP_LANG",
		setup: func() {
			_ = os.Setenv("SERVER_PORT", "rand")
			_ = os.Setenv("APP_UNITS", "metric")
			_ = os.Setenv("API_BASE_URL", "http://localhost")
			_ = os.Setenv("API_KEY", "key")
		},
	}, {
		name: "missing API_BASE_URL",
		setup: func() {
			_ = os.Setenv("SERVER_PORT", "rand")
			_ = os.Setenv("APP_UNITS", "metric")
			_ = os.Setenv("APP_LANG", "pl")
			_ = os.Setenv("API_KEY", "key")
		},
		shouldPanic: true,
	}, {
		name: "missing API_KEY",
		setup: func() {
			_ = os.Setenv("SERVER_PORT", "rand")
			_ = os.Setenv("APP_UNITS", "metric")
			_ = os.Setenv("APP_LANG", "pl")
			_ = os.Setenv("API_BASE_URL", "http://localhost")
		},
		shouldPanic: true,
	}}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			os.Clearenv()
			tc.setup()

			defer func() {
				if r := recover(); (r == nil) == tc.shouldPanic {
					t.Errorf("server panic should be [%t]", tc.shouldPanic)
				}
			}()

			go interrupt(time.Second)

			main()
		})
	}
}

func ExampleServerWithPort() {
	os.Clearenv()

	_ = os.Setenv("SERVER_PORT", "8888")
	_ = os.Setenv("APP_UNITS", "metric")
	_ = os.Setenv("APP_LANG", "pl")
	_ = os.Setenv("API_BASE_URL", "http://localhost")
	_ = os.Setenv("API_KEY", "key")

	go interrupt(time.Second)

	main()

	// Output:
	// DEBUG {ServerAddr::8888 Units:metric Lang:pl ApiBaseURL:http://localhost ApiKey:key}
	// DEBUG received "interrupt" signal
	// DEBUG shutting down server
}

func ExampleDefaultServer() {
	os.Clearenv()

	_ = os.Setenv("API_BASE_URL", "http://localhost")
	_ = os.Setenv("API_KEY", "key")

	go interrupt(time.Second)

	main()

	// Output:
	// DEBUG {ServerAddr::8080 Units:imperial Lang:en ApiBaseURL:http://localhost ApiKey:key}
	// DEBUG received "interrupt" signal
	// DEBUG shutting down server
}

func interrupt(ttl time.Duration) {
	time.Sleep(ttl)
	_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
}
