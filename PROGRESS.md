# cmail — Progress

Persistent checklist of the project's phases and tasks. Tick items off as they land
so we always know where we are — even after a context reset. Update at the end of
every PR/phase. Full rationale lives in the approved plan.

**Legend:** `[ ]` todo · `[~]` in progress · `[x]` done

---

## Phase 0 — Git bootstrap & repo hygiene  *(on master)*  ✅
- [x] Extend `.gitignore` (CLAUDE.md, design handoff, Go/Node/Wails artifacts)
- [x] Add `LICENSE` (MIT, Johan Lundgren, 2026)
- [x] Add `README.md` (overview, stack, dev setup, workflow)
- [x] Add `PROGRESS.md` (this file)
- [x] Update `CLAUDE.md` (structure + decided rules; stays gitignored/local)
- [x] Sync with existing remote commit, commit, push to master

## Phase 1 — Wails scaffold + tooling + CI  *(feat/scaffold → PR)*
- [x] `wails init -n cmail -t react-ts` (moved into repo root)
- [x] Go tooling: golangci-lint (v2), gofmt/goimports, `go test`
- [x] Frontend tooling: ESLint 9, Prettier 3, `tsc --noEmit`, Vitest 2 (modernized: Vite 6, TS 5)
- [x] lefthook pre-commit (runs lint + test for both sides)
- [x] Makefile unifying `lint`, `test`, `build`, `dev`, `fmt`, `install-hooks`
- [x] GitHub Actions CI on PR (frontend + go + wails smoke build)
- [x] Local verification: make lint + make test green; frontend + go build pass; pre-commit blocks broken code
- [ ] Verify CI runs green on the PR
- [ ] `gh auth login` + (optional) branch protection on master

## Phase 2 — Domain model + mock data  *(feat/mock-data)*  ✅
- [x] Go domain types in `internal/mail` (Mail, Category, AgentAnalysis, Fact, Draft, DraftLine, Activity, Rule, Account)
- [x] Default category taxonomy (extension via settings UI comes in Phase 3)
- [x] Mock provider behind a `Store` interface (full prototype dataset ported, returns copies)
- [x] Exposed via Wails `MailService` bindings; TS models generated (`wails generate module`)
- [x] Go tests on the store + frontend test guarding the TS↔Go contract
- [ ] Live fetch in UI deferred to Phase 3 (needs `wails dev` + webkit)

## Phase 3 — UI matching the design  *(feat/ui, split into increments)*
Increment 1 — inbox shell ✅
- [x] Port design tokens to `style.css` (OKLCH ramp, `--bgl`, ambient depth, scanlines)
- [x] `useTweaks` hook applies accent/font/density/layout/`--bgl` to `:root` (persisted)
- [x] Components: TopBar, Rail, MessageList (+Row, day groups, manual scroll), ReadingPane (agent block + draft diff + sent banner), StatusBar
- [x] State model + vim keymap (j/k/g/G/1-6/Enter/l/r/e/a/d/c/// /Esc, Mod+Enter)
- [x] Wired to mock data via `MailService` bindings (graceful empty outside Wails)
- [x] Keyboard + mouse parity for built surfaces; no `scrollIntoView`; platform-aware Mod key
- [x] Swedish copy; tabular numerals; depth preserved; frontend test for shell

Increment 2 — overlays + tweak UI ✅
- [x] Command palette (`Mod+K`) with grouped commands + filtering + keyboard nav
- [x] Settings overlay (TUI nav: konton, utseende, agent, kategorier, röst, sekretess, om)
- [x] Help/keymap overlay (`?`)
- [x] Wired accent/density/font/columns/bgl controls to `setTweak`; agent rules toggle
- [ ] Live visual check via `wails dev` (needs webkit) against screenshots

## Phase 4 — Gmail integration  *(feat/gmail)*
- [ ] OAuth2 (installed-app/PKCE), opens system browser, captures redirect
- [ ] Token storage in OS keychain (go-keyring)
- [ ] Gmail API: messages, threads, labels, drafts (never auto-send)
- [ ] Local SQLite cache + background sync
- [ ] Account flow in settings (disconnected → connecting → connected)

## Phase 5 — Local AI agent flow (Ollama)  *(feat/agent)*
- [ ] Pluggable `Classifier` interface in Go
- [ ] Ollama implementation (localhost:11434); model choice per RAM
- [ ] Detect if Ollama is running; guide setup if not
- [ ] Categorization writes AgentAnalysis + optional auto-draft per category
- [ ] Never-send locked in settings
- [ ] Verify no mail data leaves the machine

---

## Open questions (non-blocking)
- [ ] Confirm exact LICENSE copyright name
- [ ] Laptop RAM/GPU → Ollama model size (Phase 5)
- [ ] Branch protection rules on master (gh CLI)
