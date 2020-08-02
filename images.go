package main

import (
	"fmt"
	"image"
	"os"
	"sync"
	"time"

	"github.com/vitali-fedulov/images"
)

type Image struct {
	fp      string      // file path
	Mtime   time.Time   // file's last modification time
	ImgHash []float32   // similarity hash
	ImgSize image.Point // width x height
}

func makeImage(fp string) (Image, error) {
	pic, err := images.Open(fp)
	if err != nil {
		return Image{}, err
	}
	imgHash, imgSize := images.Hash(pic)
	// since it will return zero value anyway, error here does not actually matter
	mtime, _ := statMtime(fp)
	return Image{fp, mtime, imgHash, imgSize}, nil
}

func imageMaker(filesIn <-chan string, imagesOut chan<- Image, wg *sync.WaitGroup) {
	defer wg.Done()
	for fp := range filesIn {
		img, err := makeImage(fp)
		if err == nil {
			imagesOut <- img
		} else {
			fmt.Fprintf(os.Stderr, "> %s\n", err)
		}
	}
}

// MakeImages takes file pathes and concurrently makes Images for them. They
// can be used later to find duplicates with FindDups
func MakeImages(files []string) <-chan Image {
	var wg sync.WaitGroup

	filesIn := make(chan string)
	imagesOut := make(chan Image, len(files))

	for w := 1; w <= threads; w++ {
		go imageMaker(filesIn, imagesOut, &wg)
		wg.Add(1)
	}

	for _, fp := range files {
		filesIn <- fp
	}
	close(filesIn)

	wg.Wait()

	// yay, antipatterns! (actually it's ok when you sure)
	close(imagesOut)
	return imagesOut
}
