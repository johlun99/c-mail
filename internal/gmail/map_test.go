package gmail

import (
	"encoding/base64"
	"testing"
	"time"

	gmailapi "google.golang.org/api/gmail/v1"
)

func TestMapMessage(t *testing.T) {
	now := time.Date(2026, 6, 6, 12, 0, 0, 0, time.UTC)
	body := base64.URLEncoding.EncodeToString([]byte("Hej!\r\n\r\nAndra stycket."))
	msg := &gmailapi.Message{
		Id:           "abc",
		Snippet:      "Hej! Andra…",
		InternalDate: now.Add(-2 * time.Hour).UnixMilli(),
		LabelIds:     []string{"INBOX", "UNREAD", "IMPORTANT"},
		Payload: &gmailapi.MessagePart{
			MimeType: "text/plain",
			Headers: []*gmailapi.MessagePartHeader{
				{Name: "From", Value: "Anna Lindqvist <anna@nordveda.se>"},
				{Name: "Subject", Value: "Offert"},
			},
			Body: &gmailapi.MessagePartBody{Data: body},
		},
	}

	m := mapMessage(msg, now)
	if m.ID != "abc" {
		t.Errorf("ID = %q", m.ID)
	}
	if m.From != "Anna Lindqvist" || m.FromAddr != "anna@nordveda.se" {
		t.Errorf("from parsed wrong: %q / %q", m.From, m.FromAddr)
	}
	if m.Avatar != "AL" {
		t.Errorf("avatar = %q, want AL", m.Avatar)
	}
	if m.Subject != "Offert" {
		t.Errorf("subject = %q", m.Subject)
	}
	if !m.Unread {
		t.Error("expected unread")
	}
	if m.Cat != "svara" {
		t.Errorf("cat = %q, want svara (IMPORTANT)", m.Cat)
	}
	if m.Day != "idag" {
		t.Errorf("day = %q, want idag", m.Day)
	}
	if len(m.Body) != 2 || m.Body[0] != "Hej!" || m.Body[1] != "Andra stycket." {
		t.Errorf("body = %#v", m.Body)
	}
}

func TestDayLabel(t *testing.T) {
	now := time.Date(2026, 6, 6, 12, 0, 0, 0, time.UTC)
	cases := []struct {
		offset time.Duration
		want   string
	}{
		{-1 * time.Hour, "idag"},
		{-26 * time.Hour, "igår"},
		{-72 * time.Hour, "tidigare"},
	}
	for _, c := range cases {
		if got := dayLabel(now.Add(c.offset), now); got != c.want {
			t.Errorf("dayLabel(%v) = %q, want %q", c.offset, got, c.want)
		}
	}
}

func TestCategoryFromLabels(t *testing.T) {
	cases := map[string]string{
		"reklam":     "CATEGORY_PROMOTIONS",
		"svara":      "STARRED",
		"intressant": "CATEGORY_SOCIAL",
	}
	for want, label := range cases {
		if got := categoryFromLabels([]string{"INBOX", label}); got != want {
			t.Errorf("label %s → %q, want %q", label, got, want)
		}
	}
	if got := categoryFromLabels([]string{"INBOX"}); got != "intressant" {
		t.Errorf("default = %q, want intressant", got)
	}
}

func TestBuildMIME(t *testing.T) {
	got := buildMIME("a@b.se", "Hej", "kropp")
	if want := "To: a@b.se\r\nSubject: Hej\r\n"; got[:len(want)] != want {
		t.Errorf("MIME prefix = %q", got[:len(want)])
	}
}
