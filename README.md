# lsidups

### ...is a barebone tool for finding image duplicates (or just similar images) from your terminal.

### How to use
Pipe a list of files to compare into stdin or just input (`-i`) a directory you want to check. It then will output images grouped by similarity so you can process them as you please.

It uses library [images](https://github.com/vitali-fedulov/images) (MIT) under the hood, so if you want to know more about limitations and how comparison works - read [this](https://similar.pictures/algorithm-for-perceptual-image-comparison.html).

lsidups itself is just a wrapper that tries to provide a way to compare a lot (10k+) images reasonably fast from cli, and then allow you to process found duplicates in some other more convenient tool, like e.g. [sxiv](https://github.com/muennich/sxiv).

### Image formats support

At the moment of writing, it supports **only** **jpeg**, **png** and **gif**; i tried to make webp work, but [this](https://github.com/golang/go/issues/38341) prevented it. _\*very sad webp UwU\*._

### Install

Make sure you have go and git installed.

```sh
go get github.com/MahouShoujoMivutilde/lsidups
```

## Options

```
Usage of lsidups:
  -T int
        number of processing threads (default number of logical cores)

  -c    use caching (works per file path, honors mtime)

  -cache-path string
        where cache file will be stored (default "$XDG_CACHE_HOME/lsidups/" with fallback
                -> "$HOME/.cache/lsidups/" -> "$APPDATA/lsidups" -> current directory)

  -e value
        image extensions (with dots) to look for (default .jpg,.jpeg,.png,.gif)

  -i string
        directory to search (recursively) for duplicates, when set to - can take list of images
        to compare from stdin (default "-")

  -j    output duplicates as json instead of standard flat list

  -v    show time it took to complete key parts of the search
```

## Examples

find and list duplicates in ~/Pictures

```sh
lsidups -i ~/Pictures > dups.txt
```

<details>
  <summary>dups.txt</summary>

  ```sh
  /home/username/Pictures/image1.jpg
  /home/username/Pictures/dir/image1.jpg
  /home/username/Pictures/wdwd720p.jpg
  /home/username/Pictures/wdwd1080p.jpg
  /home/username/Pictures/wdwd1440p.jpg
  # ...
  ```
</details>

you could also export json

```sh
lsidups -j -i ~/Pictures > dups.json
```

<details>
  <summary>dups.txt</summary>

  ```json
  [
    [
      "/home/username/Pictures/image1.jpg",
      "/home/username/Pictures/dir/image1.jpg"
    ],
    [
      "/home/username/Pictures/wdwd720p.jpg",
      "/home/username/Pictures/wdwd1080p.jpg",
      "/home/username/Pictures/wdwd1440p.jpg"
    ]
  ]
  // ...
  ```
</details>

you can then sort images in groups e.g. by size

<details>
  <summary>sortsize.py</summary>

  ```python
  #!/usr/bin/env python3
  # sortsize.py
  # just pipe json into it

  import json
  import sys
  from os import path

  groups = json.load(sys.stdin)

  # sort files by size inside each group
  for i, group in enumerate(groups):
      groups[i] = sorted(group, key=lambda fp: path.getsize(fp), reverse=True)

  # if you want sorted json back
  # print(json.dumps(groups, indent=2))

  # if you want flat list
  for group in groups:
      for fp in group:
          print(fp)
  ```

  ```sh
  sortsize.py < dups.json > dups.txt
  ```
</details>

or compare just selected (e.g. with [fd](https://github.com/sharkdp/fd)) images

```sh
fd 'mashu' -e png --changed-within 2weeks ~/Pictures > yourlist.txt
lsidups < yourlist.txt > dups.txt
```

then process them in any image viewer that can read stdin ([sxiv](https://github.com/muennich/sxiv), [imv](https://github.com/eXeC64/imv))

```sh
sxiv -io < dups.txt
```
or

```sh
imv < dups.txt
```

Both of them allow you to map shell commands to keys, so the possibilities are endless. E.g. you could macgyver some [dmenu](https://tools.suckless.org/dmenu/)/[fzf](https://github.com/junegunn/fzf) based mover, use [trash-cli](https://github.com/andreafrancia/trash-cli) for deletion, etc.

### Caching

If you planning to run lsidups on the same directory multiple times - consider using cache to speed things up.

Note, that cache is stored in form of a hash table with pairs like _\*absoluteFilepath\*_: _\*imageProperties\*_, so you **don't** need to have different caches for different directories, because irrelevant images will be just filtered out, and new will be added to cache at the end of the run.

It is also smart enough to not use image from cache if it appears to has changed.

check for default cache file location on your system

```sh
lsidups -h
```

run with caching enabled

```sh
lsidups -c -i ~/Pictures > dups.txt
```

store cache file in the custom location (directories will be created for you if necessary)

```sh
lsidups -c -cache-path ~/where/to/store/cache.gob -i ~/Pictures > dups.txt
```
