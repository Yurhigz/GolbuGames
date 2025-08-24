## Se connecter à la DB postgreSQL 

Ouvrir un shell à partir des conteneurs 

```bash
docker-compose exec -it -u postgres postgres /bin/bash
```

Une fois dans le shell, on se connecte à la DB avec 

```bash
psql -U postgres

```

```bash

\l 
\c golbugamesdb
\dt 

```