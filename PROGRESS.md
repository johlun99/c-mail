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
- [ ] `wails init -n cmail -t react-ts`
- [ ] Go tooling: golangci-lint, gofmt/goimports, `go test`
- [ ] Frontend tooling: ESLint, Prettier, `tsc --noEmit`, Vitest
- [ ] lefthook pre-commit (runs lint + test for both sides)
- [ ] Makefile / npm scripts unifying `lint`, `test`, `build`
- [ ] GitHub Actions CI on PR (lint + test + smoke build)
- [ ] Verify CI runs green on the PR; pre-commit blocks a broken commit
- [ ] `gh auth login` + (optional) branch protection on master

## Phase 2 — Domain model + mock data  *(feat/mock-data)*
- [ ] Go domain types (Mail, Category, AgentAnalysis, Draft, Account)
- [ ] Configurable/extensible category taxonomy
- [ ] Mock provider behind a `MailStore` interface (port prototype dataset)
- [ ] Expose via Wails bindings; mirrored TS types in frontend
- [ ] Tests on the provider

## Phase 3 — UI matching the design  *(feat/ui, maybe split)*
- [ ] Port design tokens (OKLCH ramp, `--bgl`, ambient depth)
- [ ] Components: topbar, rail, message list, reading pane (agent + draft diff)
- [ ] Overlays: command palette, settings, help/keymap
- [ ] State model + vim-like keymap
- [ ] Keyboard + mouse parity; no `scrollIntoView`; platform-aware Mod key
- [ ] Swedish copy; tabular numerals; preserve depth (not flat dark)

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
