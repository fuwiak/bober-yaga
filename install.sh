#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")" && pwd)"
BIN_DIR="${YAGA_BIN_DIR:-$HOME/bin}"
mkdir -p "$BIN_DIR"
cd "$ROOT"
go build -o yaga .
chmod +x "$ROOT/run"
ln -sfn "$ROOT/run" "$BIN_DIR/yaga"
echo "installed: $BIN_DIR/yaga → $ROOT/run (Go + Bubble Tea)"
yaga version
yaga bricks list | head -20
echo "ok — yaga  |  yaga help  |  yaga profile public"
