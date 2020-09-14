# lsidups

### ...is a barebone tool for finding image duplicates (or just similar images) from your terminal.

### How to use
Pipe a list of files to compare into stdin or just input (`-i`) a directory you want to check. It then will output images grouped by similarity so you can process them as you please.

It mainly relies on [images](https://github.com/vitali-fedulov/images) (MIT) library under the hood, so if you want to know more about limitations and how comparison works - read [this](https://similar.pictures/algorithm-for-perceptual-image-comparison.html).

Phash from [goimagehash](github.com/corona10/goimagehash) (BSD-2) is used to catch cropped duplicates (with that `images` tends to struggle) and to allow for variable similarity threshold.

lsidups itself is just a wrapper that tries to provide a way to compare a lot (10k+) images reasonably fast from cli, and then allow you to process found duplicates in some other more convenient tool, like e.g. [sxiv](https://github.com/muennich/sxiv) or with some custom script (see [examples directory](examples)).

### Image formats support

At the moment of writing, it supports **only** **jpeg**, **png** and **gif**; i tried to make webp work, but [this](https://github.com/golang/go/issues/38341) prevented it. _\*very sad webp UwU\*._

### Video

If you want to find video duplicates instead - try [lsvdups](examples/lsvdups).

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

  -ct
        remove missing (on drive) files from cache

  -d int
        phash threshold distance (less = more precise match, but more false negatives) (default 8)

  -e value
        image extensions (with dots) to look for (default .jpg,.jpeg,.png,.gif)

  -g    do not merge groups if some of the items are the same

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

you can then sort images in groups e.g. by file size with [sortsize.py](examples/sortsize.py), see [examples directory](examples) for more

```sh
sortsize.py < dups.json > dups.txt
```

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

Also it is worth noting that lsidups merges groups if some of their items are the same. I think it makes sense from the user perspective, but the resulting group might contain images that are not all actually similar with each other.

Let's say we have 3 images: 1.png, 2.png, 3.png.

Hashes of 1 and 2 are similar enough to be considered related, and 2 and 3 are also similar enough, but 1 and 3 are far apart enough to be considered different.

By default they will be grouped like [1.png 2.png 3.png].

If you want to get 2 groups: [1.png 2.png] and [2.png 3.png] - pass flag `-g`.

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
