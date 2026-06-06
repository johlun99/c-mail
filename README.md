# cmail

A keyboard-first, **agentic email client** for the desktop (Linux + macOS).

cmail connects to your Gmail and uses a **local AI model** to categorize incoming
mail — cutting through noise and spam so nothing important slips through. A hard
product guarantee: **the agent never sends anything without your explicit approval.**

> Status: early development. The UI design lives in a separate design handoff; the
> application is being built from the ground up.

## Why local AI

cmail reads private email, so categorization and draft generation run **locally via
[Ollama](https://ollama.com)** by default — your mail data never leaves your machine.
The AI layer sits behind a pluggable interface, so other providers can be added later.

## Tech stack

- **[Wails v2](https://wails.io)** — desktop shell (Go + native webview)
- **Go** — native layer: Gmail API, OAuth, OS keychain, local SQLite cache, LLM proxy
- **React + TypeScript + Vite** — frontend UI
- **Ollama** — local AI for categorization & draft generation

There is no separate server; the Go side only handles the privileged operations a
browser sandbox can't.

## Development

Requires **Go**, **Node.js**, the **[Wails CLI](https://wails.io/docs/gettingstarted/installation)**,
**golangci-lint**, and **lefthook**. On Linux you also need `gtk3` + `webkit2gtk-4.1`.

Tools installed via `go install` (wails, lefthook) land in `$(go env GOPATH)/bin`
(usually `~/go/bin`) — **make sure that's on your `PATH`**, otherwise the git
pre-commit hook can't find them.

First-time setup:

```bash
cd frontend && npm install && cd ..   # frontend deps
make install-hooks                    # enable pre-commit lint+test hook
```

Everyday commands:

```bash
wails dev      # run the app in development with hot reload
wails build    # produce a distributable binary
make lint      # run all linters (Go + frontend)
make test      # run all tests (Go + frontend)
make fmt       # auto-format Go + frontend
```

## Contributing workflow

- **Branching:** `master` is the stable branch. All work happens on **feature
  branches** and lands via **pull request** — never commit directly to `master`.
- **Commits:** [Conventional Commits](https://www.conventionalcommits.org/)
  (`type(scope): message`), kept small and atomic.
- **Tests:** every feature or bugfix ships with unit tests (Go `*_test.go`,
  frontend Vitest), and existing tests are maintained as behaviour changes.
- **Quality gates:** linters and the test suite run **before every commit**
  (pre-commit hook) and again in **GitHub Actions on every PR**. Both must pass.

## License

[MIT](LICENSE) © 2026 Johan Lundgren
