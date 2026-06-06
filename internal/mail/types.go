// Package mail defines cmail's core domain types and the Store interface that
// supplies mail data to the app. The mock implementation mirrors the design
// prototype's dataset; a Gmail-backed implementation will replace it later
// behind the same interface.
package mail

// Category is a triage bucket. Key is the stable identifier; Color is a design
// token name (resolved to a real color in the frontend CSS). The taxonomy is
// user-extensible — additional categories can be added at runtime.
type Category struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Short string `json:"short"`
	Color string `json:"color"`
	Desc  string `json:"desc"`
}

// Fact is a single key/value row in the agent's inline analysis tree. A value
// containing "⚠" is rendered as a warning in the UI.
type Fact struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// AgentAnalysis is the agent's read of a mail: a one-line summary, a list of
// facts, and suggested follow-up tasks.
type AgentAnalysis struct {
	Summary string   `json:"summary"`
	Facts   []Fact   `json:"facts"`
	Tasks   []string `json:"tasks"`
}

// DraftLineKind marks how a draft line was authored.
type DraftLineKind string

const (
	// DraftAdded is an agent-authored line (rendered with a "+").
	DraftAdded DraftLineKind = "+"
	// DraftModified is a line the user edited (rendered with a "~").
	DraftModified DraftLineKind = "~"
)

// DraftLine is one line of a draft reply, with its authorship kind for the
// diff-style approval view.
type DraftLine struct {
	Text string        `json:"text"`
	Kind DraftLineKind `json:"kind"`
}

// Draft is an agent-generated reply awaiting the user's approval. The agent
// never sends — drafts are only sent after explicit user approval.
type Draft struct {
	Tone  string      `json:"tone"`
	Lines []DraftLine `json:"lines"`
}

// Mail is a single message (the head of a thread) with the agent's analysis and
// an optional pending draft reply.
type Mail struct {
	ID          string        `json:"id"`
	Cat         string        `json:"cat"`
	From        string        `json:"from"`
	FromAddr    string        `json:"fromAddr"`
	Avatar      string        `json:"avatar"`
	Subject     string        `json:"subject"`
	Snippet     string        `json:"snippet"`
	Time        string        `json:"time"`
	Day         string        `json:"day"`
	Unread      bool          `json:"unread"`
	Deadline    string        `json:"deadline,omitempty"`
	Confidence  float64       `json:"confidence"`
	ThreadCount int           `json:"threadCount"`
	Body        []string      `json:"body"`
	Agent       AgentAnalysis `json:"agent"`
	Draft       *Draft        `json:"draft"`
}

// Activity is one entry in the agent activity feed (status bar / palette context).
type Activity struct {
	Time string `json:"time"`
	Text string `json:"text"`
	Cat  string `json:"cat,omitempty"`
}

// Rule is an agent automation toggle shown in settings. Locked rules cannot be
// turned off (e.g. the never-send-without-approval guarantee).
type Rule struct {
	On     bool   `json:"on"`
	Text   string `json:"text"`
	Scope  string `json:"scope"`
	Locked bool   `json:"locked,omitempty"`
}

// AccountStatus is the connection state of a linked mail account.
type AccountStatus string

const (
	StatusConnected    AccountStatus = "connected"
	StatusConnecting   AccountStatus = "connecting"
	StatusDisconnected AccountStatus = "disconnected"
)

// Account is a linked mail account shown in settings.
type Account struct {
	Email    string        `json:"email"`
	Provider string        `json:"provider"`
	Status   AccountStatus `json:"status"`
	SyncedAt string        `json:"syncedAt"`
	Scopes   []string      `json:"scopes"`
}
