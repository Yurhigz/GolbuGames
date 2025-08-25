package handlers

import (
// 	"context"
	"encoding/json"
// 	"fmt"
	"golbugames/backend/internal/api/middleware"
	"golbugames/backend/internal/sudoku/repository"
	"golbugames/backend/pkg/types"
	"golbugames/backend/pkg/utils"
	"log"
	"net/http"
	"strconv"
// 	"time"
	"golang.org/x/crypto/bcrypt"
// 	"strings"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Impossible de récupérer l'utilisateur", http.StatusInternalServerError)
		return
	}

	// Génération du JWT d'accès
	jwtToken, err := middleware.GenerateJWT(strconv.Itoa(user.ID), user.Username, []string{"user"})
	if err != nil {
		http.Error(w, "Impossible de générer le token d'accès", http.StatusInternalServerError)
		return
	}

	// Génération du refresh token
	refreshToken, err := middleware.GenerateRefreshToken(strconv.Itoa(user.ID))
	if err != nil {
		http.Error(w, "Impossible de générer le refresh token", http.StatusInternalServerError)
		return
	}

	// Mettre le refresh token dans un cookie HttpOnly
	http.SetCookie(w, &http.Cookie{
        Name:     "refresh_token",
        Value:    refreshToken,
        HttpOnly: true,
        Secure:   true, // pas besoin en local http
        Path:     "/",
        SameSite: http.SameSiteNoneMode, // accepte frontend/backend séparés en local
        MaxAge:   7 * 24 * 60 * 60,
    })



	// Retourner uniquement le JWT dans le corps JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message":      "Utilisateur créé avec succès",
		"access_token": jwtToken,
		"username":     userReg.Username,
		"user_id":      strconv.Itoa(user.ID),
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

    user, err := repository.GetUserIdDB(r.Context(), userRegistration.Username, userRegistration.Accountname)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message":     "User created successfully",
		"username":    userRegistration.Username,
		"accountname": userRegistration.Accountname,
		"user_id":  strconv.Itoa(user.ID),
	})
}

// Compléter avec l'ajout de la génération du token JWT + Token de refraichissement
// sert à vérifier si l'utilisateur existe et à générer un token JWT
// à ajouter dans le handler login

func UserLogin(w http.ResponseWriter, r *http.Request) {
	var user types.User
	var userLogin *types.User

	// Décoder la requête
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "invalid data format", http.StatusBadRequest)
		return
	}

	if user.Username == "" || user.Password == "" || user.Accountname == "" {
		http.Error(w, "username and password are required", http.StatusBadRequest)
		return
	}

	// Récupérer l'utilisateur dans la DB
	userLogin, err = repository.GetUserIdDB(r.Context(), user.Username, user.Accountname)
	if err != nil {
		log.Printf("Error retrieving user %s: %v", user.Username, err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Vérifier le mot de passe
	if !utils.HashPasswordCompare(userLogin.Password, user.Password) {
		log.Printf("Password mismatch for user %s", user.Username)
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Générer le JWT d'accès
	jwtToken, err := middleware.GenerateJWT(strconv.Itoa(userLogin.ID), userLogin.Username, []string{"user"})
	if err != nil {
		http.Error(w, "Impossible de générer le token d'accès", http.StatusInternalServerError)
		return
	}

	// Générer le refresh token
	refreshToken, err := middleware.GenerateRefreshToken(strconv.Itoa(userLogin.ID))
	if err != nil {
		http.Error(w, "Impossible de générer le refresh token", http.StatusInternalServerError)
		return
	}

	// Stocker le refresh token dans un cookie HttpOnly
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,                  // true si HTTPS
		SameSite: http.SameSiteNoneMode,
		MaxAge:   7 * 24 * 60 * 60,      // 7 jours
	})

	// Envoyer uniquement le JWT d'accès côté client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token": jwtToken,
		"user_id":      strconv.Itoa(userLogin.ID),
		"username":     userLogin.Username,
	})
}

func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
    // Récupère le cookie HttpOnly
//     cookie, err := r.Cookie("refreshToken")
    cookie, err := r.Cookie("refresh_token")
    if err != nil {
        http.Error(w, "Refresh token manquant", http.StatusUnauthorized)
        return
    }

    // Vérifie et extrait les claims du refresh token
    claims, err := middleware.VerifyRefreshToken(cookie.Value)
    if err != nil {
        http.Error(w, "Refresh token invalide ou expiré", http.StatusUnauthorized)
        return
    }

    // Génère un nouveau JWT
    newJWT, err := middleware.GenerateJWT(claims.Subject, "", nil)
    if err != nil {
        http.Error(w, "Impossible de générer le nouveau token", http.StatusInternalServerError)
        return
    }

    // Génère un nouveau refresh token pour rotation
    newRefreshToken, err := middleware.GenerateRefreshToken(claims.Subject)
    if err != nil {
        http.Error(w, "Impossible de générer le refresh token", http.StatusInternalServerError)
        return
    }

    // Met à jour le cookie HttpOnly
    http.SetCookie(w, &http.Cookie{
        Name:     "refresh_token",
        Value:    newRefreshToken,
        HttpOnly: true,
        Secure:   true, // pas besoin en local http
        Path:     "/",
        SameSite: http.SameSiteNoneMode, // accepte frontend/backend séparés en local
        MaxAge:   7 * 24 * 60 * 60,
    })


    // Renvoie seulement le nouveau JWT
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "access_token": newJWT,
        "message":      "Nouveau token généré",
    })
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
    // Supprime le cookie refresh_token
    http.SetCookie(w, &http.Cookie{
        Name:     "refresh_token",
        Value:    "",
        Path:     "/",
        HttpOnly: true,
        Secure:   true,                     // à adapter selon ton environnement (false en local HTTP)
        SameSite: http.SameSiteNoneMode,    // correspond à ce que tu utilisais pour le créer
        MaxAge:   -1,                        // supprime le cookie
    })

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Déconnexion réussie",
    })
}
