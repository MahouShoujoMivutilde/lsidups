# lsidups

### ...is a barebone tool for finding image duplicates (or just similar images) from your terminal.

### How to use
Pipe a list of files to compare into stdin or just input (`-i`) a directory you want to check. It then will output images grouped by similarity so you can process them as you please.

It mainly relies on [images](https://github.com/vitali-fedulov/images) (MIT).

Phash from [goimagehash](github.com/corona10/goimagehash) (BSD-2) is used to catch cropped duplicates (with that `images` tends to struggle) and to allow for variable similarity threshold.

lsidups itself is just a wrapper that tries to provide a way to compare a lot (10k+) images reasonably fast from cli, and then allow you to process found duplicates in some other more convenient tool, like e.g. [nsxiv](https://codeberg.org/nsxiv/nsxiv/) or with some custom script (see [examples directory](examples)).

### Image formats support

At the moment of writing, it supports **only** **jpeg**, **png**, **gif** and **webp** (_but [not some of the ones with ICC profile](https://github.com/golang/go/issues/60437#issuecomment-1563939784)_).

### Video

If you want to find video duplicates instead - try [lsvdups](examples/lsvdups) (it's not very good, though).

### Install

_NOTE:_ If upgrading - consider deleting cache:

```sh
rm $XDG_CACHE_HOME/lsidups/*
```

#### Arch way

[in AUR](https://aur.archlinux.org/packages/lsidups-git/):

```
lsidups-git
```

#### Go way

Make sure you have go and git installed, and `$(go env GOPATH)/bin` is in your `$PATH`.

```sh
go install github.com/MahouShoujoMivutilde/lsidups@latest
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
        remove missing/changed (on drive) files from cache and exit

  -d int
        phash threshold distance (less = more precise match, but more false negatives) (default 8)

  -e value
        image extensions (with dots) to look for (default .jpg,.jpeg,.png,.gif,.webp)

  -g    do not merge groups if some of the items are the same (default will merge)

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

then process them in any image viewer that can read stdin ([nsxiv](https://codeberg.org/nsxiv/nsxiv/), [imv](https://github.com/eXeC64/imv))

```sh
nsxiv -io < dups.txt
```
or

```sh
imv < dups.txt
```

Both of them allow you to map shell commands to keys, so the possibilities are endless. E.g. you could macgyver some [dmenu](https://tools.suckless.org/dmenu/)/[fzf](https://github.com/junegunn/fzf) based mover, use [trash-cli](https://github.com/andreafrancia/trash-cli) for deletion, etc.

Or a more complex example - find images present in `folderA`, but not in `folderB`:

```sh
comm -23 \
    <(fd -t f -e png -e jpeg -e jpg -e webp . ~/pics/folderA | sort) \
    <(fd -t f -e png -e jpeg -e jpg -e webp . ~/pics/folderA ~/folderB | lsidups -c | sort)

```

Also it is worth noting that lsidups merges groups if some of their items are the same. I think it makes sense from the user perspective, but the resulting group might contain images that are not all actually similar with each other.

Let's say we have 3 images: 1.png, 2.png, 3.png.

Hashes of 1 and 2 are similar enough to be considered related, and 2 and 3 are also similar enough, but 1 and 3 are far apart enough to be considered different.

By default they will be grouped like [1.png 2.png 3.png].

If you want to get 2 groups: [1.png 2.png] and [2.png 3.png] - pass flag `-g`.

### Caching

If you're planning to run lsidups on the same directory multiple times - consider using cache to speed things up.

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

Cache from older versions _might_ become invalid after upgrades.
