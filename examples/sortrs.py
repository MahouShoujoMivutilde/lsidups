#!/usr/bin/env python3

# sort images inside each group by quality, depends on Pillow
# USAGE:
# just pipe json into it

import json
import sys
from os import path

from PIL import Image

groups = json.load(sys.stdin)

# no input
if groups is None:
    sys.exit(0)


def quality(fp):
    with Image.open(fp) as img:
        return img.width * img.height * path.getsize(fp)


for i, group in enumerate(groups):
    groups[i] = sorted(group, key=quality, reverse=True)

# if you want sorted json back
# print(json.dumps(groups, indent=2))

# flat list
for group in groups:
    for fp in group:
        print(fp)
