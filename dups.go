package main

import (
	"sync"

	"github.com/vitali-fedulov/images"
)

func dupsSearcher(ipics <-chan Image, jpics *[]Image, dupInChan chan<- map[string][]string) {
	dups := make(map[string][]string)
	for ipic := range ipics {
		for _, jpic := range *jpics {
			if jpic.fp != ipic.fp {
				if images.Similar(jpic.ImgHash, ipic.ImgHash, jpic.ImgSize, ipic.ImgSize) {
					dups[ipic.fp] = append(dups[ipic.fp], jpic.fp)
				}
			}
		}
	}
	dupInChan <- dups
}

func dupsMerger(dupIn <-chan map[string][]string, dupOut chan<- []string) {
	dups := make(map[string][]string)

	for partialMap := range dupIn {
		for k, v := range partialMap {
			// NOTE keys are unique, see dupsSearcher
			dups[k] = v
		}
	}

	// get rid of mirror pairs
	for _, v := range dups {
		for _, fp := range v {
			delete(dups, fp)
		}
	}

	for k, v := range dups {
		dupOut <- append(v, k)
	}
	close(dupOut)
}

// FindDups takes slice of Images and concurrently searches for duplicates in
// them, and returns 2d slice of groups of duplicates. This should be stage 2
// of the search.
func FindDups(pics []Image) <-chan []string {
	var wg sync.WaitGroup
	picsChan := make(chan Image, len(pics))

	pairChan := make(chan map[string][]string, len(pics))
	dupGroupsChan := make(chan []string, len(pics))

	go dupsMerger(pairChan, dupGroupsChan)

	for w := 1; w <= threads; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			dupsSearcher(picsChan, &pics, pairChan)
		}()
	}

	for _, pic := range pics {
		picsChan <- pic
	}
	close(picsChan)

	go func() {
		wg.Wait()
		// compared everything, no more pairs will be sent, dupsHolder should
		// finish what left
		close(pairChan)
	}()

	return dupGroupsChan
}
