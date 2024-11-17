package config

import (
	"fmt"
	"os"
)

func readEnv(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		fmt.Printf("[CONFIG] %s not set\n", key)
		return val
	} else {
		fmt.Printf("[CONFIG] %s=%s\n", key, val)
		return ""
	}
}

func readEnvWithFallback(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
