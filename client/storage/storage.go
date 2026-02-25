// Package storage defines the token storage interface used by the auth layer.
package storage

import "time"

// TokenStorage is the interface for storing and retrieving OAuth access tokens.
type TokenStorage interface {
	// SaveAccessToken stores the access token and its expiration time.
	SaveAccessToken(token string, expiresAt time.Time) error
	// GetAccessToken returns the stored access token and expiration time.
	GetAccessToken() (string, time.Time, error)
}