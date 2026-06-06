package main

import "testing"

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
