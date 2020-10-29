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
	if err != nil {
		return nil, nil, err
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)

	if err := retryPingPostgres(db); err != nil {
		return nil, nil, err
	}

	return db, func() {
		if err := db.Close(); err != nil {
			log.Fatal(err)
		}
	}, nil
}

const maxRetriesPingPostgres = 3

func retryPingPostgres(db *sql.DB) error {
	for retry := 0; retry <= maxRetriesPingPostgres; retry++ {
		err := pingPostgres(db)
		if err != nil {
			if retry == maxRetriesPingPostgres {
				return err
			}
			log.Println(err)
			retryingIn := time.Duration(5*(retry+1)) * time.Second
			log.Printf("Couldn't ping postgres, retrying in %s", retryingIn)
			time.Sleep(retryingIn)
			continue
		}
		return nil
	}
	return nil
}

func pingPostgres(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return db.PingContext(ctx)
}
