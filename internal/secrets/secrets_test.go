package secrets

import (
	"errors"
	"testing"

	"github.com/zalando/go-keyring"
)

func TestRefreshTokenRoundTrip(t *testing.T) {
	keyring.MockInit()

	if _, err := LoadRefreshToken("a@b.se"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
	if err := SaveRefreshToken("a@b.se", "tok-123"); err != nil {
		t.Fatal(err)
	}
	got, err := LoadRefreshToken("a@b.se")
	if err != nil {
		t.Fatal(err)
	}
	if got != "tok-123" {
		t.Errorf("got %q, want tok-123", got)
	}
	if err := DeleteRefreshToken("a@b.se"); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadRefreshToken("a@b.se"); !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestDeleteMissingTokenIsNoError(t *testing.T) {
	keyring.MockInit()
	if err := DeleteRefreshToken("missing@b.se"); err != nil {
		t.Errorf("deleting missing token should be a no-op, got %v", err)
	}
}

func TestCacheKeyIsStableAndCorrectSize(t *testing.T) {
	keyring.MockInit()
	k1, err := CacheKey()
	if err != nil {
		t.Fatal(err)
	}
	if len(k1) != 32 {
		t.Errorf("cache key len = %d, want 32", len(k1))
	}
	k2, err := CacheKey()
	if err != nil {
		t.Fatal(err)
	}
	if string(k1) != string(k2) {
		t.Error("cache key should be stable across calls")
	}
}
