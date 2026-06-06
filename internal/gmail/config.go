// Package gmail authenticates against Gmail (desktop loopback OAuth) and maps
// Gmail messages into cmail's domain model.
//
// Never-send guarantee: the agent never sends autonomously. Drafts are created
// automatically, but a message is sent only in response to the user's explicit
// approval. This is enforced in application logic (Send is only ever called from
// a user-approval action), not by withholding the scope.
package gmail

import (
	"errors"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gmailapi "google.golang.org/api/gmail/v1"
)

// scopes: gmail.modify covers reading mail, creating/managing drafts, managing
// labels, and sending (only ever invoked on explicit user approval).
var scopes = []string{gmailapi.GmailModifyScope}

// ErrNoCredentials is returned when the OAuth client credentials are not set.
var ErrNoCredentials = errors.New("gmail: GMAIL_CLIENT_ID/GMAIL_CLIENT_SECRET not set")

// Configured reports whether OAuth client credentials are available in the env.
func Configured() bool {
	return os.Getenv("GMAIL_CLIENT_ID") != "" && os.Getenv("GMAIL_CLIENT_SECRET") != ""
}

func oauthConfig(redirectURL string) (*oauth2.Config, error) {
	id := os.Getenv("GMAIL_CLIENT_ID")
	secret := os.Getenv("GMAIL_CLIENT_SECRET")
	if id == "" || secret == "" {
		return nil, ErrNoCredentials
	}
	return &oauth2.Config{
		ClientID:     id,
		ClientSecret: secret,
		Endpoint:     google.Endpoint,
		RedirectURL:  redirectURL,
		Scopes:       scopes,
	}, nil
}
