package api

import (
	"encoding/base64"
	"minecrafthead/internal/mojang"
	"minecrafthead/internal/service"
	"net/http"
)

type Handler struct {
	skinStore service.MinecraftSkinService
}

func NewHandler(skinStore service.MinecraftSkinService) *Handler {
	return &Handler{
		skinStore: skinStore,
	}
}

func (h *Handler) HeadHandler(w http.ResponseWriter, r *http.Request) {
	nick := r.PathValue("nickname")
	if nick == "" {
		http.Error(w, "missing nickname", http.StatusBadRequest)
		return
	}

	id, err := mojang.GetUUIDByNickname(nick)
	if err != nil {
		http.Error(w, "nickname not found: "+err.Error(), http.StatusNotFound)
		return
	}

	pngB64, err := h.skinStore.GetHead(r.Context(), id, 256, 256, false)
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
