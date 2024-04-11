package core

import (
	"os"
)

type MissingEnvError struct {
	Key string
}

func (e MissingEnvError) Error() string {
	return "missing required environment variable " + e.Key
}

func RequireEnv(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return "", MissingEnvError{Key: key}
	}
	return val, nil
}
