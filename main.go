package main

// TODO tests?

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	flag.Parse()

	start := time.Now()

	var files []string
	var pics []Image

	if input == "-" {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			file := strings.TrimSpace(scanner.Text())
			if fpAbs, err := filepath.Abs(file); err == nil {
				files = append(files, fpAbs)
			}
		}
	} else {
		if link, err := filepath.EvalSymlinks(input); err == nil {
			files, _ = GetFiles(link)
		} else {
			fmt.Fprintf(os.Stderr, "> %s\n", err)
			os.Exit(1)
		}
	}

	files = FilterFiles(files, func(fp string) bool {
		// making sure it's image formats go supports & not system files
		ext := strings.ToLower(filepath.Ext(fp))
		return ContainsStr(searchExt, ext) && !IsHidden(fp)
	})

	if verbose {
		fmt.Fprintf(os.Stderr, "> found %d images, took %s\n", len(files), time.Since(start))
	}

	if len(files) <= 1 {
		os.Exit(0)
	}

	start = time.Now()

	usefulCache := 0
	if usecache {
		cachedPics, err := loadCache(cachepath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "> tried to load cache, but got - %s\n", err)
		} else {
			files, pics = filterCache(files, cachedPics)
			if verbose {
				usefulCache = len(pics)
				fmt.Fprintf(os.Stderr,
					"> loaded from cache %d images total, %d will be used, took %s\n",
					len(cachedPics), usefulCache, time.Since(start))
			}
		}
	}

	start = time.Now()

	// calculating image similarity hashes
	for pic := range MakeImages(files) {
		pics = append(pics, pic)
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "> processed %d images from disk, %d from cache, took %s\n",
			len(pics)-usefulCache, usefulCache, time.Since(start))
	}

	start = time.Now()

	if usecache && len(files) > 0 {
		err := storeCache(cachepath, pics)
		if err != nil {
			fmt.Fprintf(os.Stderr, "> tried to update cache, but got - %s\n", err)
		} else {
			if verbose {
				fmt.Fprintf(os.Stderr, "> updated cache, took %s\n", time.Since(start))
			}
		}
	}

	start = time.Now()

	// searching for similar images
	count := 0
	var groups [][]string
	for group := range findDups(pics) {
		groups = append(groups, group)
		count += len(group)
	}

	if exportjson {
		byt, _ := json.MarshalIndent(groups, "", "  ")
		fmt.Println(string(byt))
	} else {
		for _, group := range groups {
			for _, fp := range group {
				fmt.Println(fp)
				count++
			}
		}
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "> found %d similar images, took %s\n",
			count, time.Since(start))
	}
}
