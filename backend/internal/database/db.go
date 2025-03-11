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

	// Configurer et créer le pool de connexions
	poolConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return fmt.Errorf("failed to parse DB config: %w", err)
	}

	poolConfig.MaxConns = 20               // Maximum 20 connexions actives
	poolConfig.MinConns = 5                // Minimum 5 connexions actives
	poolConfig.MaxConnLifetime = time.Hour // Ferme une connexion après 1h

	DBPool, err = pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Vérifier la connexion
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
