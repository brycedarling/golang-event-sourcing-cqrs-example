package viewing

import (
	"github.com/brycedarling/go-practical-microservices/internal/domain/viewing"
	"github.com/gomodule/redigo/redis"
)

// NewQueryRedis ...
func NewQueryRedis(pool *redis.Pool) viewing.Query {
	return &queryRedis{pool}
}

type queryRedis struct {
	pool *redis.Pool
}

var _ viewing.Query = (*queryRedis)(nil)

const (
	viewingKey           string = "viewing"
	videosWatchedKey     string = "videos_watched"
	lastViewProcessedKey string = "last_view_processed"
)

// Initialize ...
func (q *queryRedis) Initialize() error {
	conn := q.pool.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", viewingKey))
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	v := viewing.NewViewing()
	_, err = conn.Do("HMSET", viewingKey,
		videosWatchedKey, v.VideosWatched,
		lastViewProcessedKey, v.LastViewProcessed,
	)
	if err != nil {
		return err
	}
	return nil
}

// Find ...
func (q *queryRedis) Find() (*viewing.Viewing, error) {
	conn := q.pool.Get()
	defer conn.Close()

	values, err := redis.Values(conn.Do("HGETALL", viewingKey))
	if err != nil {
		return nil, err
	} else if len(values) == 0 {
		return nil, viewing.ErrViewingNotFound
	}

	var v viewing.Viewing
	err = redis.ScanStruct(values, &v)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

// IncrementVideosWatched ...
func (q *queryRedis) IncrementVideosWatched(globalPosition int) error {
	conn := q.pool.Get()
	defer conn.Close()

	reply, err := redis.Ints(conn.Do("HMGET", viewingKey, lastViewProcessedKey))
	if err != nil {
		return err
	}
	if len(reply) != 1 {
		return nil
	}
	lastViewProcessed := reply[0]
	if lastViewProcessed >= globalPosition {
		return nil
	}
	err = conn.Send("HINCRBY", viewingKey, videosWatchedKey, 1)
	if err != nil {
		return err
	}
	err = conn.Send("HMSET", viewingKey, lastViewProcessedKey, globalPosition)
	if err != nil {
		return err
	}
	if err := conn.Flush(); err != nil {
		return err
	}
	/*
		rval, err := conn.Receive()
		if err != nil {
			return err
		}
	*/
	return nil
}
