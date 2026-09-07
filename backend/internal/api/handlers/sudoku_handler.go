package handlers

import (
	"encoding/json"
	api_errors "golbugames/internal/api/errors"
	"golbugames/internal/sudoku"
	"golbugames/internal/sudoku/repository"
	"golbugames/pkg/types"
	"golbugames/pkg/utils"
	"log/slog"
	"net/http"
)

func AddGrid(w http.ResponseWriter, r *http.Request) {

	var req types.GridRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		slog.Error("Failed to decode request body", "error", err)
		api_errors.WriteError(w, err)
		return
	}

	difficulty := req.Difficulty
	if difficulty == "" {
		difficulty = "easy"
	}

	validDifficulties := map[string]bool{
		"easy":         true,
		"intermediate": true,
		"hard":         true,
		"expert":       true,
	}

	if !validDifficulties[difficulty] {
		slog.Error("Invalid difficulty level", "difficulty", difficulty)
		api_errors.WriteError(w, api_errors.ErrInvalidDifficulty)
		return
	}

	solvedGrid, err := sudoku.GenerateSolvedGrid()
	if err != nil {
		slog.Error("Failed to generate solved grid", "error", err)
		api_errors.WriteError(w, err)
		return
	}
	savedSolvedGrid := solvedGrid
	playableGrid, err := sudoku.GeneratePlayableGrid(solvedGrid, difficulty)
	if err != nil {
		slog.Error("Failed to generate playable grid", "error", err)
		api_errors.WriteError(w, err)
		return
	}

	boardStr := utils.GridTransformer(playableGrid)
	solutionStr := utils.GridTransformer(savedSolvedGrid)

	err = repository.AddGridDB(r.Context(), boardStr, solutionStr, difficulty)
	if err != nil {
		slog.Error("Failed to save grid to DB", "error", err)
		api_errors.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(
		map[string]string{
			"message":    "Grid successfully created",
			"board":      boardStr,
			"solution":   solutionStr,
			"difficulty": difficulty,
		}); err != nil {
		slog.Error("Failed to encode response", "error", err)
		api_errors.WriteError(w, err)
		return
	}

}

// En cours de traitement avec le Websocket
func GetGrid(w http.ResponseWriter, r *http.Request) {

	difficulty := r.URL.Query().Get("difficulty")
	if difficulty == "" {
		difficulty = "easy"
	}

	validDifficulties := map[string]bool{
		"easy":         true,
		"intermediate": true,
		"hard":         true,
		"expert":       true,
	}

	if !validDifficulties[difficulty] {
		http.Error(w, "Invalid difficulty level", http.StatusBadRequest)
		return
	}

	sudokuGrid, err := repository.GetRandomGridDB(r.Context(), difficulty)

	if err != nil {
		http.Error(w, "Internal retrieval error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Grid sucessfully retrieved",
		// 		"board":      sudokuGrid.Board,
		"difficulty": sudokuGrid.Difficulty,
	})
}

func GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	leaderboard, err := repository.GetLeaderboard(r.Context())
	if err != nil {
		slog.Error("Failed to retrieve leaderboard", "error", err)
		api_errors.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(leaderboard); err != nil {
		slog.Error("Failed to encode response", "error", err)
		api_errors.WriteError(w, err)
		return
	}
}

// func GetUserHistory(w http.ResponseWriter, r *http.Request) {
// 	// Historique des parties
// 	// Progression
// }

// func SaveGameProgress(w http.ResponseWriter, r *Request) {
// 	// Sauvegarde l'état actuel
// 	// Permet de reprendre plus tard
// }
