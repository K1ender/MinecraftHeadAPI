package main

import (
	"fmt"
	"log/slog"
	"minecrafthead/internal/api"
	"minecrafthead/internal/service"
	"minecrafthead/internal/util"
	"net/http"
	"os"

	"github.com/redis/go-redis/v9"
)

func main() {
	ip := util.GetIP(os.Getenv("REAL_IP") != "false")

	redisURL := os.Getenv("REDIS_URL")
	parsedURL, err := redis.ParseURL(redisURL)
	if err != nil {
		slog.Error(err.Error())
	}
	redis := redis.NewClient(parsedURL)
	NoCacheSkinStore := service.NewRedisCacheSkinStore(redis)
	handler := api.NewHandler(NoCacheSkinStore)

	http.HandleFunc("/head/{nickname}", handler.HeadHandler)
	slog.Info(
		fmt.Sprintf(
			"listening on http://%s:8080",
			ip,
		),
	)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		slog.Error("failed to start HTTP server", "err", err)
		os.Exit(1)
	}
}
