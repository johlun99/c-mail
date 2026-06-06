package main

import (
	"context"
	"testing"

	"cmail/internal/agent"
	"cmail/internal/mail"
)

// fakeClassifier is a deterministic Classifier for testing the service without a
// live Ollama backend.
type fakeClassifier struct {
	avail bool
	draft mail.Draft
}

func (f fakeClassifier) Available(context.Context) bool { return f.avail }
func (f fakeClassifier) Categorize(context.Context, agent.MailInput, []mail.Category) (agent.Classification, error) {
	return agent.Classification{Category: "svara", Confidence: 0.5}, nil
}
func (f fakeClassifier) DraftReply(context.Context, agent.MailInput, string) (mail.Draft, error) {
	return f.draft, nil
}

func TestMockModeWhenNotConnected(t *testing.T) {
	s := NewMailService()

	if len(s.GetMails()) == 0 {
		t.Error("expected mock mails when not connected")
	}
	if got := len(s.GetAccounts()); got != 0 {
		t.Errorf("GetAccounts() = %d, want 0 when not connected", got)
	}
	if s.GmailConnected() {
		t.Error("GmailConnected() should be false initially")
	}
	if len(s.GetCategories()) == 0 {
		t.Error("expected categories in both modes")
	}
}

func TestStateRoundTrip(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	got, err := loadState()
	if err != nil {
		t.Fatal(err)
	}
	if got != "" {
		t.Errorf("loadState() with no file = %q, want empty", got)
	}

	if err := saveState("johan@nordveda.se"); err != nil {
		t.Fatal(err)
	}
	got, err = loadState()
	if err != nil {
		t.Fatal(err)
	}
	if got != "johan@nordveda.se" {
		t.Errorf("loadState() = %q, want johan@nordveda.se", got)
	}

	if err := saveState(""); err != nil {
		t.Fatal(err)
	}
	if got, _ := loadState(); got != "" {
		t.Errorf("after clearing, loadState() = %q, want empty", got)
	}
}

func TestSendReplyRequiresConnection(t *testing.T) {
	s := NewMailService()
	if err := s.SendReply("a@b.se", "Re: hej", "kropp"); err == nil {
		t.Error("SendReply should fail when no account is connected")
	}
}

func TestAIAvailableReflectsClassifier(t *testing.T) {
	s := NewMailService()
	s.classifier = fakeClassifier{avail: true}
	if !s.AIAvailable() {
		t.Error("AIAvailable() = false, want true")
	}
	s.classifier = fakeClassifier{avail: false}
	if s.AIAvailable() {
		t.Error("AIAvailable() = true, want false")
	}
}

func TestGenerateDraft(t *testing.T) {
	want := mail.Draft{Tone: "varm", Lines: []mail.DraftLine{{Text: "Hej", Kind: mail.DraftAdded}}}
	s := NewMailService()
	s.classifier = fakeClassifier{avail: true, draft: want}

	mails := s.GetMails()
	if len(mails) == 0 {
		t.Fatal("no mock mails")
	}
	d, err := s.GenerateDraft(mails[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	if d.Tone != "varm" || len(d.Lines) != 1 || d.Lines[0].Text != "Hej" {
		t.Errorf("got %+v", d)
	}
}

func TestGenerateDraftUnknownMail(t *testing.T) {
	s := NewMailService()
	s.classifier = fakeClassifier{avail: true}
	if _, err := s.GenerateDraft("does-not-exist"); err == nil {
		t.Error("expected error for unknown mail id")
	}
}
