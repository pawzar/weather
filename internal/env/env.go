package env

import (
	"fmt"
	"os"
)

func NonEmptyEnvVar(key string) string {
	s := os.Getenv(key)

	if s == "" {
		panic(fmt.Sprintf("%s env var is empty", key))
	}

	return s
}

func EnvVar(key, def string) string {
	s := os.Getenv(key)

	if s == "" {
		return def
	}

	return s
}
