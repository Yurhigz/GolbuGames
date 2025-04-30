package handlers

import (
	"encoding/json"
	"golbugames/backend/internal/sudoku/repository"
	"golbugames/backend/pkg/types"
	"log"
	"net/http"
	"strconv"
)

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
