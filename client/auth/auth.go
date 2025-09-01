package auth

import (
	"context"
	"time"
)

// Authenticator defines an interface for obtaining and managing access tokens.
// Implementations should provide a method to retrieve an access token, as well as
// a method to check the current authentication status and token expiration time.
//
// The Status() method returns the current authentication status as a string,
// the token's expiration time as a time.Time, and an error if applicable.
type Authenticator interface {
	GetAccessToken(ctx context.Context) (string, error)
	Status() (string, time.Time, error)
}
