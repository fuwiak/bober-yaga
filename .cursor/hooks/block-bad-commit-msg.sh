#!/usr/bin/env bash
# Cursor hook: warn/block git commit if -m lacks TYPE: or mentions tools.
set -euo pipefail
input="$(cat)"
cmd="$(printf '%s' "$input" | sed -n 's/.*"command"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' | head -1)"
# Fail open if we cannot parse
if [[ -z "$cmd" ]]; then
  echo '{"permission":"allow"}'
  exit 0
fi

if ! [[ "$cmd" =~ git[[:space:]]+commit ]]; then
  echo '{"permission":"allow"}'
  exit 0
fi

# Extract -m / --message if present
msg=""
if [[ "$cmd" =~ --message[=[:space:]]+\"([^\"]+)\" ]]; then
  msg="${BASH_REMATCH[1]}"
elif [[ "$cmd" =~ -m[[:space:]]+\"([^\"]+)\" ]]; then
  msg="${BASH_REMATCH[1]}"
elif [[ "$cmd" =~ --message[=[:space:]]+\'([^\']+)\' ]]; then
  msg="${BASH_REMATCH[1]}"
elif [[ "$cmd" =~ -m[[:space:]]+\'([^\']+)\' ]]; then
  msg="${BASH_REMATCH[1]}"
fi

if [[ -n "$msg" ]]; then
  if echo "$msg" | grep -Eiq 'cursor|codex|copilot|chatgpt|co-authored-by'; then
    echo '{"permission":"deny","user_message":"Commit message must not mention Cursor/Codex/Copilot. Use ENH:/BUG:/CI: …"}'
    exit 0
  fi
  if ! [[ "$msg" =~ ^(ENH|BUG|PERF|REF|BLD|CI|DOC|TST|STY|API|DIST|MAINT|WEB|TYP|SEC|Merge|Revert|Release): ]]; then
    echo '{"permission":"deny","user_message":"Use a pandas-style prefix, e.g. ENH: add brick X"}'
    exit 0
  fi
fi

echo '{"permission":"allow"}'
