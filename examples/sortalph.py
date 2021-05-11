#!/usr/bin/env python3

# sort images inside each group alphabetically
# USAGE:
# just pipe json into it

import json
import sys

groups = json.load(sys.stdin)

# no input
if groups is None:
    sys.exit(0)

for i, group in enumerate(groups):
    groups[i] = sorted(group)

# if you want sorted json back
# print(json.dumps(groups, indent=2))

# flat list
groups = sorted(groups)
for group in groups:
    for fp in group:
        print(fp)
