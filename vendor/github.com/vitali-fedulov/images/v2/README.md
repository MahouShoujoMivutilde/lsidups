# Comparing images in Go &#10132; [LATEST version](https://github.com/vitali-fedulov/images4)

Near duplicates and resized images can be found with the module.

**Demo**: [similar image search and clustering](https://vitali-fedulov.github.io/similar.pictures/) (deployed [from](https://github.com/vitali-fedulov/similar.pictures)).

**Semantic versions**:
- v1/v2 (/images) - this repository,
- [v3](https://github.com/vitali-fedulov/images3) (/images3),
- [v4](https://github.com/vitali-fedulov/images4) (/images4) - latest recommended.

All versions will be kept available indefinitely.

# About this repo

There are no dependencies: only the Golang standard library is used. Supported image types: GIF, JPEG and PNG (golang.org/pkg/image/ as in October 2018).

`Similar` function gives a verdict whether 2 images are similar or not. The library also contains wrapper functions to open/save images and basic image resampling/resizing.

`SimilarCustom` function allows your own similarity metric thresholds.

Documentation: [godoc](https://pkg.go.dev/github.com/vitali-fedulov/images/v2).

## Example of comparing 2 photos

To test this example go-file, you need to initialize modules from command line, because the latest version (v2) uses them:

`go mod init foo`

Here `foo` can be anything for testing purposes. Then get the required import:

`go get github.com/vitali-fedulov/images/v2`

Now you are ready to run or build the example.

```go
package main

import (
	"fmt"

	// Notice v2, which is module-based most recent version.
	// Explanation: https://go.dev/blog/v2-go-modules
	"github.com/vitali-fedulov/images/v2"
)

func main() {
	
	// Open photos.
	imgA, err := images.Open("photoA.jpg")
	if err != nil {
		panic(err)
	}
	imgB, err := images.Open("photoB.jpg")
	if err != nil {
		panic(err)
	}
	
	// Calculate hashes and image sizes.
	hashA, imgSizeA := images.Hash(imgA)
	hashB, imgSizeB := images.Hash(imgB)
	
	// Image comparison.
	if images.Similar(hashA, hashB, imgSizeA, imgSizeB) {
		fmt.Println("Images are similar.")
	} else {
		fmt.Println("Images are distinct.")
	}
}
```

## Algorithm for image comparison

[Detailed explanation](https://vitali-fedulov.github.io/similar.pictures/algorithm-for-perceptual-image-comparison.html), also as a [PDF](https://github.com/vitali-fedulov/research/blob/main/Algorithm%20for%20perceptual%20image%20comparison%20OLD.pdf).

Summary: In the algorithm images are resized to small squares of fixed size.
A number of masks representing several sample pixels are run against the resized
images to calculate average color values. Then the values are compared to
give the similarity verdict. Also image proportions are used to avoid matching
images of distinct shape.
