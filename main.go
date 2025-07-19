package main

import (
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/redis/go-redis/v9"
)

func main() {
	ip := GetIP(os.Getenv("REAL_IP") != "false")

	redisURL := os.Getenv("REDIS_URL")
	parsedURL, err := redis.ParseURL(redisURL)
	if err != nil {
		slog.Error(err.Error())
	}
	redis := redis.NewClient(parsedURL)
	NoCacheSkinStore := NewRedisCacheSkinStore(redis)
	handler := NewHandler(NoCacheSkinStore)

	http.HandleFunc("/head/{nickname}", handler.headHandler)
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

type Handler struct {
	skinStore MinecraftSkinService
}

func NewHandler(skinStore MinecraftSkinService) *Handler {
	return &Handler{
		skinStore: skinStore,
	}
}

func (h *Handler) headHandler(w http.ResponseWriter, r *http.Request) {
	nick := r.PathValue("nickname")
	if nick == "" {
		http.Error(w, "missing nickname", http.StatusBadRequest)
		return
	}

	id, err := getUUIDByNickname(nick)
	if err != nil {
		http.Error(w, "nickname not found: "+err.Error(), http.StatusNotFound)
		return
	}

	pngB64, err := h.skinStore.getHead(r.Context(), id, 256, 256, false)
	if err != nil {
		http.Error(w, "failed to render head: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := base64.StdEncoding.DecodeString(pngB64)
	if err != nil {
		http.Error(w, "failed to decode base64", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Write(data)
}
