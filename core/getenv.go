package core

import (
	"github.com/pkg/errors"
	"os"
)

func RequireEnv(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return "", errors.Errorf("missing required environment variable %s", key)
	}
	return val, nil
}
