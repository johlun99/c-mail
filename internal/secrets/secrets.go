// Package secrets stores sensitive values — the OAuth refresh token and the
// local-cache encryption key — in the OS keyring (GNOME Keyring / Secret Service
// on Linux, Keychain on macOS). Nothing sensitive is written to disk in
// plaintext.
package secrets

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/zalando/go-keyring"
)

// service is the keyring service name all cmail secrets are filed under.
const service = "cmail"

const cacheKeyName = "cache:key"

// ErrNotFound is returned when a requested secret does not exist.
var ErrNotFound = errors.New("secrets: not found")

func refreshKey(email string) string { return "gmail:refresh:" + email }

// SaveRefreshToken stores the OAuth refresh token for an account.
func SaveRefreshToken(email, token string) error {
	return keyring.Set(service, refreshKey(email), token)
}

// LoadRefreshToken returns the stored refresh token for an account, or
// ErrNotFound if none is stored.
func LoadRefreshToken(email string) (string, error) {
	t, err := keyring.Get(service, refreshKey(email))
	if errors.Is(err, keyring.ErrNotFound) {
		return "", ErrNotFound
	}
	return t, err
}

// DeleteRefreshToken removes the stored refresh token for an account. Deleting a
// missing token is not an error.
func DeleteRefreshToken(email string) error {
	err := keyring.Delete(service, refreshKey(email))
	if errors.Is(err, keyring.ErrNotFound) {
		return nil
	}
	return err
}

// CacheKey returns the 32-byte AES key used to encrypt the local mail cache,
// generating and storing it in the keyring on first use.
func CacheKey() ([]byte, error) {
	v, err := keyring.Get(service, cacheKeyName)
	if err == nil {
		return base64.StdEncoding.DecodeString(v)
	}
	if !errors.Is(err, keyring.ErrNotFound) {
		return nil, err
	}
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	if err := keyring.Set(service, cacheKeyName, base64.StdEncoding.EncodeToString(key)); err != nil {
		return nil, fmt.Errorf("secrets: store cache key: %w", err)
	}
	return key, nil
}
