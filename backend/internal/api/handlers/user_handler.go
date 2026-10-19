package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"golbugames/internal/api/middleware"
	"golbugames/internal/sudoku/repository"
	"golbugames/pkg/types"
	"golbugames/pkg/utils"
	"log"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	createUser(w, r, nil)
}

func CreateUserWithJWT(manager *middleware.JWTManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		createUser(w, r, manager)
	}
}

func createUser(w http.ResponseWriter, r *http.Request, manager *middleware.JWTManager) {

	var userReg types.UserRegistration
	err := json.NewDecoder(r.Body).Decode(&userReg)
	if err != nil {
		http.Error(w, "invalid data format", http.StatusBadRequest)
		return
	}
	if userReg.Username == "" || userReg.Password == "" || userReg.Accountname == "" {
		http.Error(w, "username, account name and password are required", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userReg.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Erreur lors du hash du mot de passe", http.StatusInternalServerError)
		return
	}

	err = repository.AddUserDB(r.Context(), userReg.Username, userReg.Accountname, string(hashedPassword))
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "Error while adding a new user to the database", http.StatusInternalServerError)
		return
	}

	user, err := repository.GetUserIdDB(r.Context(), userReg.Username, userReg.Accountname)
	if err != nil {
		http.Error(w, "user was created but could not be retrieved", http.StatusInternalServerError)
		return
	}
	if manager == nil {
		http.Error(w, "JWT signer is not configured", http.StatusInternalServerError)
		return
	}

	accessToken, err := manager.Generate(strconv.Itoa(user.ID), user.Username, []string{"user"})
	if err != nil {
		http.Error(w, "could not create access token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message":      "Utilisateur créé avec succès",
		"username":     userReg.Username,
		"user_id":      strconv.Itoa(user.ID),
		"access_token": accessToken,
	})

}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		http.Error(w, "authentication required", http.StatusUnauthorized)
		return
	}

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
	if strconv.Itoa(id) != claims.UserID {
		http.Error(w, "cannot delete another user's account", http.StatusForbidden)
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
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		http.Error(w, "authentication required", http.StatusUnauthorized)
		return
	}

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
	pwdUpdate.ID, err = strconv.Atoi(claims.UserID)
	if err != nil {
		http.Error(w, "invalid authenticated user ID", http.StatusUnauthorized)
		return
	}

	// Vérifier que l'utilisateur existe
	_, err = repository.GetUserDB(r.Context(), pwdUpdate.ID)
	if err != nil {
		log.Printf("Error retrieving user %d: %v", pwdUpdate.ID, err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	pwdUpdate.NewPassword, err = utils.HashPassword(pwdUpdate.NewPassword)
	if err != nil {
		http.Error(w, "Error while hashing password", http.StatusInternalServerError)
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

	user, err := repository.GetUserIdDB(r.Context(), userRegistration.Username, userRegistration.Accountname)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message":     "User created successfully",
		"username":    userRegistration.Username,
		"accountname": userRegistration.Accountname,
		"user_id":     strconv.Itoa(user.ID),
	})
}

func UserLogin(w http.ResponseWriter, r *http.Request) {
	userLogin(w, r, nil)
}

func UserLoginWithJWT(manager *middleware.JWTManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userLogin(w, r, manager)
	}
}

func userLogin(w http.ResponseWriter, r *http.Request, manager *middleware.JWTManager) {
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

	userLogin, err = repository.GetUserIdDB(r.Context(), user.Username, user.Accountname)
	if err != nil {
		log.Printf("Error retrieving user %s: %v", user.Username, err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(userLogin.Password), []byte(user.Password)); err != nil {
		log.Printf("Password mismatch for user %s", user.Username)
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}
	if manager == nil {
		http.Error(w, "JWT signer is not configured", http.StatusInternalServerError)
		return
	}

	jwtToken, err := manager.Generate(strconv.Itoa(userLogin.ID), userLogin.Username, []string{"user"})
	if err != nil {
		http.Error(w, "could not create access token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"access_token": jwtToken,
		"user_id":      strconv.Itoa(userLogin.ID),
		"username":     userLogin.Username,
	})
}

func RefreshToken(parentsContext context.Context, refreshToken string) (string, error) {
	_, cancel := context.WithTimeout(parentsContext, 2*time.Second)
	defer cancel()

	return "", fmt.Errorf("[RefreshToken] Refresh token functionality is not implemented yet")

}
