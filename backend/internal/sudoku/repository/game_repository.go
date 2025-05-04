package repository

import (
	"context"
	"fmt"
	"golbugames/backend/internal/database"
	"golbugames/backend/internal/sudoku"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

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

func SubmitSoloGameDB(parentsContext context.Context, userId, completionTime int) error {
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
		return fmt.Errorf("[SubmitSoloGame] cannot submit the game result: %w", err)
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

func SubmitMultiGameDB(parentsContext context.Context, user1, user2 int, result, completionTime int) error {
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
		`INSERT INTO games_scores (user_id, opponent_id, game_mode, result, completion_time) 
		 VALUES ($1, $2, '1v1', $3, $4)`,
		user1, user2, result, completionTime)
	if err != nil {
		return fmt.Errorf("[SubmitMultiGame] cannot submit the game result: %w", err)
	}

	// Mettre à jour les stats des deux joueurs selon le résultat
	isWin1 := result == 0
	isDraw := result == 1
	isWin2 := result == 2

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

func GetEloDB(parentsContext context.Context, userId int) (int, error) {
	ctx, cancel := context.WithTimeout(parentsContext, 2*time.Second)
	defer cancel()

	query := `SELECT elo_score FROM leaderboard WHERE id = $1`

	var eloScore int
	err := database.DBPool.QueryRow(ctx, query, userId).Scan(&eloScore)
	if err != nil {
		return 0, fmt.Errorf("[GetElo] Error retrieving Elo for user (ID: %d): %v", userId, err)
	}

	log.Printf("Elo retrieved successfully for user (ID: %d)", userId)
	return eloScore, nil
}

func UpdateEloDB(parentsContext context.Context, userId1, userId2 int, result string) error {
	ctx, cancel := context.WithTimeout(parentsContext, 2*time.Second)
	defer cancel()

	userElo1, err := GetEloDB(ctx, userId1)
	if err != nil {
		return fmt.Errorf("[UpdateElo] Error retrieving Elo for user (ID: %d): %v", userId1, err)
	}
	userElo2, err := GetEloDB(ctx, userId2)
	if err != nil {
		return fmt.Errorf("[UpdateElo] Error retrieving Elo for user (ID: %d): %v", userId2, err)
	}

	NewElo1, NewElo2 := sudoku.EloCalculation(userElo1, userElo2, result)

	query := `UPDATE leaderboard SET elo_score = $1 WHERE id = $2`

	_, err = database.DBPool.Exec(ctx, query, NewElo1, userId1)
	if err != nil {
		return fmt.Errorf("[UpdateElo] Error updating Elo for user (ID: %d): %v", userId1, err)
	}
	_, err = database.DBPool.Exec(ctx, query, NewElo2, userId2)
	if err != nil {
		return fmt.Errorf("[UpdateElo] Error updating Elo for user (ID: %d): %v", userId1, err)
	}

	log.Printf("Elo updated successfully for both users (ID: %d, %d)", userId1, userId2)
	return nil
}

func GetLeaderboard(parentsContext context.Context) (*[]sudoku.Leaderboard, error) {
	ctx, cancel := context.WithTimeout(parentsContext, 2*time.Second)
	defer cancel()

	var leaderboardList []sudoku.Leaderboard
	query := `SELECT user_id, elo_score, RANK() OVER (ORDER BY elo_score DESC) AS rank FROM leaderboard`

	rows, _ := database.DBPool.Query(ctx, query)
	leaderboardList, err := pgx.CollectRows(rows, pgx.RowToStructByName[sudoku.Leaderboard])
	if err != nil {
		return nil, fmt.Errorf("[GetLeaderboard] Error retrieving leaderboard: %v", err)
	}

	return &leaderboardList, nil
}

func GetUserHistory(parentsContext context.Context, userId int) (*[]sudoku.GameScore, error) {
	ctx, cancel := context.WithTimeout(parentsContext, 2*time.Second)
	defer cancel()

	var gameHistory []sudoku.GameScore
	query := `SELECT id as game_id, user_id, opponent_id, game_mode, results, completion_time FROM games_scores WHERE user_id = $1`

	rows, _ := database.DBPool.Query(ctx, query, userId)
	gameHistory, err := pgx.CollectRows(rows, pgx.RowToStructByName[sudoku.GameScore])
	if err != nil {
		return nil, fmt.Errorf("[GetUserHistory] Error retrieving game history for user (ID: %d): %v", userId, err)
	}

	return &gameHistory, nil
}
