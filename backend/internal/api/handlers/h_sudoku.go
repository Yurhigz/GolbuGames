package handlers

import (
	"encoding/json"
	"golbugames/backend/internal/games/sudoku"
	"golbugames/backend/pkg/types"
	"golbugames/backend/pkg/utils"
	"log"
	"net/http"
	"strconv"
)

//  Mettre en place un système de structures personnalisé pour utliser la fonction handle plutôt que handlefunc

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var userReg types.UserRegistration
	err := json.NewDecoder(r.Body).Decode(&userReg)
	if err != nil {
		http.Error(w, "invalid data format", http.StatusBadRequest)
		return
	}
	// revoir les vérifications pour les usernames et passwords
	if userReg.Username == "" || userReg.Password == "" || userReg.Accountname == "" {
		http.Error(w, "username, account name and password are required", http.StatusBadRequest)
		return
	}

	err = sudoku.AddUser(r.Context(), userReg.Username, userReg.Accountname, userReg.Password)
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

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	var userDel types.UserDeletion
	err := json.NewDecoder(r.Body).Decode(&userDel)
	if err != nil {
		http.Error(w, "invalid data format", http.StatusBadRequest)
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

func GetUser(w http.ResponseWriter, r *http.Request) {

	var user *types.User

	strId := r.PathValue("id")

	id, err := strconv.Atoi(strId)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "id must be a number", http.StatusBadRequest)
		return
	}

	user, err = sudoku.GetUser(r.Context(), id)

	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "user id is invalid", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		map[string]string{
			"message":  "User succesfully retrieved",
			"userid":   strconv.Itoa(user.ID),
			"username": user.Username,
		})

}

func UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	var user types.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "invalid data format", http.StatusBadRequest)
		return
	}
	if user.ID <= 0 {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = sudoku.UpdateUserPassword(r.Context(), user.ID, user.Password)
	if err != nil {
		log.Printf("Error updating password for user %d: %v", user.ID, err)
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Password updated successfully",
		"userid":   strconv.Itoa(user.ID),
		"username": user.Username,
	})
}

func AddGrid(w http.ResponseWriter, r *http.Request) {

	var req types.GridRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		log.Printf("Erreur : %v", err)
		http.Error(w, "invalid data format", http.StatusBadRequest)
		return
	}
	difficulty := req.Difficulty
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

	solvedGrid, _ := sudoku.GenerateSolvedGrid()
	savedSolvedGrid := solvedGrid
	playableGrid, _ := sudoku.GeneratePlayableGrid(solvedGrid, difficulty)

	boardStr := utils.GridTransformer(playableGrid)
	solutionStr := utils.GridTransformer(savedSolvedGrid)

	err = sudoku.AddGrid(r.Context(), boardStr, solutionStr, difficulty)
	if err != nil {
		log.Printf("Error saving grid to DB: %v", err)
		http.Error(w, "Failed to save grid", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		map[string]string{
			"message":    "Grid successfully created",
			"board":      boardStr,
			"solution":   solutionStr,
			"difficulty": difficulty,
		})

}

func GetGrid(w http.ResponseWriter, r *http.Request) {

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

	sudokuGrid, err := sudoku.GetRandomGrid(r.Context(), difficulty)

	if err != nil {
		http.Error(w, "Internal retrieval error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Grid sucessfully retrieved",
		"board":   sudokuGrid.Board,
	})
}

func SubmitSoloGame(w http.ResponseWriter, r *http.Request) {
	var game types.Game
	err := json.NewDecoder(r.Body).Decode(&game)
	if err != nil {
		http.Error(w, "invalid data format", http.StatusBadRequest)
		return
	}
	if game.GameMode == "1v1" || game.OpponentID != nil || game.Results != nil {
		http.Error(w, "Invalid game mode or results for solo game", http.StatusBadRequest)
		return
	}

	err = sudoku.SubmitSoloGame(r.Context(), game.UserID, game.Completion_time)
	if err != nil {
		log.Printf("Error submitting game for user %d: %v", game.UserID, err)
		http.Error(w, "Failed to submit game", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Game submitted successfully",
		"userId":  strconv.Itoa(game.UserID),
		"score":   strconv.Itoa(game.Completion_time),
	})

}

func SubmitMultiGame(w http.ResponseWriter, r *http.Request) {
	// Enregistre le score final
	// Calcule le temps
	// Met à jour les statistiques dans la BDD
	// Met à jour le classement
}

func GetUserStats(w http.ResponseWriter, r *http.Request) {
	// Nombre de parties jouées
	// Temps moyen
	// Taux de réussite
	// Niveau préféré
}

func GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	// Classement des meilleurs joueurs
	// Filtrage par difficulté et ELO
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Modification du mot de passe
	// Mise à jour des préférences
}

func GetUserHistory(w http.ResponseWriter, r *http.Request) {
	// Historique des parties
	// Progression
}

func SaveGameProgress(w http.ResponseWriter, r *http.Request) {
	// Sauvegarde l'état actuel
	// Permet de reprendre plus tard
}
