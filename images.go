package main

import (
	"fmt"
	"image"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"github.com/corona10/goimagehash"
	"github.com/vitali-fedulov/images/v2"

	_ "golang.org/x/image/webp"
)

type Image struct {
	fp       string      // file path
	Mtime    time.Time   // file's last modification time
	ImgHash  []float32   // similarity hash from fedulov's images
	ImgSize  image.Point // width x height
	Phash    uint64
	HashKind goimagehash.Kind
}

func makeImage(fp string) (img Image, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("likely upstream bug: %s\n%s",
				au.Red(r), au.Index(210, debug.Stack()))
		}
	}()
	pic, err := images.Open(fp)
	if err != nil {
		return Image{}, err
	}
	imgHash, imgSize := images.Hash(pic)
	pHash, err := goimagehash.PerceptionHash(pic)
	if err != nil {
		return Image{}, err
	}
	hash := pHash.GetHash()
	kind := pHash.GetKind()
	// since it will return zero value anyway, error here does not actually matter
	mtime, _ := statMtime(fp)
	return Image{fp, mtime, imgHash, imgSize, hash, kind}, nil
}

func imageMaker(filesIn <-chan string, imagesOut chan<- Image) {
	for fp := range filesIn {
		img, err := makeImage(fp)
		if err == nil {
			imagesOut <- img
		} else {
			fmt.Fprintf(os.Stderr, "> %s - %s\n", au.Index(117, fp), au.Red(err))
		}
	}
}

// makeImages takes file pathes and concurrently makes Images for them. They
// can be used later to find duplicates with FindDups
func makeImages(files []string) <-chan Image {
	var wg sync.WaitGroup

	filesIn := make(chan string)
	imagesOut := make(chan Image, len(files))

	for w := 1; w <= gThreads; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			imageMaker(filesIn, imagesOut)
		}()
	}

	for _, fp := range files {
		filesIn <- fp
	}
	close(filesIn)

	go func() {
		wg.Wait()
		// processed everything, no new images will be sent
		close(imagesOut)
	}()

	return imagesOut
}
