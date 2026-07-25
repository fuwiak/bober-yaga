#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
git config core.hooksPath .githooks
git config commit.template .gitmessage
chmod +x .githooks/commit-msg .githooks/prepare-commit-msg
echo "ok — core.hooksPath=.githooks  commit.template=.gitmessage"
echo "commits need prefixes like: ENH: …  BUG: …  CI: …"
echo "Cursor/Codex attribution in messages is rejected"
