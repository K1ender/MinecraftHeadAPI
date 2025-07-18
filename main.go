package main

import (
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	ip := GetIP(os.Getenv("REAL_IP") != "false")

	http.HandleFunc("/head/{nickname}", headHandler)
	slog.Info(
		fmt.Sprintf(
			"listening on http://%s:8080",
			ip,
		),
	)

	http.ListenAndServe(":8080", nil)
}

func headHandler(w http.ResponseWriter, r *http.Request) {
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

	pngB64, err := getHead64(id, 256, 256, true)
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
