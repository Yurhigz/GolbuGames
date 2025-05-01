### Étapes pour mieux maîtriser les JWT et les claims

1. **Commencez par l'implémentation basique** :
   - Créez une simple fonction de login générant un JWT
   - Implémentez un middleware de vérification minimal
   - Testez avec une route protégée simple

2. **Expérimentez avec les outils** :
   - Utilisez [jwt.io](https://jwt.io/) pour examiner visuellement vos tokens
   - Observez comment le contenu des claims apparaît dans la partie décodée

3. **Progressez par étapes** :
   - Ajoutez un claim personnalisé simple (ex: username)
   - Puis ajoutez des roles simples
   - Testez différentes vérifications

4. **Développez un cas d'utilisation complet** :
   - Implémentez le flux d'authentification entier
   - Ajoutez les vérifications d'autorisation

### Points à retenir pour la compréhension

- Le JWT est comme une "carte d'identité numérique" signée
- Les claims sont les informations sur cette carte d'identité
- La signature garantit que personne n'a modifié la carte
- On vérifie cette carte à chaque accès à une zone restreinte

### Exemple de progression pratique

```go
// ÉTAPE 1: Authentification basique avec JWT
func LoginHandler(w http.ResponseWriter, r *http.Request) {
    // Authentification simplifiée (à compléter avec votre logique)
    username := r.FormValue("username")
    password := r.FormValue("password")
    
    // Générer un token simple
    token, err := GenerateJWT(username)
    if err != nil {
        http.Error(w, "Erreur de génération du token", http.StatusInternalServerError)
        return
    }
    
    // Retourner le token
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// ÉTAPE 2: Ajout d'un middleware simple
func SimpleAuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tokenString := r.Header.Get("Authorization")
        
        // Vérification minimale
        err := VerifyToken(tokenString)
        if err != nil {
            http.Error(w, "Non autorisé", http.StatusUnauthorized)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}

// ÉTAPE 3: Extraire et utiliser un claim simple
func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
    tokenString := r.Header.Get("Authorization")
    
    // Extraire un claim simple
    username, err := ExtractUsernameFromToken(tokenString)
    if err != nil {
        http.Error(w, "Erreur d'extraction", http.StatusInternalServerError)
        return
    }
    
    // Utiliser le claim
    fmt.Fprintf(w, "Bonjour, %s !", username)
}

// Fonction auxiliaire pour extraire un claim
func ExtractUsernameFromToken(tokenString string) (string, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return secretKey, nil
    })
    
    if err != nil {
        return "", err
    }
    
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        return claims["username"].(string), nil
    }
    
    return "", fmt.Errorf("claims invalides")
}
```