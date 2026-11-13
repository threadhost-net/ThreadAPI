package platform

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type MultiRedis struct {
	conns           map[string]*redis.Client
	defaultHost     string
	defaultPassword string
	errors          []error
}

func NewMultiRedis() *MultiRedis {
	return &MultiRedis{conns: make(map[string]*redis.Client)}
}

func (r *MultiRedis) WithDefault(host, password string) *MultiRedis {
	r.defaultHost = host
	r.defaultPassword = password
	return r
}

func (r *MultiRedis) Add(alias string, db int) *MultiRedis {
	cli := redis.NewClient(&redis.Options{
		Addr:     r.defaultHost,
		Password: r.defaultPassword,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := cli.Ping(ctx).Err(); err != nil {
		r.errors = append(r.errors, err)
	} else {
		r.conns[alias] = cli
	}
	return r
}

func (r *MultiRedis) Get(alias string) *redis.Client {
	if cli, ok := r.conns[alias]; ok {
		return cli
	}
	return nil
}

func (r *MultiRedis) Errors() []error {
	return r.errors
}

func (r *MultiRedis) Close() error {
	var firstErr error
	for _, cli := range r.conns {
		if err := cli.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
