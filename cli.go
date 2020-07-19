package main

import (
	"flag"
	"fmt"
	"os"
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

var searchExt extensions
var input string
var verbose bool

var DESC string = os.Args[0] + `

  Is a tool for finding image dupicates (or just similar images).
  Outputs images grouped by similarity (one filepath per line) to stdio
  so you can process them as you please.

`

var EXAMPLES string = `
Examples:
  find duplicates in ~/Pictures
    lsidups -i ~/Pictures > dups.txt

  or compare just selected images
    fd 'mashu' -e png --changed-within 2weeks ~/Pictures > yourlist.txt
    lsidups -i - < yourlist.txt > dups.txt

  then process them in any image viewer that can read stdio (sxiv, imv...)
    sxiv -io < dups.txt
`

func usage() {
	fmt.Fprint(flag.CommandLine.Output(), DESC)

	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %[1]s:\n", os.Args[0])

	flag.PrintDefaults()
	fmt.Fprint(flag.CommandLine.Output(), EXAMPLES)
}

func init() {
	searchExt = extensions{".jpg", ".jpeg", ".png", ".gif"}
	flag.Var(&searchExt, "e", "image extensions (with dots) to look for")
	flag.StringVar(&input, "i", "-",
		"directory to search (recursively) for duplicates, when set to - can take list of images\n"+
			"to compare from stdio (one filepath per line, like from find & fd...)")
	flag.BoolVar(&verbose, "v", false,
		"show time it took to complete key parts of the search")
	flag.Usage = usage
}

