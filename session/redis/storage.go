package redis

import (
	"context"
	"errors"
	"fmt"
	"zweb/session"
	"time"

	"github.com/redis/go-redis/v9"
)

var _ session.Store = &Store{}

type Store struct {
	prefix     string
	expiration time.Duration
	client     redis.Cmdable
}

func NewStore(client redis.Cmdable) *Store {
	return &Store{
		client: client,
		expiration: 15 * time.Minute,
		prefix: "session",
	}
}

// Generate implements session.Store.
func (s *Store) Generate(ctx context.Context, id string) (session.Session, error) {
	key := redisKey(s.prefix, id)
	_, err := s.client.HSet(ctx, key, id, id).Result()
	if err != nil {
		return nil, err
	}
	_, err = s.client.Expire(ctx, key, s.expiration).Result()
	if err != nil {
		return nil, err
	}
	return &Session{id: id, key: key, client: s.client}, err
}

// Get implements session.Store.
func (s *Store) Get(ctx context.Context, id string) (session.Session, error) {
	key := redisKey(s.prefix, id)
	i, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if i < 0 {
		return nil, errors.New("session not exist")
	}
	return &Session{id: id, key: key, client: s.client}, nil

}

// Refresh implements session.Store.
func (s *Store) Refresh(ctx context.Context, id string) error {
	key := redisKey(s.prefix, id)
	_, err := s.client.Expire(ctx, key, s.expiration).Result()
	return err
}

// Remove implements session.Store.
func (s *Store) Remove(ctx context.Context, id string) error {
	key := redisKey(s.prefix, id)
	_, err := s.client.Del(ctx, key).Result()
	return err
}

func redisKey(prefix, key string) string {
	return fmt.Sprintf("%s_%s", prefix, key)
}

var _ session.Session = &Session{}

type Session struct {
	id     string
	key string
	client redis.Cmdable
}

// Get implements session.Session.
func (s *Session) Get(ctx context.Context, key string) (any, error) {
	return s.client.HGet(ctx, s.key, key).Result()
}

// ID implements session.Session.
func (s *Session) ID() string {
	return s.id
}

// Set implements session.Session.
func (s *Session) Set(ctx context.Context, key string, val any) error {
	const lua = `
if redis.call("exists", KEYS[1])
then
	return redis.call("hset", KEYS[1], ARGV[1], ARGV[2])
else
	return -1
end`
	res, err := s.client.Eval(ctx, lua, []string{s.key}, key, val).Int()
	if err != nil {
		return err
	}
	if res < 0 {
		return errors.New("session not exist")
	}
	return nil
}
