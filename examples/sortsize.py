#!/usr/bin/env python3

# sort images inside each group by file size
# USAGE:
# just pipe json into it

import json
import sys
from os import path

groups = json.load(sys.stdin)

for i, group in enumerate(groups):
    groups[i] = sorted(group, key=lambda fp: path.getsize(fp), reverse=True)

# if you want sorted json back
# print(json.dumps(groups, indent=2))

# flat list
for group in groups:
    for fp in group:
        print(fp)
