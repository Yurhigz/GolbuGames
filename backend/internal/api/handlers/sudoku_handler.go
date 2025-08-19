package handlers

import (
	"encoding/json"
	"golbugames/backend/internal/sudoku"
	"golbugames/backend/internal/sudoku/repository"
	"golbugames/backend/pkg/types"
	"golbugames/backend/pkg/utils"
	"log"
	"net/http"
)

func AddGrid(w http.ResponseWriter, r *http.Request) {

	var req types.GridRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		log.Printf("Erreur : %v", err)
		http.Error(w, "invalid data format", http.StatusBadRequest)
		return
	}
	difficulty := req.Difficulty
	if difficulty == "" {
		difficulty = "easy"
	}

	validDifficulties := map[string]bool{
		"easy":         true,
		"intermediate": true,
		"advanced":     true,
		"expert":       true,
	}

	if !validDifficulties[difficulty] {
		http.Error(w, "Invalid difficulty level", http.StatusBadRequest)
		return
	}

	solvedGrid, _ := sudoku.GenerateSolvedGrid()
	savedSolvedGrid := solvedGrid
	playableGrid, _ := sudoku.GeneratePlayableGrid(solvedGrid, difficulty)

	boardStr := utils.GridTransformer(playableGrid)
	solutionStr := utils.GridTransformer(savedSolvedGrid)

	err = repository.AddGridDB(r.Context(), boardStr, solutionStr, difficulty)
	if err != nil {
		log.Printf("Error saving grid to DB: %v", err)
		http.Error(w, "Failed to save grid", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		map[string]string{
			"message":    "Grid successfully created",
			"board":      boardStr,
			"solution":   solutionStr,
			"difficulty": difficulty,
		})

}

func GetGrid(w http.ResponseWriter, r *http.Request) {

	difficulty := r.URL.Query().Get("difficulty")
	if difficulty == "" {
		difficulty = "easy"
	}

	validDifficulties := map[string]bool{
		"easy":         true,
		"intermediate": true,
		"advanced":     true,
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
		"message":    "Grid sucessfully retrieved",
		"board":      sudokuGrid.Board,
		"difficulty": sudokuGrid.Difficulty,
	})
}

func GetLeaderboard(w http.ResponseWriter, r *http.Request) {
    leaderboard, err := repository.GetLeaderboard(r.Context())
    if err != nil {
        http.Error(w, "Internal retrieval error", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(leaderboard)
}

// func GetUserHistory(w http.ResponseWriter, r *http.Request) {
// 	// Historique des parties
// 	// Progression
// }

// func SaveGameProgress(w http.ResponseWriter, r *Request) {
// 	// Sauvegarde l'état actuel
// 	// Permet de reprendre plus tard
// }
