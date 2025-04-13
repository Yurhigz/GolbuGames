package sudoku

import (
	"context"
	"fmt"
	"golbugames/backend/internal/database"
	"golbugames/backend/pkg/types"
	"log"
	"math/rand/v2"
	"time"

	"github.com/jackc/pgx/v5"
)

func AddUser(parentsContext context.Context, username, accountname, password string) error {
	ctx, cancel := context.WithTimeout(parentsContext, 2*time.Second)
	defer cancel()

	// Usage de transaction pour garantir l'intégrité des données car il y a deux requêtes
	tx, err := database.DBPool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("[AddUser] cannot start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// vérification de l'unicité de l'utilisateur
	var exists bool
	err = tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", username).Scan(&exists)
	if exists {
		return fmt.Errorf("[AddUser] username <%v> already exists", username)
	}

	// Insérer l'utilisateur
	var userId int
	err = tx.QueryRow(ctx,
		`INSERT INTO users (username, accountname, password) 
		 VALUES ($1, $2, $3) 
		 RETURNING id`,
		username, accountname, password).Scan(&userId)
	if err != nil {
		return fmt.Errorf("[AddUser] Error inserting user [%s]: %w", username, err)
	}

	// Initialiser les stats
	_, err = tx.Exec(ctx,
		`INSERT INTO user_stats 
		 (user_id, total_games, total_wins, total_losses, total_draws, total_time, average_time) 
		 VALUES ($1, 0, 0, 0, 0, 0, 0)`,
		userId)
	if err != nil {
		return fmt.Errorf("[AddUser] Error initializing user stats: %w", err)
	}

	// Commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("[AddUser] Error committing transaction: %w", err)
	}

	log.Printf("User added successfully with initialized stats: %s", username)
	return nil
}

func DeleteUser(parentsContext context.Context, id_user int) error {
	ctx, cancel := context.WithTimeout(parentsContext, 2*time.Second)
	defer cancel()

	query := `DELETE FROM users WHERE id = $1`

	_, err := database.DBPool.Exec(ctx, query, id_user)

	if err != nil {
		return fmt.Errorf("[DeleteUser] Error deleting user (ID: %d): %v", id_user, err)
	}

	log.Printf("User (ID: %d) deleted successfully:", id_user)
	return nil
}

func GetUser(parentsContext context.Context, id_user int) (*types.User, error) {
	ctx, cancel := context.WithTimeout(parentsContext, 2*time.Second)
	defer cancel()

	var user types.User
	query := `SELECT username,accountname,password,id FROM users WHERE id = $1`

	err := database.DBPool.QueryRow(ctx, query, id_user).Scan(&user.Username, &user.Accountname, &user.Password, &user.ID)

	if err != nil {
		return nil, fmt.Errorf("[GetUser] Error retrieving user %v : %v", id_user, err)
	}

	return &user, nil
}

func UpdateUserPassword(parentsContext context.Context, id_user int, password string) error {
	ctx, cancel := context.WithTimeout(parentsContext, 2*time.Second)
	defer cancel()

	query := `UPDATE users SET password = $1 WHERE id = $2`

	_, err := database.DBPool.Exec(ctx, query, password, id_user)

	if err != nil {
		return fmt.Errorf("[UpdateUser] Error updating user %v : %v", id_user, err)
	}

	log.Printf("User (ID: %d) updated successfully:", id_user)
	return nil

}

func AddGrid(parentsContext context.Context, board, solution, difficulty string) error {
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

func GetGrid(parentsContext context.Context, id int) (string, string, error) {
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

func GetRandomGrid(parentsContext context.Context, difficulty string) (*types.SudokuGrid, error) {
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

func updateUserStats(ctx context.Context, tx pgx.Tx, userId int, win, loss, draw bool, completionTime int, isSolo bool) error {
	query := `
        UPDATE user_stats 
        SET total_games = total_games + $6,
            total_wins = total_wins + $1,
            total_losses = total_losses + $2,
            total_draws = total_draws + $3,
            total_time = total_time + $4,
            average_time = (total_time + $4) / (total_games + 1),
            total_solo_games = total_solo_games + $5,
        WHERE user_id = $7`

	winInt := 0
	lossInt := 0
	drawInt := 0
	soloInt := 0
	multiInt := 0

	if win {
		winInt = 1
	}
	if loss {
		lossInt = 1
	}
	if draw {
		drawInt = 1
	}
	if isSolo {
		soloInt = 1
	} else {
		multiInt = 1
	}

	_, err := tx.Exec(ctx, query, winInt, lossInt, drawInt, completionTime, soloInt, multiInt, userId)
	if err != nil {
		return fmt.Errorf("failed to update user stats: %w", err)
	}
	return nil
}

func SubmitSoloGame(parentsContext context.Context, userId, completionTime int) error {
	ctx, cancel := context.WithTimeout(parentsContext, 2*time.Second)
	defer cancel()
	// Usage des transactions car double requête
	tx, err := database.DBPool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("[SubmitSoloGame] cannot start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Insérer le score
	_, err = tx.Exec(ctx,
		`INSERT INTO games_scores (user_id, game_mode, completion_time) 
		 VALUES ($1, 'solo', $2)`,
		userId, completionTime)
	if err != nil {
		return fmt.Errorf("[SubmitSoloGame] cannot submit the game results: %w", err)
	}

	// Mettre à jour les stats (considéré comme une victoire en solo)
	err = updateUserStats(ctx, tx, userId, true, false, false, completionTime, true)
	if err != nil {
		return fmt.Errorf("[SubmitSoloGame] cannot update user stats: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("[SubmitSoloGame] cannot commit transaction: %w", err)
	}

	log.Printf("Solo game successfully submitted and stats updated for user %d", userId)
	return nil
}

func SubmitMultiGame(parentsContext context.Context, user1, user2 int, results, completionTime int) error {
	ctx, cancel := context.WithTimeout(parentsContext, 2*time.Second)
	defer cancel()
	// Usage des transactions car double requête
	tx, err := database.DBPool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("[SubmitMultiGame] cannot start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Insérer le score
	_, err = tx.Exec(ctx,
		`INSERT INTO games_scores (user_id, opponent_id, game_mode, results, completion_time) 
		 VALUES ($1, $2, '1v1', $3, $4)`,
		user1, user2, results, completionTime)
	if err != nil {
		return fmt.Errorf("[SubmitMultiGame] cannot submit the game results: %w", err)
	}

	// Mettre à jour les stats des deux joueurs selon le résultat
	isWin1 := results == 0
	isDraw := results == 1
	isWin2 := results == 2

	// Stats joueur 1
	err = updateUserStats(ctx, tx, user1, isWin1, isWin2, isDraw, completionTime, false)
	if err != nil {
		return fmt.Errorf("[SubmitMultiGame] cannot update user1 stats: %w", err)
	}

	// Stats joueur 2
	err = updateUserStats(ctx, tx, user2, isWin2, isWin1, isDraw, completionTime, false)
	if err != nil {
		return fmt.Errorf("[SubmitMultiGame] cannot update user2 stats: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("[SubmitMultiGame] cannot commit transaction: %w", err)
	}

	log.Printf("Multi game successfully submitted and stats updated")
	return nil
}

func GetUserStats(parentsContext context.Context, id_user int) (*types.UserStats, error) {
	ctx, cancel := context.WithTimeout(parentsContext, 2*time.Second)
	defer cancel()

	_, err := GetUser(parentsContext, id_user)
	if err != nil {
		return nil, fmt.Errorf("[GetUserStats] the user id does not exist %w", err)
	}

	query := `SELECT user_id, total_games, total_wins, total_losses, total_draws, total_time, average_time FROM user_stats WHERE user_id = $1`

	var stats types.UserStats

	err = database.DBPool.QueryRow(ctx, query, id_user).Scan(&stats.ID, &stats.Total_games, &stats.Total_wins, &stats.Total_losses, &stats.Total_draws, &stats.Total_time, &stats.Average_time)
	if err != nil {
		return nil, fmt.Errorf("[GetUserStats] Error retrieving user stats %v : %v", id_user, err)
	}

	return &stats, nil

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
