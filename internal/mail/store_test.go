package mail

import "testing"

func TestMockStoreNonEmpty(t *testing.T) {
	s := NewMockStore()
	if len(s.Mails()) == 0 {
		t.Error("expected mails, got none")
	}
	if got := len(s.Categories()); got != 6 {
		t.Errorf("Categories() = %d, want 6", got)
	}
	if len(s.Rules()) == 0 {
		t.Error("expected rules, got none")
	}
	if len(s.Accounts()) == 0 {
		t.Error("expected accounts, got none")
	}
	if len(s.Activity()) == 0 {
		t.Error("expected activity, got none")
	}
}

func TestEveryMailHasKnownCategory(t *testing.T) {
	s := NewMockStore()
	known := make(map[string]bool)
	for _, c := range s.Categories() {
		known[c.Key] = true
	}
	for _, m := range s.Mails() {
		if !known[m.Cat] {
			t.Errorf("mail %s has unknown category %q", m.ID, m.Cat)
		}
	}
}

func TestMailIDsAreUnique(t *testing.T) {
	s := NewMockStore()
	seen := make(map[string]bool)
	for _, m := range s.Mails() {
		if seen[m.ID] {
			t.Errorf("duplicate mail id %q", m.ID)
		}
		seen[m.ID] = true
	}
}

func TestReturnedSlicesAreCopies(t *testing.T) {
	s := NewMockStore()
	mails := s.Mails()
	if len(mails) == 0 {
		t.Fatal("no mails")
	}
	mails[0].Subject = "mutated"
	if s.Mails()[0].Subject == "mutated" {
		t.Error("mutating returned slice leaked into the store")
	}
}

func TestNeverSendRuleIsLockedAndOn(t *testing.T) {
	s := NewMockStore()
	var found bool
	for _, r := range s.Rules() {
		if r.Scope == "global" && r.Locked {
			found = true
			if !r.On {
				t.Error("the never-send rule must be on")
			}
		}
	}
	if !found {
		t.Error("expected a locked global rule (never send without approval)")
	}
}

func TestDraftsContainOnlyAddedLines(t *testing.T) {
	s := NewMockStore()
	for _, m := range s.Mails() {
		if m.Draft == nil {
			continue
		}
		if len(m.Draft.Lines) == 0 {
			t.Errorf("mail %s has an empty draft", m.ID)
		}
		for _, l := range m.Draft.Lines {
			if l.Kind != DraftAdded {
				t.Errorf("mail %s draft has non-added line kind %q", m.ID, l.Kind)
			}
		}
	}
}
