package main

import (
	"sync"

	"github.com/corona10/goimagehash"
	"github.com/vitali-fedulov/images/v2"
)

func dupsSearcher(ipics <-chan Image, jpics *[]Image, dupInChan chan<- map[string][]string) {
	dups := make(map[string][]string)
	for ipic := range ipics {
		for _, jpic := range *jpics {
			if jpic.fp != ipic.fp {
				iPhash := goimagehash.NewImageHash(ipic.Phash, ipic.HashKind)
				jPhash := goimagehash.NewImageHash(jpic.Phash, jpic.HashKind)
				d, _ := iPhash.Distance(jPhash)
				if images.Similar(jpic.ImgHash, ipic.ImgHash, jpic.ImgSize, ipic.ImgSize) || d < gMaxDist {
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

	if !gNoMergeGroups {
		// merge groups if some of the items are the same
		for k, v := range dups {
			for _, fp := range v {
				for k2, v2 := range dups {
					if k2 != k {
						if containsStr(v2, fp) || containsStr(v2, k) {
							dups[k] = append(append(dups[k], v2...), k2)
							dups[k] = mkSetStr(dups[k])
							delete(dups, k2)
						}
					}
				}
			}
		}
	}

	for k, v := range dups {
		dupOut <- append(v, k)
	}
	close(dupOut)
}

// findDups takes slice of Images and concurrently searches for duplicates in
// them, and returns 2d slice of groups of duplicates. This should be stage 2
// of the search.
func findDups(pics []Image) <-chan []string {
	var wg sync.WaitGroup
	picsChan := make(chan Image)

	partialMaps := make(chan map[string][]string, gThreads)
	dupGroups := make(chan []string)

	go dupsMerger(partialMaps, dupGroups)

	for w := 1; w <= gThreads; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			dupsSearcher(picsChan, &pics, partialMaps)
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
		close(partialMaps)
	}()

	return dupGroups
}
