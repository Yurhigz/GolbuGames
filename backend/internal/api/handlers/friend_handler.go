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

func GetUserFriends(w http.ResponseWriter, r *http.Request) {
	strId := r.PathValue("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		slog.Error("invalid user id", "user_id", strId, "error", err)
		api_errors.WriteError(w, api_errors.ErrInvalidUserID)
		return
	}

	friends, err := repository.GetUserFriends(r.Context(), id)
	if err != nil {
		slog.Error("failed to retrieve friends for user", "user_id", strId, "error", err)
		api_errors.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(map[string]interface{}{
		"friends": friends,
	}); err != nil {
		slog.Error("failed to encode friends for user", "user_id", strId, "error", err)
		api_errors.WriteError(w, err)
		return
	}
}

func RemoveFriend(w http.ResponseWriter, r *http.Request) {
	strId := r.PathValue("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		slog.Error("invalid user id", "user_id", strId, "error", err)
		api_errors.WriteError(w, api_errors.ErrInvalidUserID)
		return
	}

	strFriendId := r.PathValue("f_id")
	friend_id, err := strconv.Atoi(strFriendId)
	if err != nil {
		slog.Error("invalid friend id", "friend_id", strFriendId, "error", err)
		api_errors.WriteError(w, api_errors.ErrInvalidFriendID)
		return
	}

	err = repository.RemoveFriend(r.Context(), id, friend_id)
	if err != nil {
		slog.Error("failed to remove friend for user", "user_id", strId, "error", err)
		api_errors.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(map[string]string{
		"message": "Friend remove avec succèss",
	}); err != nil {
		slog.Error("failed to encode response for user", "user_id", strId, "error", err)
		api_errors.WriteError(w, err)
		return
	}
}

func AddFriend(w http.ResponseWriter, r *http.Request) {
	var req types.AddFriendRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("failed to decode add friend request", "error", err)
		api_errors.WriteError(w, api_errors.ErrBadRequest)
		return
	}

	friend, err := repository.GetUserIdDB(r.Context(), req.FriendUsername, req.FriendUsername)
	if err != nil {
		slog.Error("failed to retrieve friend by username or email", "username_or_email", req.FriendUsername, "error", err)
		api_errors.WriteError(w, api_errors.ErrFriendNotFound)
		return
	}

	if err := repository.AddFriend(r.Context(), req.UserID, friend.ID); err != nil {
		slog.Error("failed to add friend", "user_id", req.UserID, "friend_id", friend.ID, "error", err)
		api_errors.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Friend ajouté avec succès",
		"friend":  friend,
	}); err != nil {
		slog.Error("failed to encode response for user", "user_id", req.UserID, "error", err)
		api_errors.WriteError(w, err)
		return
	}
}
