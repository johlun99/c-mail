package cache

import (
	"crypto/rand"
	"path/filepath"
	"testing"

	"cmail/internal/mail"
)

func newKey(t *testing.T) []byte {
	t.Helper()
	k := make([]byte, 32)
	if _, err := rand.Read(k); err != nil {
		t.Fatal(err)
	}
	return k
}

func sample() []mail.Mail {
	return []mail.Mail{
		{
			ID: "m1", Cat: "svara", From: "Anna", Subject: "Hej", Body: []string{"rad 1", "rad 2"},
			Agent: mail.AgentAnalysis{Summary: "s", Facts: []mail.Fact{{Key: "k", Value: "v"}}},
			Draft: &mail.Draft{Tone: "varm", Lines: []mail.DraftLine{{Text: "Hej", Kind: mail.DraftAdded}}},
		},
		{ID: "m2", Cat: "faktura", From: "Fortnox", Subject: "Faktura"},
	}
}

func TestPutAndGetRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cache.db")
	key := newKey(t)
	c, err := Open(path, key)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = c.Close() }()

	if err := c.PutMails(sample()); err != nil {
		t.Fatal(err)
	}
	got, err := c.Mails()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("got %d mails, want 2", len(got))
	}
	if got[0].ID != "m1" || got[1].ID != "m2" {
		t.Errorf("order not preserved: %s, %s", got[0].ID, got[1].ID)
	}
	if got[0].Draft == nil || got[0].Draft.Lines[0].Text != "Hej" {
		t.Error("nested draft did not round-trip")
	}
	if got[0].Agent.Facts[0].Key != "k" {
		t.Error("nested agent analysis did not round-trip")
	}
}

func TestContentIsEncryptedAtRest(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cache.db")
	c, err := Open(path, newKey(t))
	if err != nil {
		t.Fatal(err)
	}
	if err := c.PutMails(sample()); err != nil {
		t.Fatal(err)
	}
	// The plaintext subject must not appear in the stored payload blob.
	var blob []byte
	if err := c.db.QueryRow(`SELECT payload FROM mails WHERE id = 'm1'`).Scan(&blob); err != nil {
		t.Fatal(err)
	}
	_ = c.Close()
	if len(blob) == 0 {
		t.Fatal("empty payload")
	}
	for i := 0; i+3 <= len(blob); i++ {
		if string(blob[i:i+3]) == "Hej" {
			t.Fatal("plaintext content found in stored blob — not encrypted")
		}
	}
}

func TestWrongKeyFailsToDecrypt(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cache.db")
	c, err := Open(path, newKey(t))
	if err != nil {
		t.Fatal(err)
	}
	if err := c.PutMails(sample()); err != nil {
		t.Fatal(err)
	}
	_ = c.Close()

	c2, err := Open(path, newKey(t)) // different key
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = c2.Close() }()
	if _, err := c2.Mails(); err == nil {
		t.Error("expected decryption to fail with the wrong key")
	}
}
