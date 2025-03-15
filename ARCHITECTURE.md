# 🌐 Architecture de GolbuGames

Ce document décrit l'architecture technique de l'application.

## 📌 1. Vue d’ensemble
GolbuGames suit une **architecture client-serveur** :
- **Frontend** : React + TypeScript
- **Backend** : Go
- **Communication** : Websockets + REST API
- **Base de données** : PostgreSQL

## 🏗 2. Organisation du backend
Le backend est organisé en plusieurs modules :
```backend/
├── cmd/        # Point d’entrée de l’application (ex: server/main.go)
├── internal/   # Code métier et logique de l’application
│   ├── game/       # Logique du jeu (ex: sudoku.go, game_logic.go)
│   ├── api/        # Handlers REST (ex: handlers.go, routes.go)
│   ├── websocket/  # Gestion des Websockets (ex: client.go, hub.go)
│   ├── database/   # Connexion et gestion de la DB (ex: db.go, models.go)
│   └── config/     # Configuration de l’application (ex: config.go)
└── tests/      # Tests unitaires et d’intégration
    ├── integration/  # Tests d’intégration (ex: game_integration_test.go)
    └── e2e/         # Tests end-to-end (ex: game_e2e_test.go) 
```

## 🔄 3. Flux de données
1. **Un joueur se connecte** via WebSocket (`/ws`)
2. **Un match est créé** (`game/create`)
3. **Chaque joueur remplit la grille de Sudoku**
4. **Le premier à finir gagne !** (`game/win`)

## 🛠 4. Technologies
- Go pour le backend
- React + TypeScript pour le frontend
- PostgreSQL pour la base de données
- Websockets pour la communication en temps réel

## 🛢️ 5. Base de données

Organisation de la base de données

Pour assurer une séparation claire des données et une meilleure scalabilité, chaque module de jeu possède sa propre base de données. Cette approche permet de :

    Isoler les données entre les différents jeux pour éviter tout conflit.
    Faciliter la gestion et l'évolution des schémas sans impacter d'autres jeux.
    Optimiser les performances en répartissant la charge sur plusieurs bases.

Exemple d'architecture des bases de données

    sudoku_db → Stocke les grilles, parties en cours et résultats des joueurs pour le Sudoku.
    chess_db → Stocke les parties, coups et historiques pour les échecs (exemple futur).

Chaque base de données est gérée indépendamment, et les connexions sont ouvertes en fonction du module de jeu en cours. L’application backend sélectionne dynamiquement la base concernée selon le contexte de la requête. Sachant que pour l'ensemble des classements cela sera effectuer par jeu et on pourra ensuite effectuer des moyennes globales pour un classement global.

## 6. Organisation base de données SudokuDB 

Création et stockage de grille de sudoku au préalable ou à l'instant t ? Si création système de complétion en interne pour éviter de redonner les mêmes grilles à un même utilisateur en mode joueur simple. 
Pour le jeu en multijoueurs sélection peu importe les grilles effectuées par les utilisateurs ? 

Fichier config => initialisation de la création des grilles (imaginons 5000 grilles pour commencer), elles sont complètes et selon la difficulté choisie par l'utilisateur on appliquera l'algorithme de suppression d'indice. 
Il faudra ajouter un système de vérification des grilles avant de les stockers en BDD, stockage en string avec un élément de séparation type un "," pour les valeurs d'une même ligne et un ";" pour un deuxième ligne.