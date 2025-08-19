# GolbuGames - Plateforme de Jeux en Ligne

Une plateforme web moderne permettant de jouer à différents jeux en multijoueur avec un système de classement intégré.

## 🎮 Fonctionnalités

### Actuelles
- Algorithme de génération & résolution d'Algorithme
- Choix du niveau de difficulté
- Initialisation de la génération des grilles
- Base de données PostgreSQL & init.db
- Fonction d'intéraction et de stockage dans la BDD
- Hashage et gestion des mots de passe
- Quelques tests unitaires

### À venir
- Interface web responsive
- Système d'authentification sécurisé
- Mode de jeu solo
- Websocket pour multijoueurs
- Système de classement global
- Nouveaux jeux (tango autres...)
- Chat en temps réel
- Profils utilisateurs personnalisés
- Suivi des performances joueurs

## 🛠 Stack Technique

### Backend
- Go
- PostgreSQL
- Authentication JWT
- Package pgx

### Frontend
- React.js
- TypeScript
- CSS Modules
- Socket.io (pour le temps réel)

## 📦 Installation

1. Cloner le repository
```bash
git clone golbugames
cd GolbuGames
```

2. Installer les dépendances 
```bash
script bash WIP
```

## 🔧 Configuration

### Variables d'environnement
Créer un fichier `.env` à la racine du projet :
```env
DB_HOST=localhost
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=golbugames
JWT_SECRET=your_jwt_secret
```

## 🎲 Jeux disponibles (En développement)

### Sudoku
- Mode multijoueur en temps réel
- Système de score basé sur le temps et la difficulté
- Plusieurs niveaux de difficulté

## 📈 Système de Classement

### En cours de développement
- Classement global
- Classement par jeu
- Système de points et de niveaux
- Badges et récompenses

## 🔐 Authentification

- Inscription/Connexion sécurisée
- Gestion des sessions avec JWT
- Récupération de mot de passe
- Profils utilisateurs

## 🚀 Roadmap

### Phase 1
- [x] Mise en place de la logique métier du sudoku
- [x] Mise en place de la connexion avec la DB et fonctions d'intéractions
- [x] Implémentation du Sudoku et intéractions BDD 
- [x] Implémentation graphique et intéraction UI
- [ ] Système d'authentification JWT tokens
      
### Phase 2
- [ ] Ajout du chat en temps réel
- [ ] Mise en place mode multijoueurs & websocket
- [ ] Système de badges
- [x] Système de classement basique

### Phase 3
- [x] Mode tournoi
- [x] Système d'amis
- [ ] Profils utilisateurs avancés & statistiques
- [ ] Application mobile

## 👥 Contribution

cf CONTRIBUTING.md

## 📝 Licence

Distributed under the MIT License. See [`LICENSE`](LICENSE) for more information.

---

*Note: Ce README sera mis à jour régulièrement avec l'évolution du projet.*

