package main

import (
	"fmt"
	"log/slog"
	"minecrafthead/internal/api"
	"minecrafthead/internal/config"
	"minecrafthead/internal/service"
	"minecrafthead/internal/util"
	"net/http"
	"os"

	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.MustInit()
	ip := util.GetIP(cfg.UseRealIP)

	var handler *api.Handler
	var CacheSkinStore service.MinecraftSkinService
	var CacheUUIDStore service.MinecraftNicknameUUIDService

	// FIXME: Thats not the best way? Use factory or something like that
	if cfg.Cache.Type == config.CacheTypeRedis {
		redisClient := redis.NewClient(&redis.Options{Addr: cfg.Cache.Redis.URL})
		CacheSkinStore = service.NewRedisCacheSkinStore(redisClient)
		CacheUUIDStore = service.NewRedisCacheUUIDService(redisClient)
	} else {
		CacheSkinStore = service.NewNoCacheSkinStore()
		CacheUUIDStore = service.NewNoCacheUUIDService()
	}

	handler = api.NewHandler(CacheSkinStore, CacheUUIDStore)

	http.HandleFunc("/head/{nickname}", handler.HeadHandler)
	slog.Info(
		fmt.Sprintf(
			"listening on http://%s:8080",
			ip,
		),
	)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		slog.Error("failed to start HTTP server", "err", err)
		os.Exit(1)
	}
}
