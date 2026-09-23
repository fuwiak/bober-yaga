# Commit messages

## Format

```
TYPE: short imperative summary
```

Use **pandas-style** prefixes (same idea as their history: `BUG:`, `BLD:`, `PERF:`, …):

| Prefix | When |
|--------|------|
| `ENH:` | New feature or brick |
| `BUG:` | Bug fix |
| `PERF:` | Performance |
| `REF:` | Refactor, no intended behavior change |
| `BLD:` | Build, `go.mod`, install scripts |
| `CI:` | GitHub Actions, hooks |
| `DOC:` | Docs / README |
| `TST:` | Tests |
| `STY:` | Formatting only |
| `API:` | Public CLI/API change |
| `DIST:` | Packaging / release artifacts |
| `MAINT:` | Chores, dependency bumps |
| `SEC:` | Security fixes |

Examples:

```
ENH: add Audience brick for segment sync
BUG: fix credentials path when HOME unset
CI: run go test on pull requests
BLD: retarget module path to github.com/fuwiak/bober-yaga
```

## Do not put in commits

- Mentions of **Cursor**, **Codex**, Copilot, ChatGPT, or “AI generated”
- `Co-authored-by:` trailers for tools/agents

Local hooks enforce this (`.githooks/`). Enable once:

```bash
./.githooks/install.sh
```

## Setup after clone

```bash
git clone https://github.com/fuwiak/bober-yaga.git
cd bober-yaga
./.githooks/install.sh
go build -o yaga .
./install.sh   # optional: ~/bin/yaga
```

Script bricks that call into **bober-ai** `scripts/*.mjs` need:

```bash
export YAGA_REPO=/path/to/bober-ai
```
