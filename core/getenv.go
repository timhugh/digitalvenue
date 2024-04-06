package core

import (
	"fmt"
	"os"
)

func RequireEnv(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return "", fmt.Errorf("missing required environment variable %s", key)
	}
	return val, nil
}
