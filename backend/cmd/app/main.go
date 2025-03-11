package main

import (
	"fmt"
	"golbugames/backend/internal/database"
	"golbugames/backend/internal/game"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	grid, err := game.GenerateGrid("easy")

	if err != nil {
		return
	}

	fmt.Println(grid)

	err = database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	defer database.CloseDB() // S'assurer que le pool est fermé à la fin

	// Gérer les interruptions système (Ctrl+C, fermeture serveur)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Server shutting down...")

}
