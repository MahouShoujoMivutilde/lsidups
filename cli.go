package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
	searchExt  extensions
	input      string
	verbose    bool
	usecache   bool
	exportjson bool
	cachepath  string
)

const DESCRIPTION string = `

  Is a tool for finding image duplicates (or just similar images).  Outputs
  images grouped by similarity (one filepath per line) to stdout so you can
  process them as you please.

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
`

func usage() {
	fmt.Fprint(flag.CommandLine.Output(), os.Args[0]+DESCRIPTION)

	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %[1]s:\n", os.Args[0])

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
	searchExt = extensions{".jpg", ".jpeg", ".png", ".gif"}
	cachepath = filepath.Join(cacheDir(), "cachemap.gob")

	flag.Var(&searchExt, "e", "image extensions (with dots) to look for")
	flag.StringVar(&input, "i", "-",
		"directory to search (recursively) for duplicates, when set to - can take list of images\n"+
			"to compare from stdin")
	flag.BoolVar(&verbose, "v", false, "show time it took to complete key parts of the search")
	flag.BoolVar(&exportjson, "j", false, "output duplicates as json instead of standard flat list")
	flag.BoolVar(&usecache, "c", false, "use caching (works per file path, honors mtime)")
	flag.StringVar(&cachepath, "cache-path", cachepath, "where cache file will be stored")

	flag.Usage = usage
}
