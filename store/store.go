package store

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Store struct {
	store *redis.Client
}

func NewStore(conf *redis.Options) *Store {
	return &Store{
		store: redis.NewClient(conf),
	}
}

func (s *Store) Set(ctx context.Context, key string, dest any, ttl time.Duration) error {
	return s.store.Set(ctx, key, dest, ttl).Err()
}

func (s *Store) Get(ctx context.Context, key string) ([]byte, error) {
	return s.store.Get(ctx, key).Bytes()
}

func (s *Store) Del(ctx context.Context, key ...string) error {
	return s.store.Del(ctx, key...).Err()
}
