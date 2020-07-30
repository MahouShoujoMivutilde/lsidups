package main

import (
	"encoding/gob"
	"os"
	"path/filepath"
)

// LoadCache takes path to .gob file and returns cache map.
func LoadCache(cachepath string) (map[string]Image, error) {
	cachedPics := make(map[string]Image)
	file, err := os.Open(cachepath)
	if err != nil {
		return cachedPics, err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&cachedPics)
	if err != nil {
		return cachedPics, err
	}

	// since we did not export fp to save space - we need to set it back
	for fp, img := range cachedPics {
		img.fp = fp
		cachedPics[fp] = img
	}

	return cachedPics, nil
}

// StoreCache takes current slice of pictures and saves it to the given path in
// gob format, if cache already exists - appends new data to it.
func StoreCache(cachepath string, pics []Image) error {
	cachedPics, _ := LoadCache(cachepath)

	for _, img := range pics {
		cachedPics[img.fp] = img
	}

	err := os.MkdirAll(filepath.Dir(cachepath), 0700)
	if err != nil {
		return err
	}

	file, err := os.Create(cachepath)
	if err != nil {
		return err
	}

	// TODO maybe also gzip it?
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(&cachedPics)
	if err != nil {
		return err
	}
	return nil
}

// FilterCache returns files that did not have cache, and slice of Image for
// files that: 1) did have cache;  2) were not changed on disk (checked by
// comparing current and cached mtime).
func FilterCache(files []string, cachedPics map[string]Image) ([]string, []Image) {
	var filteredFiles []string
	var filteredPics []Image

	for _, fp := range files {
		mtimeNow, _ := statMtime(fp)
		img, ok := cachedPics[fp]

		if ok && img.Mtime.Equal(mtimeNow) && !img.Mtime.IsZero() && !mtimeNow.IsZero() {
			// cache should be valid
			filteredPics = append(filteredPics, img)
		} else {
			filteredFiles = append(filteredFiles, fp)
		}
	}

	return filteredFiles, filteredPics
}
