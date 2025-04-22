package handlers

import (
	"encoding/json"
	"golbugames/backend/internal/sudoku/repository"
	"golbugames/backend/pkg/types"
	"log"
	"net/http"
	"strconv"
)

func SubmitSoloGame(w http.ResponseWriter, r *http.Request) {
	var game types.Game
	err := json.NewDecoder(r.Body).Decode(&game)
	if err != nil {
		http.Error(w, "invalid data format", http.StatusBadRequest)
		return
	}
	if game.GameMode == "1v1" || game.OpponentID != nil || game.Results != nil {
		http.Error(w, "Invalid game mode or results for solo game", http.StatusBadRequest)
		return
	}

	err = repository.SubmitSoloGameDB(r.Context(), game.UserID, game.Completion_time)
	if err != nil {
		log.Printf("Error submitting game for user %d: %v", game.UserID, err)
		http.Error(w, "Failed to submit game", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":         "Game submitted successfully",
		"userId":          strconv.Itoa(game.UserID),
		"completion_time": strconv.Itoa(game.Completion_time),
	})

}

func SubmitMultiGame(w http.ResponseWriter, r *http.Request) {
	var game types.Game
	err := json.NewDecoder(r.Body).Decode(&game)
	if err != nil {
		http.Error(w, "invalid data format", http.StatusBadRequest)
		return
	}
	if game.GameMode != "1v1" || game.OpponentID == nil || game.Results == nil {
		http.Error(w, "Invalid game mode or results for multiplayer game", http.StatusBadRequest)
		return
	}
	err = repository.SubmitMultiGameDB(r.Context(), game.UserID, *game.OpponentID, *game.Results, game.Completion_time)
	if err != nil {
		log.Printf("Error submitting game for user %d: %v", game.UserID, err)
		http.Error(w, "Failed to submit game", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":    "Game submitted successfully",
		"userId":     strconv.Itoa(game.UserID),
		"opponentId": strconv.Itoa(*game.OpponentID),
		"score":      strconv.Itoa(*game.Results),
	})
}
