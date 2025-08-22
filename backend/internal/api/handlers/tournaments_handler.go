package handlers

import (
    "encoding/json"
    "golbugames/backend/internal/sudoku/repository"
    "log"
    "net/http"
    "golbugames/backend/pkg/types"
)

func GetAllTournaments(w http.ResponseWriter, r *http.Request) {
    tournaments, err := repository.GetAllTournaments(r.Context())
    if err != nil {
        log.Printf("%v", err)
        http.Error(w, "Error while retrieving tournaments from the database", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "tournaments": tournaments,
    })
}

func AddTournament(w http.ResponseWriter, r *http.Request) {
    var tournament types.Tournament
    if err := json.NewDecoder(r.Body).Decode(&tournament); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    err := repository.AddTournament(r.Context(), tournament)
    if err != nil {
        log.Printf("Erreur DB : %v", err)
        http.Error(w, "Error while inserting tournament into the database", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message":   "Le tournoi a bien été créé",
    })
}
