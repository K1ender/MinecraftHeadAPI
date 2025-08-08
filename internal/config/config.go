package config

import (
	"log/slog"
	"os"
	"reflect"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	UseRealIP bool
	Cache     CacheConfig
}

type CacheType string

const (
	CacheTypeRedis CacheType = "redis"
	CacheTypeNone  CacheType = "none"
)

type CacheConfig struct {
	Type CacheType

	Redis struct {
		URL string
	}
}

func MustInit() *Config {
	if err := godotenv.Load(".env"); err != nil {
		slog.Warn("Failed to load .env file", "err", err)
	}
	cfg := Config{
		UseRealIP: getVariable("REAL_IP", false),
		Cache: CacheConfig{
			Type: getVariable("CACHE_TYPE", CacheTypeNone),
			Redis: struct {
				URL string
			}{
				URL: getVariable("REDIS_URL", ""),
			},
		},
	}
	return &cfg
}

func getVariable[T any](key string, fallback T) T {
	valStr, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	var t T
	typ := reflect.TypeOf(t)

	switch typ.Kind() {
	case reflect.Bool:
		val, err := strconv.ParseBool(valStr)
		if err != nil {
			return fallback
		}
		return any(val).(T)
	case reflect.Int:
		val, err := strconv.Atoi(valStr)
		if err != nil {
			return fallback
		}
		return any(val).(T)
	case reflect.String:
		return any(valStr).(T)
	default:
		panic("unsupported type")
	}
}
