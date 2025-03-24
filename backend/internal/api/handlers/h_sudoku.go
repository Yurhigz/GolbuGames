package handlers

import (
	"encoding/json"
	"golbugames/backend/internal/games/sudoku"
	"log"
	"net/http"
)

type UserRegistration struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func createUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "unauthorized method", http.StatusMethodNotAllowed)
		return
	}

	var userReg UserRegistration
	err := json.NewDecoder(r.Body).Decode(&userReg)
	if err != nil {
		http.Error(w, "invalid data format, must be a json", http.StatusBadRequest)
		return
	}
	// revoir les vérifications pour les usernames et password
	if userReg.Username == "" || userReg.Password == "" {
		http.Error(w, "username and password are required", http.StatusBadRequest)
		return
	}

	err = sudoku.AddUser(r.Context(), userReg.Username, userReg.Password)
	// Vérifier les duplicatas d'utilisateurs
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "Error while adding a new user to the database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Utilisateur créé avec succès",
		"username": userReg.Username,
	})

}

func deleterUser(w http.ResponseWriter, r *http.Request) {

}

func addGrid(w http.ResponseWriter, r *http.Request) {

}

func getGrid(w http.ResponseWriter, r *http.Request) {

}

func getRandomGrid(w http.ResponseWriter, r *http.Request) {

}
