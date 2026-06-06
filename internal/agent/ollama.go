package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"cmail/internal/mail"
)

const (
	defaultHost  = "http://localhost:11434"
	defaultModel = "qwen2.5:3b"
)

// Ollama is a Classifier backed by a local Ollama instance.
type Ollama struct {
	baseURL string
	model   string
	client  *http.Client
}

// NewOllama builds an Ollama classifier from OLLAMA_HOST / OLLAMA_MODEL (with
// local defaults).
func NewOllama() *Ollama {
	return NewOllamaWith(getenv("OLLAMA_HOST", defaultHost), getenv("OLLAMA_MODEL", defaultModel))
}

// NewOllamaWith builds an Ollama classifier against a specific host and model.
func NewOllamaWith(baseURL, model string) *Ollama {
	return &Ollama{
		baseURL: strings.TrimRight(baseURL, "/"),
		model:   model,
		client:  &http.Client{Timeout: 90 * time.Second},
	}
}

// Available reports whether the Ollama server responds.
func (o *Ollama) Available(ctx context.Context) bool {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, o.baseURL+"/api/tags", nil)
	if err != nil {
		return false
	}
	resp, err := o.client.Do(req)
	if err != nil {
		return false
	}
	defer func() { _ = resp.Body.Close() }()
	return resp.StatusCode == http.StatusOK
}

// Categorize asks the model to place the mail in exactly one category.
func (o *Ollama) Categorize(ctx context.Context, in MailInput, cats []mail.Category) (Classification, error) {
	var allowed []string
	known := make(map[string]bool, len(cats))
	for _, c := range cats {
		allowed = append(allowed, fmt.Sprintf("%s (%s)", c.Key, c.Label))
		known[c.Key] = true
	}
	system := "Du är en e-postassistent som sorterar inkommande mail. " +
		"Klassificera mailet i EXAKT en av dessa kategorinycklar: " + strings.Join(allowed, ", ") + ". " +
		"Svara ENDAST med JSON på formen " +
		`{"category":"<nyckel>","confidence":0.0-1.0,"summary":"en mening på svenska",` +
		`"facts":[{"key":"...","value":"..."}],"tasks":["..."]}.`
	user := fmt.Sprintf("Avsändare: %s\nÄmne: %s\n\n%s", in.From, in.Subject, in.Body)

	var out struct {
		Category   string      `json:"category"`
		Confidence float64     `json:"confidence"`
		Summary    string      `json:"summary"`
		Facts      []mail.Fact `json:"facts"`
		Tasks      []string    `json:"tasks"`
	}
	if err := o.chatJSON(ctx, system, user, &out); err != nil {
		return Classification{}, err
	}
	if !known[out.Category] && len(cats) > 0 {
		out.Category = cats[0].Key // fall back to the first category on a bad label
	}
	if out.Facts == nil {
		out.Facts = []mail.Fact{}
	}
	if out.Tasks == nil {
		out.Tasks = []string{}
	}
	return Classification(out), nil
}

// DraftReply writes a reply body in the given tone.
func (o *Ollama) DraftReply(ctx context.Context, in MailInput, tone string) (mail.Draft, error) {
	system := "Du skriver svar på e-post på svenska. Ton: " + tone + ". " +
		"Returnera ENDAST brödtexten i svaret, utan ämnesrad eller förklaringar."
	user := fmt.Sprintf("Svara på detta mail från %s.\nÄmne: %s\n\n%s", in.From, in.Subject, in.Body)

	text, err := o.chatText(ctx, system, user)
	if err != nil {
		return mail.Draft{}, err
	}
	lines := strings.Split(strings.TrimSpace(text), "\n")
	out := make([]mail.DraftLine, len(lines))
	for i, l := range lines {
		out[i] = mail.DraftLine{Text: l, Kind: mail.DraftAdded}
	}
	return mail.Draft{Tone: tone, Lines: out}, nil
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
	Stream   bool          `json:"stream"`
	Format   string        `json:"format,omitempty"`
}

type chatResponse struct {
	Message chatMessage `json:"message"`
}

func (o *Ollama) chat(ctx context.Context, system, user, format string) (string, error) {
	body, err := json.Marshal(chatRequest{
		Model: o.model,
		Messages: []chatMessage{
			{Role: "system", Content: system},
			{Role: "user", Content: user},
		},
		Stream: false,
		Format: format,
	})
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, o.baseURL+"/api/chat", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := o.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ollama: chat: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama: chat status %d", resp.StatusCode)
	}
	var cr chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&cr); err != nil {
		return "", fmt.Errorf("ollama: decode: %w", err)
	}
	return cr.Message.Content, nil
}

func (o *Ollama) chatText(ctx context.Context, system, user string) (string, error) {
	return o.chat(ctx, system, user, "")
}

func (o *Ollama) chatJSON(ctx context.Context, system, user string, target any) error {
	content, err := o.chat(ctx, system, user, "json")
	if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(content), target); err != nil {
		return fmt.Errorf("ollama: parse json reply: %w", err)
	}
	return nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
