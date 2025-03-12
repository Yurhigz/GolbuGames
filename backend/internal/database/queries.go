package database

import (
	"context"
	"fmt"
	"log"
	"time"
)

func AddUser(parentsContext context.Context, bdd string, id int, username, password string) error {
	ctx, cancel := context.WithTimeout(parentsContext, 2*time.Second)
	defer cancel()

	query := fmt.Sprintf("INSERT INTO %s (username, password) VALUES ($1, $2)", bdd)

	_, err := DBPool.Exec(ctx, query, username, password)
	if err != nil {
		log.Println("Error inserting user:", err)
		return err
	}

	log.Println("✅ User added successfully:", username)
	return nil
}

func DeleteUser() {

}

func AddGrid() {

}

func GetGrid() {

}
