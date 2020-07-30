package main

import (
	"runtime"
	"sync"

	"github.com/vitali-fedulov/images"
)

func dupsSearcher(ipics <-chan Image, jpics *[]Image, dupInChan chan<- []string, wg *sync.WaitGroup) {
	defer wg.Done()
	for ipic := range ipics {
		for _, jpic := range *jpics {
			if jpic.fp != ipic.fp {
				if images.Similar(jpic.ImgHash, ipic.ImgHash, jpic.ImgSize, ipic.ImgSize) {
					dupInChan <- []string{jpic.fp, ipic.fp}
				}
			}
		}
	}
}

func dupsHolder(dupIn <-chan []string, dupOut chan<- []string, done <-chan bool) {
	var dups [][]string
	for {
		select {
		case pair := <-dupIn:
			jpicFp, ipicFp := pair[0], pair[1]

			ipicGroup := findGroup(dups, ipicFp)
			jpicGroup := findGroup(dups, jpicFp)

			if jpicGroup == -1 && ipicGroup == -1 {
				dups = append(dups, []string{ipicFp, jpicFp})
			} else if jpicGroup != -1 && ipicGroup == -1 {
				dups[jpicGroup] = append(dups[jpicGroup], ipicFp)
			} else if jpicGroup == -1 && ipicGroup != -1 {
				dups[ipicGroup] = append(dups[ipicGroup], jpicFp)
			}
		case <-done:
			for _, group := range dups {
				dupOut <- group
			}
			close(dupOut)
			return
		}
	}
}

// findDups takes slice of Images and concurrently searches for duplicates in
// them, and returns 2d slice of groups of duplicates This should be stage 2 of
// the search.
func findDups(pics []Image) <-chan []string {
	var wg sync.WaitGroup
	picsChan := make(chan Image, len(pics))

	pairChan := make(chan []string, len(pics))
	dupGroupsChan := make(chan []string, len(pics))
	doneChan := make(chan bool)

	go dupsHolder(pairChan, dupGroupsChan, doneChan)

	for w := 1; w <= runtime.NumCPU(); w++ {
		wg.Add(1)
		go dupsSearcher(picsChan, &pics, pairChan, &wg)
	}

	for _, pic := range pics {
		picsChan <- pic
	}
	close(picsChan)

	wg.Wait()
	doneChan <- true

	return dupGroupsChan
}
