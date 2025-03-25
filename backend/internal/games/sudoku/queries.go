package sudoku

import (
	"context"
	"fmt"
	"golbugames/backend/internal/database"
	"golbugames/backend/pkg/types"
	"log"
	"math/rand/v2"
	"time"
)

func AddUser(parentsContext context.Context, username, password string) error {
	ctx, cancel := context.WithTimeout(parentsContext, 2*time.Second)
	defer cancel()

	// vérification de l'unicité de l'utilisateur
	var exists bool
	err := database.DBPool.QueryRow(parentsContext, "SELECT * FROM users WHERE username = $1 ", username).Scan(&exists)
	if err != nil {
		log.Printf("error checking user existence: %v", err)
		return err
	}
	if exists {
		log.Printf("username %v already exists", username)
		return fmt.Errorf("username already exists")
	}
	// trouver un moyen pour éviter les injections sql basiques
	query := `INSERT INTO users (username, password) VALUES ($1, $2)`

	_, err = database.DBPool.Exec(ctx, query, username, password)
	if err != nil {
		log.Printf("Error inserting user [%s] %v", username, err)
		return err
	}

	log.Println("User added successfully:", username)
	return nil
}

func DeleteUser(parentsContext context.Context, id_user int) error {
	ctx, cancel := context.WithTimeout(parentsContext, 2*time.Second)
	defer cancel()

	query := `DELETE FROM users WHERE id = $1`

	_, err := database.DBPool.Exec(ctx, query, id_user)

	if err != nil {
		log.Printf("Error deleting user (ID: %d): %v", id_user, err)
		return err
	}

	log.Printf("User (ID: %d) deleted successfully:", id_user)
	return nil
}

func GetUser(parentsContext context.Context, id_user int) (*types.User, error) {
	ctx, cancel := context.WithTimeout(parentsContext, 2*time.Second)
	defer cancel()

	var user types.User
	query := `SELECT * FROM users WHERE id = $1`

	err := database.DBPool.QueryRow(ctx, query, id_user).Scan(&user)

	if err != nil {
		log.Printf("Error retrieving user %d : %v", id_user, err)
		return nil, err
	}

	return &user, nil
}

//  Mettre en place une fonction de stockage des grilles complètes en version string pour simplifier le stockage

func AddGrid(parentsContext context.Context, board, solution, difficulty string) error {
	ctx, cancel := context.WithTimeout(parentsContext, 2*time.Second)
	defer cancel()

	query := `INSERT INTO sudoku_games (board, solution, difficulty) VALUES ($1, $2, $3)`

	_, err := database.DBPool.Exec(ctx, query, board, solution, difficulty)

	if err != nil {
		log.Printf("Error inserting sudoku grid : %v", err)
		return err
	}

	log.Printf("Sudoku grid added successfully with difficulty : %s", difficulty)
	return nil
}

// Une fonction de récupération des grilles pour l'utilisateur côté frontend

func GetGrid(parentsContext context.Context, id int) (string, string, error) {
	ctx, cancel := context.WithTimeout(parentsContext, 2*time.Second)
	defer cancel()

	query := `SELECT board,solution FROM sudoku_games WHERE id = $1`
	var board, solution string

	err := database.DBPool.QueryRow(ctx, query, id).Scan(&board, &solution)
	if err != nil {
		log.Printf("cannot find grid %d : %v", id, err)
		return "", "", err
	}

	log.Printf("Sudoku grid (id: %d) successfully retrieved", id)
	return board, solution, nil
}

// Une fonction de récupération aléatoire des grilles pour l'utilisateur a voir si on choisit aléatoirement côté back ou côté front selon les performances

func GetRandomGrid(parentsContext context.Context, difficulty string) (string, string, error) {
	ctx, cancel := context.WithTimeout(parentsContext, 2*time.Second)
	defer cancel()

	// Première requête pour compter
	var count int
	countQuery := `SELECT COUNT(*) FROM sudoku_games WHERE difficulty = $1`
	err := database.DBPool.QueryRow(ctx, countQuery, difficulty).Scan(&count)
	if err != nil {
		log.Printf("failed to count grids: %v", err)
		return "", "", err
	}

	if count == 0 {
		log.Printf("no grids found for difficulty %s", difficulty)
		return "", "", err
	}

	offset := rand.IntN(count)

	var board, solution string
	query := `
        SELECT board, solution 
        FROM sudoku_games 
        WHERE difficulty = $1 
        LIMIT 1 OFFSET $2`

	err = database.DBPool.QueryRow(ctx, query, difficulty, offset).Scan(&board, &solution)
	if err != nil {
		return "", "", fmt.Errorf("failed to get random grid: %w", err)
	}

	return board, solution, nil
}
