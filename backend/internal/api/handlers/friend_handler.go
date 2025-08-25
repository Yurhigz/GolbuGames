package handlers

import (
    "encoding/json"
    "golbugames/backend/internal/sudoku/repository"
    "log"
    "net/http"
    "strconv"
    "golbugames/backend/pkg/types"
    "golbugames/backend/internal/api/middleware"
)

func GetUserFriends(w http.ResponseWriter, r *http.Request) {

    claims, ok := r.Context().Value("claims").(*middleware.CustomClaims)
    if !ok || claims == nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    log.Printf("UserID from token: %s, expires at: %v", claims.UserID, claims.ExpiresAt.Time)

    strId := r.PathValue("id")
    id, err := strconv.Atoi(strId)
    if err != nil {
        log.Printf("%v", err)
        http.Error(w, "id must be a number", http.StatusBadRequest)
        return
    }

    friends, err := repository.GetUserFriends(r.Context(), id)
    if err != nil {
        log.Printf("%v", err)
        http.Error(w, "Error while retrieving friends from the database", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "friends": friends,
    })
}

func RemoveFriend(w http.ResponseWriter, r *http.Request) {
    strId := r.PathValue("id")
    id, err := strconv.Atoi(strId)
    if err != nil {
        log.Printf("%v", err)
        http.Error(w, "id must be a number", http.StatusBadRequest)
        return
    }

    strFId := r.PathValue("f_id")
    f_id, err := strconv.Atoi(strFId)
    if err != nil {
        log.Printf("%v", err)
        http.Error(w, "f_id must be a number", http.StatusBadRequest)
        return
    }

    err = repository.RemoveFriend(r.Context(), id, f_id)
    if err != nil {
        log.Printf("%v", err)
        http.Error(w, "Error while removing friends from the database", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Friend remove avec succèss",
	})
}

func AddFriend(w http.ResponseWriter, r *http.Request) {
    var req types.AddFriendRequest

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    friend, err := repository.GetUserIdDB(r.Context(), req.FriendUsername, req.FriendUsername)
    if err != nil {
        http.Error(w, "Friend not found", http.StatusNotFound)
        return
    }

    if err := repository.AddFriend(r.Context(), req.UserID, friend.ID); err != nil {
        log.Printf("%v", err)
        http.Error(w, "Error while adding friend", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "Friend ajouté avec succès",
        "friend":  friend,
    })
}
