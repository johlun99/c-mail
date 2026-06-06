package gmail

import "testing"

func TestConfigured(t *testing.T) {
	t.Setenv("GMAIL_CLIENT_ID", "")
	t.Setenv("GMAIL_CLIENT_SECRET", "")
	if Configured() {
		t.Error("Configured() should be false with no credentials")
	}

	t.Setenv("GMAIL_CLIENT_ID", "id")
	t.Setenv("GMAIL_CLIENT_SECRET", "secret")
	if !Configured() {
		t.Error("Configured() should be true with both credentials set")
	}

	t.Setenv("GMAIL_CLIENT_SECRET", "")
	if Configured() {
		t.Error("Configured() should be false when the secret is missing")
	}
}

func TestOAuthConfigRequiresCredentials(t *testing.T) {
	t.Setenv("GMAIL_CLIENT_ID", "")
	t.Setenv("GMAIL_CLIENT_SECRET", "")
	if _, err := oauthConfig("http://127.0.0.1:1/callback"); err == nil {
		t.Error("oauthConfig should error without credentials")
	}

	t.Setenv("GMAIL_CLIENT_ID", "id")
	t.Setenv("GMAIL_CLIENT_SECRET", "secret")
	cfg, err := oauthConfig("http://127.0.0.1:1/callback")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.ClientID != "id" || cfg.RedirectURL != "http://127.0.0.1:1/callback" {
		t.Errorf("unexpected config: %+v", cfg)
	}
}
