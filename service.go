package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	wruntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"cmail/internal/agent"
	"cmail/internal/cache"
	"cmail/internal/gmail"
	"cmail/internal/mail"
	"cmail/internal/secrets"
)

// syncInterval is how often connected mail is refreshed in the background.
const syncInterval = 2 * time.Minute

// MailService is bound to the frontend by Wails. It serves mock data until a
// Gmail account is connected, after which it serves live Gmail data backed by an
// encrypted local cache. The same surface is used in both modes.
type MailService struct {
	mock       mail.Store
	classifier agent.Classifier

	mu      sync.Mutex
	ctx     context.Context
	client  *gmail.Client
	cache   *cache.Cache
	account string
	mails   []mail.Mail
}

// NewMailService creates a MailService that starts in mock mode, with a local
// Ollama classifier for the agent flow.
func NewMailService() *MailService {
	return &MailService{mock: mail.NewMockStore(), classifier: agent.NewOllama()}
}

// AIAvailable reports whether the local AI backend (Ollama) is reachable.
func (s *MailService) AIAvailable() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return s.classifier.Available(ctx)
}

// GenerateDraft asks the local agent to draft a reply for a mail.
func (s *MailService) GenerateDraft(mailID string) (mail.Draft, error) {
	m, ok := s.findMail(mailID)
	if !ok {
		return mail.Draft{}, errors.New("mail saknas")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	return s.classifier.DraftReply(ctx, agent.MailInput{
		From:    m.From,
		Subject: m.Subject,
		Body:    strings.Join(m.Body, "\n\n"),
	}, "professionell · varm")
}

func (s *MailService) findMail(id string) (mail.Mail, bool) {
	for _, m := range s.GetMails() {
		if m.ID == id {
			return m, true
		}
	}
	return mail.Mail{}, false
}

// classify runs the local agent over fetched mail (best-effort): if the backend
// is unavailable it leaves the label-based category in place.
func (s *MailService) classify(ctx context.Context, mails []mail.Mail) []mail.Mail {
	if !s.classifier.Available(ctx) {
		return mails
	}
	cats := s.mock.Categories()
	for i := range mails {
		c, err := s.classifier.Categorize(ctx, agent.MailInput{
			From:    mails[i].From,
			Subject: mails[i].Subject,
			Body:    strings.Join(mails[i].Body, "\n\n"),
		}, cats)
		if err != nil {
			continue
		}
		mails[i].Cat = c.Category
		mails[i].Confidence = c.Confidence
		mails[i].Agent = mail.AgentAnalysis{Summary: c.Summary, Facts: c.Facts, Tasks: c.Tasks}
	}
	return mails
}

// Start records the Wails context, restores a previous session (if any) and
// begins background sync. Called from App.startup.
func (s *MailService) Start(ctx context.Context) {
	s.mu.Lock()
	s.ctx = ctx
	s.mu.Unlock()
	go s.restore()
	go s.syncLoop(ctx)
}

// restore reconnects to the last-connected account using the keyring-stored
// refresh token, so the user does not re-authenticate on every launch.
func (s *MailService) restore() {
	email, err := loadState()
	if err != nil || email == "" || !gmail.Configured() {
		return
	}
	tok, err := secrets.LoadRefreshToken(email)
	if err != nil {
		return
	}
	client, err := gmail.NewClient(context.Background(), tok)
	if err != nil {
		return
	}
	if err := s.activate(email, client); err != nil {
		return
	}
	s.emit("mails:updated")
}

func (s *MailService) syncLoop(ctx context.Context) {
	t := time.NewTicker(syncInterval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			if !s.GmailConnected() {
				continue
			}
			if err := s.RefreshMails(); err == nil {
				s.emit("mails:updated")
			}
		}
	}
}

func (s *MailService) emit(event string) {
	s.mu.Lock()
	ctx := s.ctx
	s.mu.Unlock()
	if ctx != nil {
		wruntime.EventsEmit(ctx, event)
	}
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
	_ = saveState(email)
	s.emit("mails:updated")
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
	_ = saveState("")
	s.emit("mails:updated")
	if email != "" {
		return secrets.DeleteRefreshToken(email)
	}
	return nil
}

// SendReply sends a reply. It is only ever invoked from an explicit user-approval
// action in the UI — never autonomously (the never-send guarantee).
func (s *MailService) SendReply(to, subject, body string) error {
	s.mu.Lock()
	client := s.client
	s.mu.Unlock()
	if client == nil {
		return errors.New("inget anslutet konto")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return client.Send(ctx, to, subject, body)
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	mails, err := client.FetchMails(ctx, 50)
	if err != nil {
		return err
	}
	mails = s.classify(ctx, mails)
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

// persistState records which account is connected so the session can be restored
// on the next launch (the token itself stays in the keyring).
type persistState struct {
	Account string `json:"account"`
}

func statePath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	dir = filepath.Join(dir, "cmail")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	return filepath.Join(dir, "state.json"), nil
}

func saveState(email string) error {
	p, err := statePath()
	if err != nil {
		return err
	}
	b, err := json.Marshal(persistState{Account: email})
	if err != nil {
		return err
	}
	return os.WriteFile(p, b, 0o600)
}

func loadState() (string, error) {
	p, err := statePath()
	if err != nil {
		return "", err
	}
	b, err := os.ReadFile(p)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", nil
		}
		return "", err
	}
	var st persistState
	if err := json.Unmarshal(b, &st); err != nil {
		return "", err
	}
	return st.Account, nil
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
