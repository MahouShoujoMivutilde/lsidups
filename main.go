package main

// TODO tests?
// TODO split to cmd and pkg

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

	if gTidyCache {
		purged, err := tidyCache(gCachePath)

		if err != nil {
			fmt.Fprintf(os.Stderr, "> tried to tidy cache, but got - %s\n", au.Red(err))
		} else {
			fmt.Fprintf(os.Stderr, "> purged %d from cache, took %s\n",
				au.Cyan(purged), au.Green(time.Since(start)))
		}
		os.Exit(0)
	}

	var files []string
	var pics []Image

	if gInput == "-" {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			file := strings.TrimSpace(scanner.Text())
			if fpAbs, err := filepath.Abs(file); err == nil {
				files = append(files, fpAbs)
			}
		}
	} else {
		if link, err := filepath.EvalSymlinks(gInput); err == nil {
			files, _ = getFiles(link)
		} else {
			fmt.Fprintf(os.Stderr, "> %s\n", au.Red(err))
			os.Exit(1)
		}
	}

	files = filterFiles(files, func(fp string) bool {
		// making sure it's image formats go supports & not system files
		ext := strings.ToLower(filepath.Ext(fp))
		return containsStr(gSearchExt, ext) && !isHidden(fp)
	})

	if gVerbose {
		fmt.Fprintf(os.Stderr, "> found %d images, took %s\n",
			au.Cyan(len(files)), au.Green(time.Since(start)))
	}

	if len(files) <= 1 {
		os.Exit(0)
	}

	start = time.Now()

	usefulCache := 0
	if gUseCache {
		cachedPics, err := loadCache(gCachePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "> tried to load cache, but got - %s\n", au.Red(err))
		} else {
			files, pics = filterCache(files, cachedPics)
			if gVerbose {
				usefulCache = len(pics)
				fmt.Fprintf(os.Stderr,
					"> loaded from cache %d images total, %d will be used, took %s\n",
					au.Cyan(len(cachedPics)), au.Cyan(usefulCache), au.Green(time.Since(start)))
			}
		}
	}

	start = time.Now()

	// calculating image similarity hashes
	for pic := range makeImages(files) {
		pics = append(pics, pic)
	}

	if gVerbose {
		fmt.Fprintf(os.Stderr, "> processed %d images from disk, %d from cache, took %s\n",
			au.Cyan(len(pics)-usefulCache), au.Cyan(usefulCache), au.Green(time.Since(start)))
	}

	start = time.Now()

	if gUseCache && len(pics)-usefulCache > 0 {
		cachedPics, _ := loadCache(gCachePath)

		for _, img := range pics {
			cachedPics[img.fp] = img
		}

		err := storeCache(gCachePath, cachedPics)
		if err != nil {
			fmt.Fprintf(os.Stderr, "> tried to update cache, but got - %s\n", au.Red(err))
		} else {
			if gVerbose {
				fmt.Fprintf(os.Stderr, "> updated cache, took %s\n", au.Green(time.Since(start)))
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

	if gExportJSON {
		byt, _ := json.MarshalIndent(groups, "", "  ")
		fmt.Println(string(byt))
	} else {
		for _, group := range groups {
			for _, fp := range group {
				fmt.Println(fp)
			}
		}
	}

	if gVerbose {
		fmt.Fprintf(os.Stderr, "> found %d similar images, took %s\n",
			au.Cyan(count), au.Green(time.Since(start)))
	}
}
