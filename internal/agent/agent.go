// Package agent runs the local AI flow: classifying mail into categories and
// drafting replies. It is pluggable behind the Classifier interface; the default
// implementation talks to a local Ollama instance so mail content never leaves
// the machine.
package agent

import (
	"context"

	"cmail/internal/mail"
)

// MailInput is the subset of a mail handed to the model.
type MailInput struct {
	From    string
	Subject string
	Body    string
}

// Classification is the model's read of a mail.
type Classification struct {
	Category   string
	Confidence float64
	Summary    string
	Facts      []mail.Fact
	Tasks      []string
}

// Classifier categorizes mail and drafts replies. Implementations must never
// send — drafting only.
type Classifier interface {
	// Available reports whether the backend is reachable.
	Available(ctx context.Context) bool
	// Categorize assigns one of cats to the mail, with a summary and facts.
	Categorize(ctx context.Context, in MailInput, cats []mail.Category) (Classification, error)
	// DraftReply writes a reply body in the requested tone.
	DraftReply(ctx context.Context, in MailInput, tone string) (mail.Draft, error)
}
