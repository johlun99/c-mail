package gmail

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"golang.org/x/oauth2"
	gmailapi "google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"

	cmail "cmail/internal/mail"
)

// Client is an authenticated Gmail API client.
type Client struct {
	svc *gmailapi.Service
}

// NewClient builds a Gmail client from a stored refresh token. The token source
// transparently mints fresh access tokens as needed.
func NewClient(ctx context.Context, refreshToken string) (*Client, error) {
	cfg, err := oauthConfig("")
	if err != nil {
		return nil, err
	}
	ts := cfg.TokenSource(ctx, &oauth2.Token{RefreshToken: refreshToken})
	svc, err := gmailapi.NewService(ctx, option.WithTokenSource(ts))
	if err != nil {
		return nil, fmt.Errorf("gmail: new service: %w", err)
	}
	return &Client{svc: svc}, nil
}

// Profile returns the authenticated account's email address.
func (c *Client) Profile(ctx context.Context) (string, error) {
	p, err := c.svc.Users.GetProfile("me").Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("gmail: profile: %w", err)
	}
	return p.EmailAddress, nil
}

// FetchMails fetches up to max inbox messages, mapped to the domain model.
func (c *Client) FetchMails(ctx context.Context, max int64) ([]cmail.Mail, error) {
	list, err := c.svc.Users.Messages.List("me").Q("in:inbox").MaxResults(max).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("gmail: list messages: %w", err)
	}
	now := time.Now()
	out := make([]cmail.Mail, 0, len(list.Messages))
	for _, ref := range list.Messages {
		full, err := c.svc.Users.Messages.Get("me", ref.Id).Format("full").Context(ctx).Do()
		if err != nil {
			return nil, fmt.Errorf("gmail: get message %s: %w", ref.Id, err)
		}
		out = append(out, mapMessage(full, now))
	}
	return out, nil
}

// CreateDraft creates a Gmail draft reply. Drafting is safe to do automatically;
// it never sends.
func (c *Client) CreateDraft(ctx context.Context, to, subject, body string) error {
	enc := base64.URLEncoding.EncodeToString([]byte(buildMIME(to, subject, body)))
	_, err := c.svc.Users.Drafts.Create("me", &gmailapi.Draft{
		Message: &gmailapi.Message{Raw: enc},
	}).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("gmail: create draft: %w", err)
	}
	return nil
}

// Send sends a reply. This must only be called in response to explicit user
// approval — the agent never calls it autonomously (the never-send guarantee).
func (c *Client) Send(ctx context.Context, to, subject, body string) error {
	enc := base64.URLEncoding.EncodeToString([]byte(buildMIME(to, subject, body)))
	_, err := c.svc.Users.Messages.Send("me", &gmailapi.Message{Raw: enc}).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("gmail: send: %w", err)
	}
	return nil
}

func buildMIME(to, subject, body string) string {
	return "To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"\r\n" + body
}
