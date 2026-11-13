#!/usr/bin/env python3
"""Build site/data.json from ROADMAP.md and the git history.

Run from the repo root (the GitHub Action does this on every push):
    python3 progress-page/build_data.py

Only "## Phase ..." sections count towards progress; "Ongoing" and the
other sections are left out.
"""
import json
import re
import subprocess
from datetime import datetime
from pathlib import Path
from zoneinfo import ZoneInfo

ROADMAP = Path("platform-roadmap/ROADMAP.md")
OUT = Path("progress-page/data.json")
NOTE = Path("progress-page/note.md")  # the "start of the day" note shown on the page
MAX_COMMITS = 300  # enough for "today" counts and streaks
# Commits and the progress curve only count from this day (Toulouse time).
TRACK_FROM = datetime(2026, 11, 16, tzinfo=ZoneInfo("Europe/Paris"))

TASK = re.compile(r"^\s*[-*]\s+\[( |x|X)\]\s+(.*)$")
HEADING = re.compile(r"^##\s+(.*)$")


def git(*args):
    return subprocess.run(["git", *args], capture_output=True, text=True, check=True).stdout


def parse_roadmap(text):
    """Phase sections with their checkbox items."""
    sections, current = [], None
    for line in text.splitlines():
        if m := HEADING.match(line):
            title = m.group(1).strip()
            if title.startswith("Phase"):
                phase, _, name = title.partition(" — ")
                current = {"phase": phase.strip(), "title": name.strip(), "items": []}
                sections.append(current)
            else:
                current = None
        elif current is not None and (m := TASK.match(line)):
            label, _, note = m.group(2).partition(" — ")
            current["items"].append(
                {"done": m.group(1).lower() == "x", "label": label.strip(), "note": note.strip().strip("_")}
            )
    return [s for s in sections if s["items"]]


def percent(text):
    items = [i for s in parse_roadmap(text) for i in s["items"]]
    if not items:
        return 0.0
    return round(100 * sum(i["done"] for i in items) / len(items), 1)


def commits():
    """Recent commits since TRACK_FROM, leaving out commits that only touch the note."""
    out = git("log", f"-n{MAX_COMMITS}", "--no-merges", "--format=%x1e%h%x1f%cI%x1f%s", "--name-only")
    result = []
    for record in out.split("\x1e")[1:]:
        header, _, files = record.partition("\n")
        h, d, s = header.split("\x1f")
        changed = {f for f in files.splitlines() if f.strip()}
        if changed == {str(NOTE)}:
            continue
        if datetime.fromisoformat(d) >= TRACK_FROM:
            result.append({"hash": h, "date": d, "message": s})
    return result


def note():
    """The note's text and when it was last committed (None if there is no note)."""
    if not NOTE.exists() or not NOTE.read_text(encoding="utf-8").strip():
        return None
    date = git("log", "-1", "--format=%cI", "--", str(NOTE)).strip()
    return {"text": NOTE.read_text(encoding="utf-8").strip(), "date": date or None}


def history():
    """% complete after every commit that touched ROADMAP.md (oldest first)."""
    out = git("log", "--reverse", "--follow", "--format=%H%x1f%cI", "--", str(ROADMAP))
    points = []
    for line in out.splitlines():
        sha, date = line.split("\x1f")
        try:
            text = git("show", f"{sha}:{ROADMAP}")
        except subprocess.CalledProcessError:
            continue  # file deleted/renamed in that commit
        points.append({"date": date, "percent": percent(text)})
    # Everything before TRACK_FROM collapses into one starting point.
    before = [pt for pt in points if datetime.fromisoformat(pt["date"]) < TRACK_FROM]
    after = [pt for pt in points if datetime.fromisoformat(pt["date"]) >= TRACK_FROM]
    start = {"date": TRACK_FROM.isoformat(), "percent": before[-1]["percent"] if before else 0.0}
    return [start] + after


def main():
    text = ROADMAP.read_text(encoding="utf-8") if ROADMAP.exists() else ""
    data = {
        "percent": percent(text),
        "roadmap": parse_roadmap(text),
        "trackFrom": TRACK_FROM.isoformat(),
        "history": history(),
        "commits": commits(),
        "note": note(),
    }
    OUT.parent.mkdir(parents=True, exist_ok=True)
    OUT.write_text(json.dumps(data, indent=1), encoding="utf-8")
    print(f"wrote {OUT}: {data['percent']}% done, {len(data['commits'])} commits")


if __name__ == "__main__":
    main()
