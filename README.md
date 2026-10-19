# GolbuGames

GolbuGames est un projet que j'ai commencé pour créer une plateforme de jeux en ligne autour du Sudoku. L'objectif était de proposer une expérience complète, du jeu solo au multijoueur, avec un compte joueur et quelques fonctionnalités pour retrouver une communauté autour du jeu.

Le projet est maintenant terminé dans le périmètre que je m'étais fixé. Je ne prévois pas d'y ajouter de nouveaux jeux ou de nouvelles fonctionnalités. Ce dépôt reste disponible pour garder une trace du projet et de ce que j'ai appris en le construisant.

## Ce que contient le projet

### Sudoku

- Une partie solo avec plusieurs niveaux de difficulté et l'enregistrement des résultats.
- Un mode multijoueur en temps réel avec matchmaking, validation des coups côté serveur et sauvegarde des résultats.
- Des aides de jeu, dont la mise en évidence des cases contenant le même chiffre.
- Deux styles d'interface, rétro et classique, que le joueur peut choisir.

### Comptes et fonctionnalités autour du jeu

- Inscription et connexion avec mots de passe hachés et jetons JWT.
- Profil joueur et changement de mot de passe.
- Liste d'amis et création de tournois.
- Classement et statistiques de base.
- Quelques indicateurs pour suivre l'activité du multijoueur : joueurs en attente, parties actives et parties terminées.

## Technologies

- Frontend : React et CSS.
- Backend : Go, avec une API HTTP et des WebSockets.
- Base de données : PostgreSQL, avec pgx.
- Authentification : JWT.

## Vérifications

Les vérifications effectuées sont la compilation et l'analyse statique du backend, ainsi que le build de production du frontend :

```bash
cd backend
go build ./...
go vet ./...
```

Pour le frontend :

```bash
cd frontend/my-app
npm ci
npm run build
```

Ces commandes ne remplacent pas une validation complète des parcours avec une base de données et plusieurs joueurs connectés. Les tests unitaires et d'intégration du multijoueur n'ont pas été ajoutés, le choix étant de terminer le projet sans prolonger le chantier.

Le projet peut être démarrer sous docker/docker-compose afin de faire fonctionner l'ensemble des sous parties en conteneurs et pouvoir redémarrer cela à convenance.

Le démarrage complet avec une base de données neuve n'a pas été validé. Le script [init-db.sql](./backend/internal/database/init-db.sql) supprime puis recrée la base `golbugamesdb` ; il ne faut pas l'exécuter sur une base contenant des données à conserver. Son schéma doit aussi être rapproché des modèles utilisés par le backend avant un déploiement.

## Fin du projet

Le projet est clôturé à son état actuel. Il remplit son objectif principal : proposer un Sudoku jouable en solo et en multijoueur, avec les fonctionnalités de compte et de suivi prévues à une échelle simple. Certaines parties restent volontairement basiques et le projet n'a pas vocation à devenir une plateforme de jeux plus large.

L'idée est maintenant de le laisser ainsi, d'en garder une trace.

## Licence

Ce projet est distribué sous licence MIT. Voir [LICENSE](./LICENSE).
