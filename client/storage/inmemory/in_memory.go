package inmemory

import (
	"sync"
	"time"

	"github.com/tkliner/go-gopay/client/storage"
)

// InMemoryTokenStorage is an in-memory implementation of the TokenStorage interface.
// It is safe for concurrent use.
type InMemoryTokenStorage struct {
	mu        sync.RWMutex
	token     string
	expiresAt time.Time
}

// NewInMemoryTokenStorage creates a new instance of InMemoryTokenStorage.
func NewInMemoryTokenStorage() storage.TokenStorage {
	return &InMemoryTokenStorage{}
}

// SaveAccessToken stores the access token and its expiration time in memory.
// This method is thread-safe.
func (s *InMemoryTokenStorage) SaveAccessToken(token string, expiresAt time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.token = token
	s.expiresAt = expiresAt
	return nil
}

// GetAccessToken retrieves the stored access token and its expiration time.
// This method is thread-safe.
func (s *InMemoryTokenStorage) GetAccessToken() (string, time.Time, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.token, s.expiresAt, nil
}