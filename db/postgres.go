package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose"

	_ "github.com/lib/pq"
)

type Postgres struct {
    DB *sqlx.DB
}

func PostgresConnection(connectionString string) (*Postgres, error) {
    db, err := sqlx.Open("postgres", connectionString)
    if err != nil {
        return nil, fmt.Errorf("error connecting to the database: %w", err)
    }

	// Run migrations
    err = goose.Up(db.DB, "./migrations")
    if err != nil {
        panic(err)
    }

    // Ping DB to check connection is successful
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("error pinging the database: %w", err)
    }

    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(25)

    fmt.Println("Successfully connected to the database!")
    return &Postgres{DB: db}, nil
}