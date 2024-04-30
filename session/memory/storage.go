package memory

import (
	"context"
	"errors"
	"sync"
	"zweb/session"
	"time"

	cache "github.com/patrickmn/go-cache"
)

var _ session.Store = &Store{}

type Store struct {
	mutex sync.RWMutex
	cache *cache.Cache
	expiration time.Duration
}

func NewStore(expiration time.Duration) *Store {
	return &Store{
		cache: cache.New(expiration, time.Second),
		expiration: expiration,
	}
}

// Generate implements session.Store.
func (s *Store) Generate(ctx context.Context, id string) (session.Session, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	session := &Session{
		id: id,
		data: map[string]any{},
	}
	s.cache.Set(id, session, s.expiration)
	return session, nil
}

// Get implements session.Store.
func (s *Store) Get(ctx context.Context, id string) (session.Session, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	session, ok := s.cache.Get(id)
	if !ok {
		return nil, errors.New("key not found")
	}
	return session.(*Session), nil
}

// Refresh implements session.Store.
func (s *Store) Refresh(ctx context.Context, id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	sess, ok := s.cache.Get(id)
	if !ok {
		return errors.New("key not found")
	}
	s.cache.Set(id, sess, s.expiration)
	return nil
}

// Remove implements session.Store.
func (s *Store) Remove(ctx context.Context, id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.cache.Delete(id)
	return nil
}

var _ session.Session = &Session{}

type Session struct {
	mutex sync.RWMutex
	id string
	data map[string]any
}

// Get implements session.Session.
func (s *Session) Get(ctx context.Context, key string) (any, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	val, ok := s.data[key]
	if !ok {
		return nil, errors.New("key not found")
	}
	return val, nil
}

// ID implements session.Session.
func (s *Session) ID() string {
	return s.id
}

// Set implements session.Session.
func (s *Session) Set(ctx context.Context, key string, val any) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.data[key] = val
	return nil
}
