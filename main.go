package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/vitali-fedulov/images"
)

type Image struct {
	Fp      string
	ImgHash []float32
	ImgSize image.Point
}

func makeImage(fp string) (Image, error) {
	pic, err := images.Open(fp)
	if err != nil {
		return Image{}, err
	}
	imgHash, imgSize := images.Hash(pic)
	return Image{fp, imgHash, imgSize}, nil
}

func imageMaker(jobs <-chan string, results chan<- Image, wg *sync.WaitGroup) {
	defer wg.Done()
	for fp := range jobs {
		img, err := makeImage(fp)
		if err == nil {
			results <- img
		} else {
			fmt.Fprintf(os.Stderr, "> %s - %s\n", fp, err)
		}
	}
}

func insertAfterFp(arr []string, fp string, newFp string) []string {
	after := -1
	for i, e := range arr {
		if e == fp {
			after = i
			break
		}
	}

	if after == -1 {
		return arr
	}

	// increase capacity for new element to fit
	arr = append(arr, "")

	// shift by 1 all elements after "after"
	copy(arr[after+1:], arr[after:])
	arr[after+1] = newFp
	return arr
}

func dupsSearch(pics <-chan Image, ipics *[]Image, dupInChan chan<- []string, wg *sync.WaitGroup) {
	defer wg.Done()
	for pic := range pics {
		for _, ipic := range *ipics {
			if ipic.Fp != pic.Fp {
				if images.Similar(ipic.ImgHash, pic.ImgHash, ipic.ImgSize, pic.ImgSize) {
					dupInChan <- []string{ipic.Fp, pic.Fp}
				}
			}
		}
	}
}

func dupsHolder(dupInChan <-chan []string, dupOutChan chan<- string, doneChan <-chan bool) {
	var duplicates []string
	for {
		select {
		case pair := <-dupInChan:
			ipicFp, picFp := pair[0], pair[1]
			ipicIn := ContainsStr(duplicates, ipicFp)
			picIn := ContainsStr(duplicates, picFp)

			if picIn && !ipicIn {
				duplicates = insertAfterFp(duplicates, picFp, ipicFp)
			} else if !picIn && ipicIn {
				duplicates = insertAfterFp(duplicates, ipicFp, picFp)
			} else if !picIn && !ipicIn {
				duplicates = append(duplicates, picFp, ipicFp)
			}
		case <-doneChan:
			for _, fp := range duplicates {
				dupOutChan <- fp
			}
			close(dupOutChan)
			return
		}
	}
}

func loadCache(cachepath string) (map[string]Image, error) {
	cachedPics := make(map[string]Image)
	byt, err := ioutil.ReadFile(cachepath)
	if err != nil {
		return cachedPics, err
	}

	// TODO unmarshal is slow (takes about 1s) when we have alot (like 5k+) Images
	err = json.Unmarshal(byt, &cachedPics)
	if err != nil {
		return cachedPics, err
	}

	return cachedPics, nil
}

// filterCache returns files that did not have cache, and
// slice of Image for pictures that did
func filterCache(files []string, cachedPics map[string]Image) ([]string, []Image) {
	var filteredFiles []string
	var filteredPics []Image

	for _, fp := range files {
		img, ok := cachedPics[fp]
		if ok {
			filteredPics = append(filteredPics, img)
		} else {
			filteredFiles = append(filteredFiles, fp)
		}
	}

	return filteredFiles, filteredPics
}

func storeCache(cachepath string, pics []Image) error {
	cachedPics, _ := loadCache(cachepath)

	for _, img := range pics {
		cachedPics[img.Fp] = img
	}

	serializedPics, err := json.Marshal(cachedPics)

	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(cachepath), 0700)
	if err != nil {
		return err
	}

	// TODO maybe also gzip it?
	err = ioutil.WriteFile(cachepath, serializedPics, 0600)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	flag.Parse()

	start := time.Now()

	var files []string
	// TODO make this a map
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

	if usecache {
		cachedPics, err := loadCache(cachepath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "> tried to load cache, but got - %s\n", err)
		} else {
			files, pics = filterCache(files, cachedPics)
			if verbose {
				fmt.Fprintf(os.Stderr,
					"> loaded from cache %d images total, %d will be used, took %s\n",
					len(cachedPics), len(pics), time.Since(start))
			}
		}
	}

	start = time.Now()
	// calculating image similarity hashes
	jobs := make(chan string)
	results := make(chan Image, len(files))

	var wg sync.WaitGroup
	for w := 1; w <= runtime.NumCPU(); w++ {
		wg.Add(1)
		go imageMaker(jobs, results, &wg)
	}

	for _, fp := range files {
		jobs <- fp
	}
	close(jobs)

	wg.Wait()
	// yay, antipatterns! (actually it's ok when you sure)
	close(results)

	for pic := range results {
		pics = append(pics, pic)
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "> processed %d images from disk, %d from cache, took %s\n",
			len(files), len(pics)-len(files), time.Since(start))
	}

	start = time.Now()

	if usecache {
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
	picsChan := make(chan Image, len(pics))

	dupInChan := make(chan []string, len(pics))
	dupOutChan := make(chan string, len(pics))
	doneChan := make(chan bool)

	go dupsHolder(dupInChan, dupOutChan, doneChan)

	for w := 1; w <= runtime.NumCPU(); w++ {
		wg.Add(1)
		go dupsSearch(picsChan, &pics, dupInChan, &wg)
	}

	for _, pic := range pics {
		picsChan <- pic
	}
	close(picsChan)

	wg.Wait()
	doneChan <- true

	count := 0
	for fp := range dupOutChan {
		fmt.Println(fp)
		count++
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "> found %d similar images, took %s\n",
			count, time.Since(start))
	}
}
