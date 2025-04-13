package main

import (
	"context"
	"golbugames/backend/internal/api/router"
	"golbugames/backend/internal/database"
	"log"
	"net/http"
)

// func timeHandler(w http.ResponseWriter, r *http.Request) {
// 	tm := time.Now().Format(time.RFC1123)
// 	w.Write([]byte("The time is: " + tm + "Au plaisir mon chacalito"))
// }

func main() {
	// -------------- MAIN TESTS --------------
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// solvedGrid, err := games.GenerateSolvedGrid()
	// if err != nil {
	// 	return
	// }

	// grid, err := games.GeneratePlayableGrid(solvedGrid, "easy")

	// fmt.Println(grid)

	err := database.InitDB(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	// defer database.CloseDB() // S'assurer que le pool est fermé à la fin

	// unit_tests.CreateUserTest(ctx)
	// unit_tests.DeleteUserTest(ctx)

	// Gérer les interruptions système (Ctrl+C, fermeture serveur)
	// stop := make(chan os.Signal, 1)
	// signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	// <-stop

	// cancel()
	// log.Println("Server shutting down...")
	// -------------- MAIN TESTS --------------

	// -------------- TEST SERVEUR HTTP --------------
	// type Router struct {
	// 	mux *http.ServeMux
	// }

	router := router.NewRouter()

	router.InitRoutes()

	log.Print("Listening...")
	http.ListenAndServe(":3000", router)

	// -------------- TEST SERVEUR HTTP --------------

}
