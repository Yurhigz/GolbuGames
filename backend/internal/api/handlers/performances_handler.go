package handlers

import (
	"encoding/json"
	api_errors "golbugames/internal/api/errors"
	"golbugames/internal/sudoku/repository"
	"golbugames/pkg/types"
	"log/slog"
	"net/http"
	"strconv"
)

func GetUserStats(w http.ResponseWriter, r *http.Request) {
	strId := r.PathValue("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		slog.Error("Invalid user ID", "id", strId, "error", err)
		api_errors.WriteError(w, api_errors.ErrInvalidUserID)
		return
	}

	var userStats *types.UserStats

	userStats, err = repository.GetUserStatsDB(r.Context(), id)
	if err != nil {
		slog.Error("Error retrieving stats for user", "id", userStats.ID, "error", err)
		api_errors.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(
		map[string]string{
			"message":      "User stats successfully retrieved",
			"userid":       strconv.Itoa(userStats.ID),
			"total_games":  strconv.Itoa(userStats.Total_games),
			"total_wins":   strconv.Itoa(userStats.Total_wins),
			"total_losses": strconv.Itoa(userStats.Total_losses),
			"total_draws":  strconv.Itoa(userStats.Total_draws),
			"total_time":   strconv.Itoa(userStats.Total_time),
			"average_time": strconv.Itoa(userStats.Average_time),
		}); err != nil {
		slog.Error("Failed to encode response", "error", err)
		api_errors.WriteError(w, err)
		return
	}

}
