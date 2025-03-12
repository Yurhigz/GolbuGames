package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DBPool *pgxpool.Pool

func InitDB() error {
	dbURL := "host=127.0.0.1 user=postgres password=postgres dbname=sudokudb port=5442 sslmode=disable TimeZone=Asia/Shanghai"
	if dbURL == "" {
		return fmt.Errorf("DATABASE_URL is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Configuration du pool à partir de la dburl
	poolConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return fmt.Errorf("failed to parse DB config: %w", err)
	}

	poolConfig.MaxConns = 20                      // Maximum 20 connexions actives
	poolConfig.MinConns = 5                       // Minimum 5 connexions actives
	poolConfig.MaxConnLifetime = time.Hour        // Ferme une connexion après 1h
	poolConfig.MaxConnIdleTime = 10 * time.Minute // limite du temps d'inactivité pour un pool

	//  création des pools à l'aide des éléments de configuration précédents
	DBPool, err = pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Vérifier la connexion avec la BDD
	err = DBPool.Ping(ctx)
	if err != nil {
		return fmt.Errorf("cannot reach database: %w", err)
	}

	log.Println("Database pool initialized successfully")
	return nil
}

// CloseDB ferme le pool de connexions à la base de données
func CloseDB() {
	if DBPool != nil {
		DBPool.Close()
		log.Println("Database pool closed")
	}
}

// os.Getenv("DATABASE_URL")
//  https://donchev.is/post/working-with-postgresql-in-go-using-pgx/
//  Il faudra mettre en place une gestion des connexions dans le cas de figure où la DB redémarre avec une boucle de tentative de reconnexion
// et éventuellement limiter le nombre de tentatives ou mettre un temps entre chaque tentative
