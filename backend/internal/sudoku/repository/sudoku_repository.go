package repository

import (
	"context"
	"fmt"
	"golbugames/backend/internal/database"
	"golbugames/backend/pkg/types"
	"log"
	"math/rand/v2"
	"time"
)

func AddGridDB(parentsContext context.Context, board, solution, difficulty string) error {
	ctx, cancel := context.WithTimeout(parentsContext, 2*time.Second)
	defer cancel()

	query := `INSERT INTO sudoku_games (board, solution, difficulty) VALUES ($1, $2, $3)`

	_, err := database.DBPool.Exec(ctx, query, board, solution, difficulty)

	if err != nil {
		return fmt.Errorf("[AddGrid] Error inserting sudoku grid : %v", err)
	}

	log.Printf("Sudoku grid added successfully with difficulty : %s", difficulty)
	return nil
}

func GetGridDB(parentsContext context.Context, id int) (string, string, error) {
	ctx, cancel := context.WithTimeout(parentsContext, 2*time.Second)
	defer cancel()

	query := `SELECT board,solution FROM sudoku_games WHERE id = $1`
	var board, solution string

	err := database.DBPool.QueryRow(ctx, query, id).Scan(&board, &solution)
	if err != nil {
		return "", "", fmt.Errorf("[GetGrid] cannot find grid %d : %v", id, err)
	}

	log.Printf("Sudoku grid (id: %d) successfully retrieved", id)
	return board, solution, nil
}

func GetRandomGridDB(parentsContext context.Context, difficulty string) (*types.SudokuGrid, error) {
	ctx, cancel := context.WithTimeout(parentsContext, 2*time.Second)
	defer cancel()

	// Première requête pour compter
	var count int
	countQuery := `SELECT COUNT(*) FROM sudoku_games WHERE difficulty = $1`
	err := database.DBPool.QueryRow(ctx, countQuery, difficulty).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("[GetRandomGrid] failed to count grids: %v", err)
	}

	if count == 0 {
		return nil, fmt.Errorf("[GetRandomGrid] no grid found for difficulty %s", difficulty)
	}

	offset := rand.IntN(count)

	var sudokuGrid types.SudokuGrid
	query := `
        SELECT board, solution, difficulty
        FROM sudoku_games 
        WHERE difficulty = $1 
        LIMIT 1 OFFSET $2`

	err = database.DBPool.QueryRow(ctx, query, difficulty, offset).Scan(&sudokuGrid.Board, &sudokuGrid.Solution, &sudokuGrid.Difficulty)
	if err != nil {
		return nil, fmt.Errorf("[GetRandomGrid] failed to get random grid: %w", err)
	}

	return &sudokuGrid, nil
}

// func GetLeaderboard(w http.ResponseWriter, r *http.Request) {
// 	// Classement des meilleurs joueurs
// 	// Filtrage par difficulté et ELO
// }

// func GetUserHistory(w http.ResponseWriter, r *http.Request) {
// 	// Historique des parties
// 	// Progression
// }

// func SaveGameProgress(w http.ResponseWriter, r *http.Request) {
// 	// Sauvegarde l'état actuel
// 	// Permet de reprendre plus tard
// }
