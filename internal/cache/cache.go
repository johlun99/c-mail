// Package cache is an at-rest-encrypted local cache of mail. Each mail is stored
// as an AES-256-GCM encrypted JSON blob; the key lives in the OS keyring, so the
// cache is only readable while the keyring is unlocked. Only opaque message IDs
// and an ordering index are stored in plaintext.
package cache

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"

	_ "modernc.org/sqlite" // pure-Go SQLite driver (no cgo)

	"cmail/internal/mail"
)

// Cache is an encrypted, on-disk mail cache.
type Cache struct {
	db   *sql.DB
	aead cipher.AEAD
}

// Open opens (or creates) the cache at path, using key (32 bytes for AES-256)
// to encrypt mail content.
func Open(path string, key []byte) (*Cache, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("cache: %w", err)
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("cache: %w", err)
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("cache: open db: %w", err)
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS mails (id TEXT PRIMARY KEY, ord INTEGER, payload BLOB)`); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("cache: init schema: %w", err)
	}
	return &Cache{db: db, aead: aead}, nil
}

// Close closes the underlying database.
func (c *Cache) Close() error { return c.db.Close() }

func (c *Cache) seal(plain []byte) ([]byte, error) {
	nonce := make([]byte, c.aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	return c.aead.Seal(nonce, nonce, plain, nil), nil
}

func (c *Cache) unseal(blob []byte) ([]byte, error) {
	ns := c.aead.NonceSize()
	if len(blob) < ns {
		return nil, fmt.Errorf("cache: ciphertext too short")
	}
	return c.aead.Open(nil, blob[:ns], blob[ns:], nil)
}

// PutMails atomically replaces the cache contents with mails, preserving order.
func (c *Cache) PutMails(mails []mail.Mail) error {
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.Exec(`DELETE FROM mails`); err != nil {
		return err
	}
	stmt, err := tx.Prepare(`INSERT INTO mails (id, ord, payload) VALUES (?, ?, ?)`)
	if err != nil {
		return err
	}
	defer func() { _ = stmt.Close() }()

	for i, m := range mails {
		j, err := json.Marshal(m)
		if err != nil {
			return err
		}
		blob, err := c.seal(j)
		if err != nil {
			return err
		}
		if _, err := stmt.Exec(m.ID, i, blob); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// Mails returns the cached mails in stored order, decrypting each.
func (c *Cache) Mails() ([]mail.Mail, error) {
	rows, err := c.db.Query(`SELECT payload FROM mails ORDER BY ord`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var out []mail.Mail
	for rows.Next() {
		var blob []byte
		if err := rows.Scan(&blob); err != nil {
			return nil, err
		}
		plain, err := c.unseal(blob)
		if err != nil {
			return nil, fmt.Errorf("cache: decrypt: %w", err)
		}
		var m mail.Mail
		if err := json.Unmarshal(plain, &m); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}
