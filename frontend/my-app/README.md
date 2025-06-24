# Première installation
docker-compose build

# Démarrer le conteneur avec un shell
docker-compose up -d

# Se connecter au shell
docker exec -it react_dev_env sh

# Dans le shell du conteneur
npm install  # Si nécessaire
npm start    # Pour lancer l'app

# Pour arrêter le conteneur
docker-compose down

# Dépendances à installer 
```
npm install react-icons
```