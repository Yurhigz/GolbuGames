package handlers

import (
	"encoding/json"
	"golbugames/backend/internal/sudoku/repository"
	"golbugames/backend/pkg/types"
	"log"
	"net/http"
	"strconv"
)

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

	err = repository.AddUserDB(r.Context(), userReg.Username, userReg.Accountname, userReg.Password)
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

	strId := r.PathValue("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "id must be a number", http.StatusBadRequest)
		return
	}

	if id <= 0 {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	_, err = repository.GetUserDB(r.Context(), id)
	if err != nil {
		log.Printf("Error retrieving user %d: %v", id, err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	err = repository.DeleteUserDB(r.Context(), id)

	if err != nil {
		log.Printf("Error deleting user %d: %v", id, err)
		http.Error(w, "the user cannot be deleted", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User sucessfully deleted",
		"userId":  strconv.Itoa(id),
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

	user, err = repository.GetUserDB(r.Context(), id)

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

	var pwdUpdate types.PasswordUpdate
	err := json.NewDecoder(r.Body).Decode(&pwdUpdate)
	if err != nil {
		http.Error(w, "invalid data format", http.StatusBadRequest)
		return
	}

	if pwdUpdate.NewPassword == "" {
		http.Error(w, "new password is required", http.StatusBadRequest)
		return
	}

	// Vérifier que l'utilisateur existe
	_, err = repository.GetUserDB(r.Context(), pwdUpdate.ID)
	if err != nil {
		log.Printf("Error retrieving user %d: %v", pwdUpdate.ID, err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	err = repository.UpdateUserPasswordDB(r.Context(), pwdUpdate.ID, pwdUpdate.NewPassword)
	if err != nil {
		log.Printf("Error updating password for user %d: %v", pwdUpdate.ID, err)
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Password updated successfully",
		"userid":  strconv.Itoa(pwdUpdate.ID),
	})
}

func GetUserStats(w http.ResponseWriter, r *http.Request) {
	strId := r.PathValue("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "id must be a number", http.StatusBadRequest)
		return
	}

	_, err = repository.GetUserDB(r.Context(), id)
	if err != nil {
		log.Printf("Error retrieving user %d: %v", id, err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var userStats *types.UserStats

	userStats, err = repository.GetUserStatsDB(r.Context(), id)
	if err != nil {
		log.Printf("Error retrieving stats for user %d: %v", userStats.ID, err)
		http.Error(w, "Failed to retrieve user stats", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		map[string]string{
			"message":      "User stats successfully retrieved",
			"userid":       strconv.Itoa(userStats.ID),
			"total_games":  strconv.Itoa(userStats.Total_games),
			"total_wins":   strconv.Itoa(userStats.Total_wins),
			"total_losses": strconv.Itoa(userStats.Total_losses),
			"total_draws":  strconv.Itoa(userStats.Total_draws),
			"total_time":   strconv.Itoa(userStats.Total_time),
			"average_time": strconv.Itoa(userStats.Average_time),
		})

}
