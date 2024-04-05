package core

import (
	"github.com/rs/zerolog/log"
	"os"
)

func Getenv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatal().Msgf("Environment variable %s is required", key)
	}
	return val
}
