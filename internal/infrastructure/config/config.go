package config

import (
	"database/sql"
	"log"
	"os"

	"github.com/brycedarling/go-practical-microservices/internal/domain/identity"
	"github.com/brycedarling/go-practical-microservices/internal/domain/viewing"
	"github.com/brycedarling/go-practical-microservices/internal/eventstore"
	"github.com/gomodule/redigo/redis"
)

// Config ...
type Config struct {
	Env            *Env
	DB             *sql.DB
	EventStore     eventstore.Store
	RedisPool      *redis.Pool
	IdentityQuery  identity.Query
	ViewingQuery   viewing.Query
	PasswordHasher identity.PasswordHasher
}

// NewConfig ...
func NewConfig(
	env *Env,
	db *sql.DB,
	store eventstore.Store,
	pool *redis.Pool,
	iq identity.Query,
	vq viewing.Query,
	ph identity.PasswordHasher,
) *Config {
	return &Config{env, db, store, pool, iq, vq, ph}
}

// InitializeTestEnvConfig ...
func InitializeTestEnvConfig() *Config {
	os.Setenv("APP_ENV", "test")
	os.Setenv("PORT", "8888")
	os.Setenv("EVENT_STORE_CONNECTION_STRING", "dbname=micro user=message_store password=postgres")
	os.Setenv("QUERY_CONNECTION_STRING", ":6379")

	env, err := InitializeEnv()
	if err != nil {
		log.Fatal(err)
	}
	conf, _, err := InitializeConfig(env)
	if err != nil {
		log.Fatal(err)
	}
	return conf
}
