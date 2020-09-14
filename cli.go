package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/logrusorgru/aurora"
	"github.com/mattn/go-isatty"
)

type extensions []string

func (e *extensions) String() string {
	return strings.Join(*e, ",")
}

func (e *extensions) Set(val string) error {
	*e = extensions{}
	for _, ext := range strings.Split(val, ",") {
		*e = append(*e, ext)
	}
	return nil
}

var (
	gSearchExt     extensions
	gInput         string
	gVerbose       bool
	gUseCache      bool
	gTidyCache     bool
	gExportJSON    bool
	gNoMergeGroups bool
	gCachePath     string
	gThreads       int
	gMaxDist       int
	au             aurora.Aurora
)

const DESCRIPTION string = `

  Is a tool for finding image duplicates (or just similar images).  Outputs
  images grouped by similarity to stdout so you can process them as you please.

  It uses https://github.com/vitali-fedulov/images (for its precise matching)
  together with phash from https://github.com/corona10/goimagehash (to detect
  cropped images and allow for variable similarity threshold)

`

const EXAMPLES string = `
Examples:
  find duplicates in ~/Pictures
    lsidups -i ~/Pictures > dups.txt

  or compare just selected images
    fd 'mashu' -e png --changed-within 2weeks ~/Pictures > yourlist.txt
    lsidups < yourlist.txt > dups.txt

  then process them in any image viewer that can read stdin (sxiv, imv...)
    sxiv -io < dups.txt

  or you could export json instead
    lsidups -j < yourlist.txt > dups.json

  if you planning to run lsidups on the same directory multiple times
  - consider using cache to speed things up
    lsidups -c -i ~/Pictures > dups.txt

  if you want to save cache file to the custom location (directories will be created
  for you if necessary)
    lsidups -c -cache-path ~/where/to/store/cache.gob -i ~/Pictures > dups.txt

  also it is worth noting that lsidups merges groups if some of their items are the same.

  i think it makes sense from the user perspective, but the resulting group
  might contain images that are not all actually similar with each other: let's
  say we have 3 images: [1.png 2.png 3.png], 1 and 2 hashes are similar enough
  to be consider related images, and 2 and 3 also similar enough, but 1 and 3
  are far apart enough to be consider different.  If you want to get 2 groups:
  [1.png 2.png] and [2.png 3.png] - pass flag -g
`

func usage() {
	name := filepath.Base(os.Args[0])
	fmt.Fprint(flag.CommandLine.Output(), name+DESCRIPTION)

	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %[1]s:\n", name)

	flag.PrintDefaults()
	fmt.Fprint(flag.CommandLine.Output(), EXAMPLES)
}

func cacheDir() string {
	base := "."
	if d := os.Getenv("XDG_CACHE_HOME"); d != "" {
		base = d
		base = filepath.Join(d, "lsidups")
	} else if d := os.Getenv("HOME"); d != "" {
		base = filepath.Join(d, ".cache", "lsidups")
	} else if d := os.Getenv("APPDATA"); d != "" {
		base = filepath.Join(d, "lsidups")
	} else {
		base, _ = os.Getwd()
	}
	return base
}

func init() {
	gSearchExt = extensions{".jpg", ".jpeg", ".png", ".gif"}
	gCachePath = filepath.Join(cacheDir(), "cachemap.gob")

	flag.Var(&gSearchExt, "e", "image extensions (with dots) to look for")
	flag.StringVar(&gInput, "i", "-",
		"directory to search (recursively) for duplicates, when set to - can take list of images\n"+
			"to compare from stdin")
	flag.BoolVar(&gVerbose, "v", false, "show time it took to complete key parts of the search")
	flag.BoolVar(&gExportJSON, "j", false, "output duplicates as json instead of standard flat list")
	flag.BoolVar(&gUseCache, "c", false, "use caching (works per file path, honors mtime)")
	flag.BoolVar(&gTidyCache, "ct", false, "remove missing (on drive) files from cache")
	flag.BoolVar(&gNoMergeGroups, "g", false, "do not merge groups if some of the items are the same")
	flag.StringVar(&gCachePath, "cache-path", gCachePath, "where cache file will be stored")
	flag.IntVar(&gThreads, "T", runtime.NumCPU(), "number of processing threads")
	flag.IntVar(&gMaxDist, "d", 8, "phash threshold distance (less = more precise match, but more false negatives)")

	if gThreads < 1 {
		gThreads = runtime.NumCPU()
	}

	// color output only when stderr is terminal
	au = aurora.NewAurora(isatty.IsTerminal(os.Stderr.Fd()))

	flag.Usage = usage
}
