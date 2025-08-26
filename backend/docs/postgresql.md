# Guide PostgreSQL

## 1. Se connecter à la base depuis Docker

Pour ouvrir un shell dans le conteneur PostgreSQL :

```bash
docker-compose exec -it -u postgres postgres /bin/bash
```

Puis, pour se connecter à PostgreSQL :

```bash
psql -U postgres
```

## 2. Commandes utiles dans psql

### Lister les bases de données

```sql
\l
```

### Se connecter à une base spécifique

```sql
\c golbugamesdb
```

### Lister les tables

```sql
\dt
```

### Afficher la structure d'une table

```sql
\d sudoku_games
```

### Exécuter des requêtes simples

* Compter le nombre de grilles par difficulté :

```sql
SELECT difficulty, COUNT(*) FROM sudoku_games GROUP BY difficulty;
```

* Afficher toutes les grilles d’une difficulté spécifique :

```sql
SELECT * FROM sudoku_games WHERE difficulty = 'hard';
```

* Afficher toutes les grilles :

```sql
SELECT * FROM sudoku_games;
```

### Quelques commandes supplémentaires utiles

* Lister les utilisateurs PostgreSQL :

```sql
\du
```

* Quitter psql :

```sql
\q
```

* Exécuter un script SQL depuis le conteneur :

```bash
psql -U postgres -d golbugamesdb -f /chemin/vers/script.sql
```

* Vérifier la version de PostgreSQL :

```sql
SELECT version();
```

---

