package main

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"cmail/internal/cache"
	"cmail/internal/gmail"
	"cmail/internal/mail"
	"cmail/internal/secrets"
)

// MailService is bound to the frontend by Wails. It serves mock data until a
// Gmail account is connected, after which it serves live Gmail data backed by an
// encrypted local cache. The same surface is used in both modes.
type MailService struct {
	mock mail.Store

	mu      sync.Mutex
	client  *gmail.Client
	cache   *cache.Cache
	account string
	mails   []mail.Mail
}

// NewMailService creates a MailService that starts in mock mode.
func NewMailService() *MailService {
	return &MailService{mock: mail.NewMockStore()}
}

// GetMails returns live Gmail mails when connected, otherwise mock mails.
func (s *MailService) GetMails() []mail.Mail {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.account != "" {
		return append([]mail.Mail(nil), s.mails...)
	}
	return s.mock.Mails()
}

// GetCategories returns the triage taxonomy (app-defined, identical in both modes).
func (s *MailService) GetCategories() []mail.Category { return s.mock.Categories() }

// GetActivity returns the agent activity feed.
func (s *MailService) GetActivity() []mail.Activity { return s.mock.Activity() }

// GetRules returns the agent automation rules.
func (s *MailService) GetRules() []mail.Rule { return s.mock.Rules() }

// GetAccounts returns the connected account, or an empty list when not connected.
func (s *MailService) GetAccounts() []mail.Account {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.account == "" {
		return []mail.Account{}
	}
	return []mail.Account{{
		Email:    s.account,
		Provider: "Gmail",
		Status:   mail.StatusConnected,
		SyncedAt: "synkad nyss",
		Scopes: []string{
			"läser e-post",
			"skapar utkast, skickar bara med ditt godkännande",
			"hanterar etiketter",
		},
	}}
}

// GmailConfigured reports whether OAuth client credentials are available.
func (s *MailService) GmailConfigured() bool { return gmail.Configured() }

// GmailConnected reports whether an account is currently connected.
func (s *MailService) GmailConnected() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.account != ""
}

// ConnectGmail runs the OAuth flow, persists the refresh token in the keyring,
// and loads mail. It blocks until the user completes (or cancels) the browser
// consent, or the timeout elapses.
func (s *MailService) ConnectGmail() (mail.Account, error) {
	if !gmail.Configured() {
		return mail.Account{}, gmail.ErrNoCredentials
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	tok, err := gmail.Authenticate(ctx, openURL)
	if err != nil {
		return mail.Account{}, err
	}
	client, err := gmail.NewClient(context.Background(), tok.RefreshToken)
	if err != nil {
		return mail.Account{}, err
	}
	email, err := client.Profile(ctx)
	if err != nil {
		return mail.Account{}, err
	}
	if err := secrets.SaveRefreshToken(email, tok.RefreshToken); err != nil {
		return mail.Account{}, err
	}
	if err := s.activate(email, client); err != nil {
		return mail.Account{}, err
	}
	return s.GetAccounts()[0], nil
}

// DisconnectGmail clears the session, closes the cache and deletes the token.
func (s *MailService) DisconnectGmail() error {
	s.mu.Lock()
	email := s.account
	if s.cache != nil {
		_ = s.cache.Close()
	}
	s.client, s.cache, s.account, s.mails = nil, nil, "", nil
	s.mu.Unlock()
	if email != "" {
		return secrets.DeleteRefreshToken(email)
	}
	return nil
}

// RefreshMails re-fetches inbox mail from Gmail and updates the cache. No-op when
// not connected.
func (s *MailService) RefreshMails() error {
	s.mu.Lock()
	client, c := s.client, s.cache
	s.mu.Unlock()
	if client == nil {
		return nil
	}
	mails, err := client.FetchMails(context.Background(), 50)
	if err != nil {
		return err
	}
	if c != nil {
		_ = c.PutMails(mails)
	}
	s.mu.Lock()
	s.mails = mails
	s.mu.Unlock()
	return nil
}

// activate opens the encrypted cache and loads mail for a connected account,
// falling back to cached data if the initial network fetch fails.
func (s *MailService) activate(email string, client *gmail.Client) error {
	key, err := secrets.CacheKey()
	if err != nil {
		return err
	}
	path, err := cachePath()
	if err != nil {
		return err
	}
	c, err := cache.Open(path, key)
	if err != nil {
		return err
	}
	s.mu.Lock()
	s.client, s.cache, s.account = client, c, email
	s.mu.Unlock()

	if err := s.RefreshMails(); err != nil {
		if cached, cerr := c.Mails(); cerr == nil && len(cached) > 0 {
			s.mu.Lock()
			s.mails = cached
			s.mu.Unlock()
			return nil
		}
		return err
	}
	return nil
}

func cachePath() (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	dir = filepath.Join(dir, "cmail")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	return filepath.Join(dir, "cache.db"), nil
}

func openURL(target string) error {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", target).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", target).Start()
	default:
		return exec.Command("xdg-open", target).Start()
	}
}
