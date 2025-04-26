package repository

import (
	"context"
	"fmt"
	"golbugames/backend/internal/database"
	"golbugames/backend/pkg/types"
	"log"
	"time"
)

func AddUserDB(parentsContext context.Context, username, accountname, password string) error {
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

func DeleteUserDB(parentsContext context.Context, id_user int) error {
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

func GetUserDB(parentsContext context.Context, id_user int) (*types.User, error) {
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

func UserLoginDB(parentsContext context.Context, username, password string) (bool, error) {
	ctx, cancel := context.WithTimeout(parentsContext, 2*time.Second)
	defer cancel()

	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 AND password = $2)`

	err := database.DBPool.QueryRow(ctx, query, username, password)

	if err != nil {
		return false, fmt.Errorf("[UserLogin] Error checking user %v : %v", username, err)
	}

	return true, nil

}

func GetUserIdDB(parentsContext context.Context, username, accountname string) (*types.User, error) {
	ctx, cancel := context.WithTimeout(parentsContext, 2*time.Second)
	defer cancel()

	var user types.User
	query := `SELECT username,accountname,password,id FROM users WHERE username = $1 AND accountname = $2`

	err := database.DBPool.QueryRow(ctx, query, username, accountname).Scan(&user.Username, &user.Accountname, &user.Password, &user.ID)

	if err != nil {
		return nil, fmt.Errorf("[GetUserId] Error retrieving user %v : %v", username, err)
	}

	return &user, nil
}

func UpdateUserPasswordDB(parentsContext context.Context, id_user int, password string) error {
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

func GetUserStatsDB(parentsContext context.Context, id_user int) (*types.UserStats, error) {
	ctx, cancel := context.WithTimeout(parentsContext, 2*time.Second)
	defer cancel()

	_, err := GetUserDB(parentsContext, id_user)
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
