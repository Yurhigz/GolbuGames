package handlers

import (
	"encoding/json"
	"golbugames/backend/internal/games/sudoku"
	"golbugames/backend/pkg/types"
	"log"
	"net/http"
	"strconv"
)

func createUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "unauthorized method", http.StatusMethodNotAllowed)
		return
	}

	var userReg types.UserRegistration
	err := json.NewDecoder(r.Body).Decode(&userReg)
	if err != nil {
		http.Error(w, "invalid data format, must be a json", http.StatusBadRequest)
		return
	}
	// revoir les vérifications pour les usernames et passwords
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Utilisateur créé avec succès",
		"username": userReg.Username,
	})

}

func deleterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "unauthorized method", http.StatusMethodNotAllowed)
		return
	}

	var userDel types.UserDeletion
	err := json.NewDecoder(r.Body).Decode(&userDel)
	if err != nil {
		http.Error(w, "invalid data format, must be a json", http.StatusBadRequest)
		return
	}

	if userDel.ID <= 0 {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = sudoku.DeleteUser(r.Context(), userDel.ID)

	if err != nil {
		log.Printf("Error deleting user %d: %v", userDel.ID, err)
		http.Error(w, "the user cannot be deleted", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User sucessfully deleted",
		"userId":  strconv.Itoa(userDel.ID),
	})
}

func getUser(w http.ResponseWriter, r *http.Request) {

}

func getGrid(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "unauthorized method", http.StatusMethodNotAllowed)
		return
	}

	difficulty := r.URL.Query().Get("difficulty")
	if difficulty == "" {
		difficulty = "easy"
	}

	validDifficulties := map[string]bool{
		"easy":         true,
		"intermediate": true,
		"advanced":     true,
		"expert":       true,
	}

	if !validDifficulties[difficulty] {
		http.Error(w, "Invalid difficulty level", http.StatusBadRequest)
		return
	}

	board, _, err := sudoku.GetRandomGrid(r.Context(), difficulty)

	if err != nil {
		http.Error(w, "Internal retrieval error", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Grid sucessfully retrieved",
		"board":   board,
	})

}
