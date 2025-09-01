package storage

import "time"

// TokenStorage definuje rozhraní pro ukládání a načítání přístupových tokenů.
type TokenStorage interface {
	// SaveAccessToken uloží přístupový token s danou dobou platnosti.
	SaveAccessToken(token string, expiresAt time.Time) error
	// GetAccessToken načte uložený přístupový token.
	GetAccessToken() (string, time.Time, error)
}