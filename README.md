
# lsidups

### ...is a barebone tool for finding image duplicates (or just similar images) from your terminal.

### How to use
Pipe a list of files to compare into stdio or just input (`-i`) a directory you want to check. It then will output images grouped by similarity so you can process them as you please.

It uses library [images](https://github.com/vitali-fedulov/images) (MIT) under the hood, so if you want to know more about  limitations  and how comparison works - read [this](https://similar.pictures/algorithm-for-perceptual-image-comparison.html).

lsidups itself is just a wrapper that tries to provide a way to compare a lot (10k+) images reasonably fast from cli, and then allow you to process found duplicates in some other more convenient tool, like e.g. [sxiv](https://github.com/muennich/sxiv).

### Image formats support

At the moment of writing, it supports **only** **jpeg**, **png** and **gif**; i tried to make webp work, but [this](https://github.com/golang/go/issues/38341) prevented it.  _\*very sad webp UwU\*._

### Install

Make sure you have go and git installed.

```
go get github.com/MahouShoujoMivutilde/lsidups
```

## Options

```
Usage of lsidups:
  -c    cache similarity hashes per image path
  -cache-path string
        where cache file will be stored (default "$XDG_CACHE_HOME/lsidups/" with fallback
                -> "$HOME/.cache/lsidups/" -> $APPDATA/lsidups -> current directory)
  -e value
        image extensions (with dots) to look for (default .jpg,.jpeg,.png,.gif)
  -i string
        directory to search (recursively) for duplicates, when set to - can take list of images
        to compare from stdin (default "-")
  -v    show time it took to complete key parts of the search
```

## Examples

find duplicates in ~/Pictures

```
lsidups -i ~/Pictures > dups.txt
```

or compare just selected images
```
fd 'mashu' -e png --changed-within 2weeks ~/Pictures > yourlist.txt
lsidups -i - < yourlist.txt > dups.txt
```

then process them in any image viewer that can read stdio ([sxiv](https://github.com/muennich/sxiv), [imv](https://github.com/eXeC64/imv))

```
sxiv -io < dups.txt
```
or

```
imv < dups.txt
```

Both of them allow you to map shell commands to keys, so the possibilities are endless. E.g. you could macgyver some [dmenu](https://tools.suckless.org/dmenu/) based mover, use [trash-cli](https://github.com/andreafrancia/trash-cli) for deletion.
