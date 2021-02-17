package main

import (
	"encoding/gob"
	"os"
	"path/filepath"
)

// loadCache takes path to .gob file and returns cache map.
func loadCache(cachepath string) (map[string]Image, error) {
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

// storeCache takes cache map and saves it to the given path in gob format
func storeCache(cachepath string, cachedPics map[string]Image) error {
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

// filterCache returns files that did not have cache, and slice of Image for
// files that: 1) did have cache;  2) were not changed on disk (checked by
// comparing current and cached mtime).
func filterCache(files []string, cachedPics map[string]Image) ([]string, []Image) {
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

// tidyCache cleans cache from files that are missing on drive
func tidyCache(cachepath string) (int, error) {
	cachedPics, _ := loadCache(cachepath)
	initalCount := len(cachedPics)

	for fp, img := range cachedPics {
		stat, err := os.Stat(fp)
		if os.IsNotExist(err) || stat.ModTime().After(img.Mtime) {
			delete(cachedPics, fp)
		}
	}

	if purged := initalCount - len(cachedPics); purged > 0 {
		err := storeCache(cachepath, cachedPics)
		return purged, err
	}
	return 0, nil
}
