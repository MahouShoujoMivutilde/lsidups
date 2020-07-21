package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/vitali-fedulov/images"
)

type Image struct {
	fp      string
	imgHash []float32
	imgSize image.Point
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
			if ipic.fp != pic.fp {
				if images.Similar(ipic.imgHash, pic.imgHash, ipic.imgSize, pic.imgSize) {
					dupInChan <- []string{ipic.fp, pic.fp}
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

func main() {
	flag.Parse()

	start := time.Now()

	var files []string

	if input == "-" {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			file := strings.TrimSpace(scanner.Text())
			if fpAbs, err := filepath.Abs(file); err == nil {
				files = append(files, fpAbs)
			}
		}
	} else {
		files, _ = GetFiles(input)
	}

	files = FilterFiles(files, func(fp string) bool {
		// making sure it's image formats go supports & not system files
		ext := strings.ToLower(filepath.Ext(fp))
		return ContainsStr(searchExt, ext) && !IsHidden(fp)
	})

	if verbose {
		fmt.Fprintf(os.Stderr, "> found %d images, took %s\n", len(files), time.Since(start))
	}

	start = time.Now()

	// calculating image similarity hashes
	numJobs := len(files)
	jobs := make(chan string, numJobs)
	results := make(chan Image, numJobs)

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

	if verbose {
		fmt.Fprintf(os.Stderr, "> processed images, took %s\n", time.Since(start))
	}

	var pics []Image
	for pic := range results {
		pics = append(pics, pic)
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
