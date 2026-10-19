package handlers

import (
	"encoding/json"
	"golbugames/internal/websocket/multiplayer"
	"log/slog"
	"net/http"
)

func GetMultiplayerMetrics(hubManager *multiplayer.HubManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(hubManager.MetricsSnapshot()); err != nil {
			slog.Error("Failed to encode multiplayer metrics", "error", err)
			http.Error(w, "failed to encode multiplayer metrics", http.StatusInternalServerError)
		}
	}
}
