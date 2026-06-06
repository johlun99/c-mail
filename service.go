package main

import "cmail/internal/mail"

// MailService is bound to the frontend by Wails and exposes the mail data layer.
// It currently wraps the mock store; a Gmail-backed Store will replace it later
// without changing this surface.
type MailService struct {
	store mail.Store
}

// NewMailService creates a MailService backed by the mock dataset.
func NewMailService() *MailService {
	return &MailService{store: mail.NewMockStore()}
}

// GetMails returns all mails.
func (s *MailService) GetMails() []mail.Mail { return s.store.Mails() }

// GetCategories returns the triage taxonomy.
func (s *MailService) GetCategories() []mail.Category { return s.store.Categories() }

// GetActivity returns the agent activity feed.
func (s *MailService) GetActivity() []mail.Activity { return s.store.Activity() }

// GetRules returns the agent automation rules.
func (s *MailService) GetRules() []mail.Rule { return s.store.Rules() }

// GetAccounts returns the linked mail accounts.
func (s *MailService) GetAccounts() []mail.Account { return s.store.Accounts() }
