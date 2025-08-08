package api

import (
	"encoding/base64"
	"minecrafthead/internal/service"
	"net/http"
	"strconv"
)

type Handler struct {
	skinStore service.MinecraftSkinService
	uuidStore service.MinecraftNicknameUUIDService
}

func NewHandler(skinStore service.MinecraftSkinService, uuidStore service.MinecraftNicknameUUIDService) *Handler {
	return &Handler{
		skinStore: skinStore,
		uuidStore: uuidStore,
	}
}

func (h *Handler) HeadHandler(w http.ResponseWriter, r *http.Request) {
	nick := r.PathValue("nickname")
	if nick == "" {
		http.Error(w, "missing nickname", http.StatusBadRequest)
		return
	}

	overlayQuery := r.URL.Query().Get("overlay")
	if overlayQuery == "" {
		overlayQuery = "false"
	}
	overlay, err := strconv.ParseBool(overlayQuery)
	if err != nil {
		http.Error(w, "invalid overlay", http.StatusBadRequest)
		return
	}

	id, err := h.uuidStore.GetUUIDByNickname(r.Context(), nick)
	if err != nil {
		http.Error(w, "nickname not found: "+err.Error(), http.StatusNotFound)
		return
	}

	pngB64, err := h.skinStore.GetHead(r.Context(), id, 256, 256, overlay)
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
