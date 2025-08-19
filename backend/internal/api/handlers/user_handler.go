package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"golbugames/backend/internal/api/middleware"
	"golbugames/backend/internal/sudoku/repository"
	"golbugames/backend/pkg/types"
	"golbugames/backend/pkg/utils"
	"log"
	"net/http"
	"strconv"
	"time"
	"golang.org/x/crypto/bcrypt"
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

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userReg.Password), bcrypt.DefaultCost)
    if err != nil {
        // gérer l'erreur
        http.Error(w, "Erreur lors du hash du mot de passe", http.StatusInternalServerError)
        return
    }

	err = repository.AddUserDB(r.Context(), userReg.Username, userReg.Accountname, string(hashedPassword))
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

func GetUserId(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	accountname := r.URL.Query().Get("accountname")

	if username == "" || accountname == "" {
		http.Error(w, "username and accountname are required", http.StatusBadRequest)
		return
	}

	var user *types.User

	user, err := repository.GetUserIdDB(r.Context(), username, accountname)

	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "username or accountname are invalid", http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		map[string]string{
			"message": "User succesfully retrieved",
			"userid":  strconv.Itoa(user.ID),
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

// Il faudra intégrer la validation par mail de l'inscription de l'utilisateur.
func UserSignin(w http.ResponseWriter, r *http.Request) {
	var userRegistration types.UserRegistration

	err := json.NewDecoder(r.Body).Decode(&userRegistration)

	if err != nil {
		http.Error(w, "invalid data format", http.StatusBadRequest)
		return
	}

	if userRegistration.Username == "" || userRegistration.Password == "" || userRegistration.Accountname == "" {
		http.Error(w, "username, account name and password are required", http.StatusBadRequest)
		return
	}

	// hashage du mot de passe avant sauvegarde dans la DB

	userRegistration.Password, err = utils.HashPassword(userRegistration.Password)
	if err != nil {
		log.Printf("Error hashing password for user %s: %v", userRegistration.Username, err)
		http.Error(w, "Error while hashing password", http.StatusInternalServerError)
		return
	}

	err = repository.AddUserDB(r.Context(), userRegistration.Username, userRegistration.Accountname, userRegistration.Password)

	if err != nil {
		log.Printf("Error adding user %s: %v", userRegistration.Username, err)
		http.Error(w, "Error while adding a new user to the database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message":     "User created successfully",
		"username":    userRegistration.Username,
		"accountname": userRegistration.Accountname,
	})
}

// Compléter avec l'ajout de la génération du token JWT + Token de refraichissement
// sert à vérifier si l'utilisateur existe et à générer un token JWT
// à ajouter dans le handler login

func UserLogin(w http.ResponseWriter, r *http.Request) {
	var user types.User
	var userLogin *types.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "invalid data format", http.StatusBadRequest)
		return
	}

	if user.Username == "" || user.Password == "" || user.Accountname == "" {
		http.Error(w, "username and password are required", http.StatusBadRequest)
		return
	}

	// récupération de l'utilisateur dans la base de données à l'aide de son nom d'utilisateur et ID
	userLogin, err = repository.GetUserIdDB(r.Context(), user.Username, user.Accountname)
	if err != nil {
		log.Printf("Error retrieving user %s: %v", user.Username, err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	if !utils.HashPasswordCompare(userLogin.Password, user.Password) {
		log.Printf("Password mismatch for user %s", user.Username)
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Gérer la partie refresh token avec la fonction RefreshToken
	jwtToken, err := middleware.GenerateJWT(strconv.Itoa(userLogin.ID), userLogin.Username, []string{"user"})

	w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "access_token": jwtToken,
    })
// 	refreshToken, err := RefreshToken(r.Context(), strconv.Itoa(userLogin.ID))
//
// 	repository.StoreRefreshToken(r.Context(), userLogin.ID, refreshToken)
//
// 	if err != nil {
// 		log.Printf("Error generating JWT for user %s: %v", userLogin.Username, err)
// 		http.Error(w, "Error generating token", http.StatusInternalServerError)
// 		return
// 	}
//
// 	// Code à vérifier pour le token de refraichissement et à corriger côté DB
// 	http.SetCookie(w, &http.Cookie{
// 		Name:     "refresh_token",
// 		Value:    refreshToken,
// 		HttpOnly: true,
// 		Secure:   true,
// 		SameSite: http.SameSiteStrictMode,
// 		Path:     "/refresh",
// 		Expires:  time.Now().Add(30 * 24 * time.Hour),
// 	})
//
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(map[string]string{
// 		"access_token": jwtToken,
// 	})
//
//     w.Header().Set("Content-Type", "application/json")
//     json.NewEncoder(w).Encode(map[string]string{
//         "access_token": jwtToken,
//     })

}

// sert à rafraîchir le token JWT à l'aide du token de rafraîchissement
// mais n'est pas encore implémenté
func RefreshToken(parentsContext context.Context, refreshToken string) (string, error) {
	_, cancel := context.WithTimeout(parentsContext, 2*time.Second)
	defer cancel()

	return "", fmt.Errorf("[RefreshToken] Refresh token functionality is not implemented yet")

}
