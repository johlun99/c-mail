package gmail

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"

	"golang.org/x/oauth2"
)

// Authenticate runs the desktop loopback OAuth flow: it starts a localhost
// server on an ephemeral port, asks openBrowser to open Google's consent page,
// captures the redirect, and exchanges the code (with PKCE) for a token.
//
// The returned token carries the long-lived refresh token, which the caller
// should persist securely (the OS keyring).
func Authenticate(ctx context.Context, openBrowser func(url string) error) (*oauth2.Token, error) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("gmail: listen: %w", err)
	}
	defer func() { _ = ln.Close() }()

	redirect := fmt.Sprintf("http://127.0.0.1:%d/callback", ln.Addr().(*net.TCPAddr).Port)
	cfg, err := oauthConfig(redirect)
	if err != nil {
		return nil, err
	}

	state, err := randString()
	if err != nil {
		return nil, err
	}
	verifier := oauth2.GenerateVerifier()
	authURL := cfg.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
		oauth2.ApprovalForce,
		oauth2.S256ChallengeOption(verifier),
	)

	type result struct {
		code string
		err  error
	}
	resCh := make(chan result, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("state") != state {
			http.Error(w, "state mismatch", http.StatusBadRequest)
			resCh <- result{err: fmt.Errorf("gmail: state mismatch")}
			return
		}
		if e := q.Get("error"); e != "" {
			http.Error(w, e, http.StatusBadRequest)
			resCh <- result{err: fmt.Errorf("gmail: authorization denied: %s", e)}
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte("<p>cmail: inloggning klar — du kan stänga den här fliken.</p>"))
		resCh <- result{code: q.Get("code")}
	})
	srv := &http.Server{Handler: mux}
	go func() { _ = srv.Serve(ln) }()
	defer func() { _ = srv.Shutdown(context.Background()) }()

	if err := openBrowser(authURL); err != nil {
		return nil, fmt.Errorf("gmail: open browser: %w", err)
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-resCh:
		if res.err != nil {
			return nil, res.err
		}
		tok, err := cfg.Exchange(ctx, res.code, oauth2.VerifierOption(verifier))
		if err != nil {
			return nil, fmt.Errorf("gmail: token exchange: %w", err)
		}
		return tok, nil
	}
}

func randString() (string, error) {
	b := make([]byte, 24)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
