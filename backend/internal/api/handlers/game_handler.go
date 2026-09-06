package handlers

import (
	"encoding/json"
	api_errors "golbugames/internal/api/errors"
	"golbugames/internal/sudoku/repository"
	"golbugames/pkg/types"
	"log/slog"
	"net/http"
	"strconv"
)

func SubmitSoloGame(w http.ResponseWriter, r *http.Request) {
	var game types.Game
	err := json.NewDecoder(r.Body).Decode(&game)
	if err != nil {
		slog.Error("Failed to decode request body", "error", err)
		api_errors.WriteError(w, err)
		return
	}

	if game.UserID == "" || game.GameMode != "solo" || game.Completion_time <= 0 || game.OpponentID != "" || game.Results != nil || game.Difficulty == "" {
		slog.Error("Missing required fields for solo game", "game_mode", game.GameMode, "completion_time", game.Completion_time, "opponent_id", game.OpponentID, "results", game.Results)
		api_errors.WriteError(w, api_errors.ErrBadRequest)
		return
	}

	err = repository.SubmitSoloGameDB(r.Context(), game.UserID, game.Completion_time)
	if err != nil {
		slog.Error("Error submitting game for user %v: %v", game.UserID, err)
		api_errors.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(map[string]string{
		"message":         "Game submitted successfully",
		"userId":          game.UserID,
		"completion_time": strconv.Itoa(game.Completion_time),
	}); err != nil {
		slog.Error("Failed to encode response", "error", err)
		api_errors.WriteError(w, err)
		return
	}

}

func SubmitMultiGame(w http.ResponseWriter, r *http.Request) {
	var game types.Game
	err := json.NewDecoder(r.Body).Decode(&game)
	if err != nil {
		slog.Error("Failed to decode request body", "error", err)
		api_errors.WriteError(w, api_errors.ErrBadRequest)
		return
	}
	if game.UserID == "" || game.GameMode != "1v1" || game.OpponentID == "" || game.Completion_time <= 0 || game.Results == nil || game.Difficulty == "" {
		slog.Error("Missing required fields for multiplayer game", "game_mode", game.GameMode, "opponent_id", game.OpponentID, "completion_time", game.Completion_time, "results", game.Results)
		api_errors.WriteError(w, api_errors.ErrBadRequest)
		return
	}
	err = repository.SubmitMultiGameDB(r.Context(), game.UserID, game.OpponentID, *game.Results, game.Completion_time)
	if err != nil {
		slog.Error("Error submitting game for user %v: %v", game.UserID, err)
		api_errors.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(map[string]string{
		"message":    "Game submitted successfully",
		"userId":     game.UserID,
		"opponentId": game.OpponentID,
		"score":      strconv.Itoa(*game.Results),
	}); err != nil {
		slog.Error("Failed to encode response", "error", err)
		api_errors.WriteError(w, err)
		return
	}
}
