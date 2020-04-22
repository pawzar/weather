// +build integration

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"weather/internal/app"
	"weather/internal/env"
	"weather/internal/log"
	"weather/internal/net"
	"weather/internal/run"
)

var testingBaseURL string

func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())
	config := app.Config{
		ServerAddr: ":" + strconv.Itoa(net.Must(net.GetFreePort())),
		Units:      "metric",
		Lang:       "pl",
		ApiBaseURL: "https://api.openweathermap.org/data/2.5",
		ApiKey:     env.NonEmptyEnvVar("API_KEY"),
	}
	logger := &log.QuietLogger{}

	go run.Server(ctx, config, logger)

	testingBaseURL = fmt.Sprintf("http://localhost%s", config.ServerAddr)
	for _, err := http.Get(testingBaseURL); err != nil; {
		fmt.Printf("waiting for %s\n", testingBaseURL)
		time.Sleep(time.Second)
	}

	exitCode := m.Run()

	cancel()

	os.Exit(exitCode)
}

func Test_ExpectedStatuses(t *testing.T) {
	tcs := []struct {
		name          string
		url           string
		code          int
		needEmptyBody bool
	}{
		{
			name:          "bad request",
			url:           testingBaseURL,
			code:          400,
			needEmptyBody: true,
		},
		{
			name:          "not found",
			url:           fmt.Sprintf("%s/?c=%s", testingBaseURL, "hobbiton"),
			code:          200,
			needEmptyBody: true,
		},
		{
			name: "ok",
			url:  fmt.Sprintf("%s/?c=%s", testingBaseURL, "warsaw"),
			code: 200,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			response, err := http.Get(tc.url)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
				return
			}

			if response.StatusCode != tc.code {
				t.Errorf("unexpected status code: %d instead of %d", response.StatusCode, tc.code)
				return
			}

			if tc.needEmptyBody && response.ContentLength > 0 {
				t.Errorf("unexpected content length: %d instead of 0", response.ContentLength)
				return
			}

			if !tc.needEmptyBody && response.ContentLength == 0 {
				t.Error("unexpectedly empty body")
				return
			}

		})
	}
}
