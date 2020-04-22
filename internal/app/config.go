package app

import (
	"strconv"

	"weather/internal/env"
	"weather/internal/net"
)

type Config struct {
	ServerAddr string
	Units      string
	Lang       string
	ApiBaseURL string
	ApiKey     string
}

func Configure() Config {
	port := env.EnvVar("SERVER_PORT", "8080")
	if port == "rand" {
		port = strconv.Itoa(net.Must(net.GetFreePort()))
	}

	return Config{
		ServerAddr: ":" + port,
		Units:      env.EnvVar("APP_UNITS", "imperial"),
		Lang:       env.EnvVar("APP_LANG", "en"),
		ApiBaseURL: env.NonEmptyEnvVar("API_BASE_URL"),
		ApiKey:     env.NonEmptyEnvVar("API_KEY"),
	}
}
