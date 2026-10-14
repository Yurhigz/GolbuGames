package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DBPool *pgxpool.Pool

type PoolSettings struct {
	URL                  string
	MinConnections       int32
	MaxConnections       int32
	MaxConnectionAge     time.Duration
	MaxConnectionIdle    time.Duration
	ConnectTimeout       time.Duration
	QueryTimeout         time.Duration
	BulkOperationTimeout time.Duration
}

var (
	QueryTimeout         = 2 * time.Second
	BulkOperationTimeout = 10 * time.Second
)

func InitDB(ctx context.Context, settings PoolSettings) error {
	dbURL := settings.URL
	if dbURL == "" {
		return fmt.Errorf("DATABASE_URL is not set")
	}

	// Configuration du pool à partir de la dburl
	poolConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return fmt.Errorf("failed to parse DB config: %w", err)
	}

	poolConfig.MaxConns = settings.MaxConnections
	poolConfig.MinConns = settings.MinConnections
	poolConfig.MaxConnLifetime = settings.MaxConnectionAge
	poolConfig.MaxConnIdleTime = settings.MaxConnectionIdle
	QueryTimeout = settings.QueryTimeout
	BulkOperationTimeout = settings.BulkOperationTimeout

	//  création des pools à l'aide des éléments de configuration précédents
	connectCtx, cancel := context.WithTimeout(ctx, settings.ConnectTimeout)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(connectCtx, poolConfig)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Vérifier la connexion avec la BDD
	err = pool.Ping(connectCtx)
	if err != nil {
		pool.Close()
		return fmt.Errorf("cannot reach database: %w", err)
	}
	DBPool = pool

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
