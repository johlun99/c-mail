package gmail

import (
	"encoding/base64"
	"net/mail"
	"strings"
	"time"

	gmailapi "google.golang.org/api/gmail/v1"

	cmail "cmail/internal/mail"
)

// mapMessage converts a fully-fetched Gmail message into cmail's domain model.
// Agent analysis, category confidence and drafts are left empty/zero — those are
// produced by the local AI layer (Phase 5). A coarse category is derived from
// Gmail's own labels as a stopgap until then.
func mapMessage(m *gmailapi.Message, now time.Time) cmail.Mail {
	h := headers(m)
	name, addr := parseFrom(h["from"])
	ts := time.UnixMilli(m.InternalDate)

	return cmail.Mail{
		ID:          m.Id,
		Cat:         categoryFromLabels(m.LabelIds),
		From:        name,
		FromAddr:    addr,
		Avatar:      initials(name, addr),
		Subject:     h["subject"],
		Snippet:     m.Snippet,
		Time:        ts.Format("15:04"),
		Day:         dayLabel(ts, now),
		Unread:      contains(m.LabelIds, "UNREAD"),
		ThreadCount: 1,
		Body:        extractBody(m.Payload),
		Agent:       cmail.AgentAnalysis{Facts: []cmail.Fact{}, Tasks: []string{}},
	}
}

func headers(m *gmailapi.Message) map[string]string {
	out := map[string]string{}
	if m.Payload == nil {
		return out
	}
	for _, hdr := range m.Payload.Headers {
		out[strings.ToLower(hdr.Name)] = hdr.Value
	}
	return out
}

func parseFrom(raw string) (name, addr string) {
	if raw == "" {
		return "", ""
	}
	if a, err := mail.ParseAddress(raw); err == nil {
		name = a.Name
		addr = a.Address
		if name == "" {
			name = addr
		}
		return name, addr
	}
	return raw, raw
}

func initials(name, addr string) string {
	src := strings.TrimSpace(name)
	if src == "" {
		src = addr
	}
	fields := strings.Fields(src)
	switch {
	case len(fields) >= 2:
		return strings.ToUpper(string([]rune(fields[0])[:1]) + string([]rune(fields[1])[:1]))
	case len(fields) == 1 && len([]rune(fields[0])) >= 2:
		return strings.ToUpper(string([]rune(fields[0])[:2]))
	case len(fields) == 1:
		return strings.ToUpper(fields[0])
	default:
		return "?"
	}
}

func dayLabel(t, now time.Time) string {
	y, m, d := now.Date()
	today := time.Date(y, m, d, 0, 0, 0, 0, now.Location())
	switch {
	case !t.Before(today):
		return "idag"
	case !t.Before(today.AddDate(0, 0, -1)):
		return "igår"
	default:
		return "tidigare"
	}
}

// categoryFromLabels is a coarse stopgap mapping from Gmail labels to cmail
// categories, replaced by the local AI classifier in Phase 5.
func categoryFromLabels(labels []string) string {
	switch {
	case contains(labels, "CATEGORY_PROMOTIONS"):
		return "reklam"
	case contains(labels, "IMPORTANT"), contains(labels, "STARRED"):
		return "svara"
	case contains(labels, "CATEGORY_FORUMS"), contains(labels, "CATEGORY_SOCIAL"),
		contains(labels, "CATEGORY_UPDATES"):
		return "intressant"
	default:
		return "intressant"
	}
}

func extractBody(p *gmailapi.MessagePart) []string {
	raw := findPlainText(p)
	if raw == "" {
		return nil
	}
	parts := strings.Split(strings.ReplaceAll(raw, "\r\n", "\n"), "\n\n")
	out := make([]string, 0, len(parts))
	for _, para := range parts {
		if s := strings.TrimSpace(para); s != "" {
			out = append(out, s)
		}
	}
	return out
}

func findPlainText(p *gmailapi.MessagePart) string {
	if p == nil {
		return ""
	}
	if p.MimeType == "text/plain" && p.Body != nil && p.Body.Data != "" {
		if data, err := base64.URLEncoding.DecodeString(p.Body.Data); err == nil {
			return string(data)
		}
	}
	for _, child := range p.Parts {
		if s := findPlainText(child); s != "" {
			return s
		}
	}
	return ""
}

func contains(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}
