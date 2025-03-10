package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgres struct {
	db *pgxpool.Pool
}

func ConnectDB() {
	conn, err := pgx.Connect(context.Background(), "postgresql://host:port/database")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

}

// os.Getenv("DATABASE_URL")
//  https://donchev.is/post/working-with-postgresql-in-go-using-pgx/
