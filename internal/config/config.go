package config

import (
	"os"
	"strconv"
	"time"
)

const (
	PORT                   = "PORT"
	DB_CONN_STRING         = "DB_CONN_STRING"
	NATS_URL               = "NATS_URL"
	CACHE_TTL              = "CACHE_TTL"
	CACHE_CLEANUP_INTERVAL = "CACHE_CLEANUP_INTERVAL"
	DEBUG                  = "DEBUG"

	defaultPort                 = "8090"
	defaultConnDBString         = "host=postgres user=user password=secret dbname=orders_db sslmode=disable"
	defaultNatsURL              = "nats://nats:4222"
	defaultCacheTTL             = time.Duration(-1) // no expiration by default
	defaultCacheCleanupInterval = time.Duration(-1)

	defaultDebug = false
)

type Config struct {
	Port                 string
	DBConnString         string
	NatsURL              string
	CacheTTL             time.Duration
	CacheCleanupInterval time.Duration
	Debug                bool
}

func NewConfig() (*Config, error) {
	port := getEnv(PORT, defaultPort)
	dBconnStr := getEnv(DB_CONN_STRING, defaultConnDBString)
	natsURL := getEnv(NATS_URL, defaultNatsURL)
	cacheTTL := getEnvDuration(CACHE_TTL, defaultCacheTTL)
	cacheCleanupInterval := getEnvDuration(CACHE_CLEANUP_INTERVAL, defaultCacheCleanupInterval)

	debug := getEnvBool(DEBUG, defaultDebug)

	return &Config{
		Port:                 port,
		DBConnString:         dBconnStr,
		NatsURL:              natsURL,
		CacheTTL:             cacheTTL,
		CacheCleanupInterval: cacheCleanupInterval,
		Debug:                debug,
	}, nil
}

func getEnv(key string, defaultVal string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defaultVal
}

func getEnvBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}

func getEnvDuration(name string, defaultVal time.Duration) time.Duration {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseInt(valStr, 10, 64); err == nil {
		return time.Duration(val)
	}

	return defaultVal
}
