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
from pathlib import Path

ROADMAP = Path("platform-roadmap/ROADMAP.md")
OUT = Path("progress-page/data.json")
MAX_COMMITS = 300  # enough for "today" counts and streaks

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
    out = git("log", f"-n{MAX_COMMITS}", "--no-merges", "--format=%h%x1f%cI%x1f%s")
    rows = [line.split("\x1f") for line in out.splitlines() if line]
    return [{"hash": h, "date": d, "message": s} for h, d, s in rows]


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
    return points


def main():
    text = ROADMAP.read_text(encoding="utf-8") if ROADMAP.exists() else ""
    data = {
        "percent": percent(text),
        "roadmap": parse_roadmap(text),
        "history": history(),
        "commits": commits(),
    }
    OUT.parent.mkdir(parents=True, exist_ok=True)
    OUT.write_text(json.dumps(data, indent=1), encoding="utf-8")
    print(f"wrote {OUT}: {data['percent']}% done, {len(data['commits'])} commits")


if __name__ == "__main__":
    main()
