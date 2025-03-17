package main

import (
	"context"
	"fmt"
	"golbugames/backend/internal/database"
	"golbugames/backend/internal/games"
	"golbugames/backend/tests/unit_tests"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	solvedGrid, err := games.GenerateSolvedGrid()
	if err != nil {
		return
	}

	grid, err := games.GeneratePlayableGrid(solvedGrid, "easy")

	fmt.Println(grid)

	err = database.InitDB(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB() // S'assurer que le pool est fermé à la fin

	unit_tests.CreateUserTest(ctx)
	unit_tests.DeleteUserTest(ctx)

	// Gérer les interruptions système (Ctrl+C, fermeture serveur)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	cancel()
	log.Println("Server shutting down...")

}
