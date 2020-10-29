package config

import (
	"context"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

// NewRedisPool ...
func NewRedisPool(env *Env) (*redis.Pool, func(), error) {
	pool := &redis.Pool{
		// MaxIdle:   50,
		// MaxActive: 10000,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", env.QueryConnectionString)
		},
	}

	if err := retryPingRedis(pool); err != nil {
		return nil, nil, err
	}

	return pool, func() {
		if err := pool.Close(); err != nil {
			log.Fatal(err)
		}
	}, nil
}

const maxRetriesPingRedis = 3

func retryPingRedis(pool *redis.Pool) error {
	for retry := 0; retry <= maxRetriesPingRedis; retry++ {
		err := pingRedis(pool)
		if err != nil {
			if retry == maxRetriesPingRedis {
				return err
			}
			retryingIn := time.Duration(5*(retry+1)) * time.Second
			log.Printf("Couldn't ping redis, retrying in %s", retryingIn)
			time.Sleep(retryingIn)
			continue
		}
		return nil
	}
	return nil
}

func pingRedis(pool *redis.Pool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := pool.GetContext(ctx)
	defer conn.Close()
	return err
}
