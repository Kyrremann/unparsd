#!/usr/bin/env python3
"""Split untappd.json into one per-year file (e.g. checkins/2016.json)."""

import json
import os
import sys

input_file = sys.argv[1] if len(sys.argv) > 1 else "untappd.json"
output_dir = sys.argv[2] if len(sys.argv) > 2 else "checkins"

with open(input_file) as f:
    checkins = json.load(f)

by_year = {}
for c in checkins:
    year = c["created_at"][:4]
    by_year.setdefault(year, []).append(c)

os.makedirs(output_dir, exist_ok=True)

for year, items in sorted(by_year.items()):
    # Sort newest-first to match the fetch command's output format.
    items.sort(key=lambda c: c["created_at"], reverse=True)
    path = os.path.join(output_dir, f"{year}.json")
    with open(path, "w") as f:
        json.dump(items, f, indent=2, ensure_ascii=False)
    print(f"Wrote {len(items):4d} check-ins to {path}")
