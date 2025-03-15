package database

import (
	"context"
	"log"
	"time"
)

func AddUser(parentsContext context.Context, username, password string) error {
	ctx, cancel := context.WithTimeout(parentsContext, 2*time.Second)
	defer cancel()
	// trouver un moyen pour éviter les injections sql basiques
	query := `INSERT INTO users (username, password) VALUES ($1, $2)`

	_, err := DBPool.Exec(ctx, query, username, password)
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

	_, err := DBPool.Exec(ctx, query, id_user)

	if err != nil {
		log.Printf("Error deleting user (ID: %d): %v", id_user, err)
		return err
	}

	log.Printf("User (ID: %d) deleted successfully:", id_user)
	return nil
}

//  Mettre en place une fonction de stockage des grilles complètes en version string pour simplifier le stockage

// func AddGrid(parentsContext context.Context) {

// }

// Une fonction de récupération des grilles pour l'utilisateur côté frontend

// func GetGrid() {

// }
