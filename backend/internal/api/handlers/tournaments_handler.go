package handlers

import (
	"encoding/json"
	api_errors "golbugames/internal/api/errors"
	"golbugames/internal/sudoku/repository"
	"golbugames/pkg/types"
	"log/slog"
	"net/http"
)

func GetAllTournaments(w http.ResponseWriter, r *http.Request) {
	tournaments, err := repository.GetAllTournaments(r.Context())
	if err != nil {
		slog.Error("Failed to retrieve tournaments from the database", "error", err)
		api_errors.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(map[string]interface{}{
		"tournaments": tournaments,
	}); err != nil {
		slog.Error("Failed to encode response", "error", err)
		api_errors.WriteError(w, err)
		return
	}
}

func AddTournament(w http.ResponseWriter, r *http.Request) {
	var tournament types.Tournament
	if err := json.NewDecoder(r.Body).Decode(&tournament); err != nil {
		slog.Error("Failed to decode request body", "error", err)
		api_errors.WriteError(w, api_errors.ErrBadRequest)
		return
	}

	err := repository.AddTournament(r.Context(), tournament)
	if err != nil {
		slog.Error("Failed to insert tournament into the database", "error", err)
		api_errors.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Le tournoi a bien été créé",
	}); err != nil {
		slog.Error("Failed to encode response", "error", err)
		api_errors.WriteError(w, err)
		return
	}
}
