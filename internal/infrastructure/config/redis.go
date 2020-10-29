package config

import (
	"log"

	"github.com/gomodule/redigo/redis"
)

// NewRedisPool ...
func NewRedisPool(env *Env) (*redis.Pool, func()) {
	pool := &redis.Pool{
		// MaxIdle:   50,
		// MaxActive: 10000,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", env.QueryConnectionString)
		},
	}
	return pool, func() {
		if err := pool.Close(); err != nil {
			log.Fatal(err)
		}
	}
}
