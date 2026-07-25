# yaga

**yaga** (Яга) — modular Yandex ops CLI (Go + Charm Bubble Tea).

Each service is a **brick** in `bricks.go`. Profiles hide private bricks.

> Developed here: [`github.com/fuwiak/yaga`](https://github.com/fuwiak/yaga).  
> Script bricks that call `scripts/*.mjs` expect a **bober-ai** checkout via `YAGA_REPO`.

## Install

```bash
git clone https://github.com/fuwiak/yaga.git
cd yaga
go build -o yaga .
./install.sh          # ~/bin/yaga → ./run
./.githooks/install.sh  # commit style hooks (contributors)
```

## Usage

```bash
yaga                         # TUI
yaga webmaster status
yaga webmaster oauth
yaga webmaster seo
yaga metrika status
yaga direct campaigns status
yaga bricks
yaga profile public
yaga doctor
yaga credentials
```

Point at bober-ai scripts when needed:

```bash
export YAGA_REPO=/path/to/bober-ai
```

## Commit style

Pandas-style prefixes (`ENH:`, `BUG:`, `CI:`, …). No Cursor/Codex attribution.

See [CONTRIBUTING.md](CONTRIBUTING.md).

## CI / CD

- **CI** — `go vet`, build, test, smoke on push/PR; PR commits checked for `TYPE:` prefix
- **Release** — tag `v*` → multi-arch binaries on GitHub Releases

```bash
git tag v0.2.1
git push origin v0.2.1
```

## Stack

- Go 1.22+
- bubbletea + lipgloss + bubbles
