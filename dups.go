package main

import (
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

func removeGroup(s [][]string, index int) [][]string {
	ret := make([][]string, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
}

func dupsHolder(dupIn <-chan []string, dupOut chan<- []string) {
	var dups [][]string
	for pair := range dupIn {
		jpicFp, ipicFp := pair[0], pair[1]

		ipicGroup := findGroup(dups, ipicFp)
		jpicGroup := findGroup(dups, jpicFp)

		if jpicGroup == -1 && ipicGroup == -1 {
			dups = append(dups, []string{ipicFp, jpicFp})
		} else if jpicGroup != -1 && ipicGroup == -1 {
			dups[jpicGroup] = append(dups[jpicGroup], ipicFp)
		} else if jpicGroup == -1 && ipicGroup != -1 {
			dups[ipicGroup] = append(dups[ipicGroup], jpicFp)
		} else if jpicGroup != -1 && ipicGroup != -1 && jpicGroup != ipicGroup {
			// both found, but in different groups, so we merge 2 groups
			dups[ipicGroup] = append(dups[ipicGroup], dups[jpicGroup]...)
			dups = removeGroup(dups, jpicGroup)
		}
	}

	for _, group := range dups {
		dupOut <- group
	}
	close(dupOut)
}

// FindDups takes slice of Images and concurrently searches for duplicates in
// them, and returns 2d slice of groups of duplicates. This should be stage 2
// of the search.
func FindDups(pics []Image) <-chan []string {
	var wg sync.WaitGroup
	picsChan := make(chan Image, len(pics))

	pairChan := make(chan []string, len(pics))
	dupGroupsChan := make(chan []string, len(pics))

	go dupsHolder(pairChan, dupGroupsChan)

	for w := 1; w <= threads; w++ {
		wg.Add(1)
		go dupsSearcher(picsChan, &pics, pairChan, &wg)
	}

	for _, pic := range pics {
		picsChan <- pic
	}
	close(picsChan)

	wg.Wait()
	// compared everything, no more pairs will be sent, dupsHolder should
	// finish what left
	close(pairChan)

	return dupGroupsChan
}
