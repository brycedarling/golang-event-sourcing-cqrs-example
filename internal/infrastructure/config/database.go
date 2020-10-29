package config

import (
	"context"
	"database/sql"
	"log"
	"time"

	// Postgres
	_ "github.com/lib/pq"
)

// NewDB ...
func NewDB(env *Env) (*sql.DB, func(), error) {
	db, err := sql.Open("postgres", env.EventStoreConnectionString)
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)

	err = ping(db)
	if err != nil {
		return nil, nil, err
	}

	return db, func() {
		if err := db.Close(); err != nil {
			log.Fatal(err)
		}
	}, err
}

func ping(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return db.PingContext(ctx)
}
