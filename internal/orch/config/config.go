package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port string

	TimeAdditionMs        int
	TimeSubtractionMs     int
	TimeMultiplicationsMs int
	TimeDivisionsMs       int
}

func getEnvDefault(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func FromEnv() *Config {
	config := &Config{}
	config.Port = getEnvDefault("PORT", "8080")

	config.TimeAdditionMs, _ = strconv.Atoi(getEnvDefault("TIME_ADDITION_MS", "1000"))
	config.TimeSubtractionMs, _ = strconv.Atoi(getEnvDefault("TIME_SUBTRACTION_MS", "1000"))
	config.TimeMultiplicationsMs, _ = strconv.Atoi(getEnvDefault("TIME_MULTIPLICATIONS_MS", "1000"))
	config.TimeDivisionsMs, _ = strconv.Atoi(getEnvDefault("TIME_DIVISIONS_MS", "1000"))

	return config
}
