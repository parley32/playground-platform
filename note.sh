#!/usr/bin/env bash
# Write the note shown at the top of the progress page, then commit and push it.
#   ./note.sh "Motivation 7/10. Today: finish the TCP server timeouts."
#   ./note.sh            (opens your editor for a longer, multi-line note)
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"
NOTE=progress-page/note.md
if [ $# -gt 0 ]; then
  printf '%s\n' "$*" > "$NOTE"
else
  "${EDITOR:-nano}" "$NOTE"
fi
git add "$NOTE"
git commit -q -m "note: $(head -n1 "$NOTE" | cut -c1-60)"
git push -q
echo "Note is on its way to the page (about a minute)."
