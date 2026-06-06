package agent

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"cmail/internal/mail"
)

var cats = []mail.Category{
	{Key: "svara", Label: "att svara på"},
	{Key: "reklam", Label: "reklam"},
}

func reply(t *testing.T, content string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/tags":
			w.WriteHeader(http.StatusOK)
		case "/api/chat":
			_ = json.NewEncoder(w).Encode(chatResponse{Message: chatMessage{Role: "assistant", Content: content}})
		default:
			http.NotFound(w, r)
		}
	}))
}

func TestAvailable(t *testing.T) {
	srv := reply(t, "")
	o := NewOllamaWith(srv.URL, "m")
	if !o.Available(context.Background()) {
		t.Error("expected Available true for a live server")
	}
	srv.Close()
	if o.Available(context.Background()) {
		t.Error("expected Available false after server closed")
	}
}

func TestCategorize(t *testing.T) {
	content := `{"category":"svara","confidence":0.9,"summary":"Kund vill ha svar.","facts":[{"key":"deadline","value":"fre"}],"tasks":["Svara"]}`
	srv := reply(t, content)
	defer srv.Close()

	o := NewOllamaWith(srv.URL, "m")
	c, err := o.Categorize(context.Background(), MailInput{From: "Anna", Subject: "Offert"}, cats)
	if err != nil {
		t.Fatal(err)
	}
	if c.Category != "svara" || c.Confidence != 0.9 {
		t.Errorf("got %+v", c)
	}
	if len(c.Facts) != 1 || c.Facts[0].Key != "deadline" {
		t.Errorf("facts = %+v", c.Facts)
	}
	if len(c.Tasks) != 1 {
		t.Errorf("tasks = %+v", c.Tasks)
	}
}

func TestCategorizeFallsBackOnUnknownCategory(t *testing.T) {
	srv := reply(t, `{"category":"hittepå","confidence":0.5,"summary":"x"}`)
	defer srv.Close()

	o := NewOllamaWith(srv.URL, "m")
	c, err := o.Categorize(context.Background(), MailInput{}, cats)
	if err != nil {
		t.Fatal(err)
	}
	if c.Category != "svara" {
		t.Errorf("unknown category should fall back to first; got %q", c.Category)
	}
	if c.Facts == nil || c.Tasks == nil {
		t.Error("facts/tasks should be non-nil")
	}
}

func TestDraftReply(t *testing.T) {
	srv := reply(t, "Hej Anna,\n\nTack för ditt mail.")
	defer srv.Close()

	o := NewOllamaWith(srv.URL, "m")
	d, err := o.DraftReply(context.Background(), MailInput{From: "Anna"}, "varm")
	if err != nil {
		t.Fatal(err)
	}
	if d.Tone != "varm" {
		t.Errorf("tone = %q", d.Tone)
	}
	if len(d.Lines) != 3 {
		t.Fatalf("lines = %d, want 3", len(d.Lines))
	}
	if d.Lines[0].Text != "Hej Anna," || d.Lines[0].Kind != mail.DraftAdded {
		t.Errorf("first line = %+v", d.Lines[0])
	}
}
