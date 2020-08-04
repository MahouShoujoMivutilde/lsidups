#!/usr/bin/env python3

# sort images inside each group by resolution, depends on Pillow
# USAGE:
# just pipe json into it

import json
import sys

from PIL import Image

groups = json.load(sys.stdin)


def resolution(fp):
    with Image.open(fp) as img:
        return img.width * img.height


for i, group in enumerate(groups):
    groups[i] = sorted(group, key=resolution, reverse=True)

# if you want sorted json back
# print(json.dumps(groups, indent=2))

# flat list
for group in groups:
    for fp in group:
        print(fp)
