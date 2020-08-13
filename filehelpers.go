package main

import (
	"os"
	"path/filepath"
	"time"
)

// getFiles recursively walks dir tree and returns all files inside (absolute
// paths)
func getFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(fp string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			abs, err := filepath.Abs(fp)
			if err == nil {
				files = append(files, abs)
			}
		}
		return nil
	})
	return files, err
}

// filterFiles takes a slice of files and function to evaluate each file with
func filterFiles(slice []string, condition func(string) bool) []string {
	var newSlice []string
	for _, element := range slice {
		if condition(element) {
			newSlice = append(newSlice, element)
		}
	}
	return newSlice
}

// statMtime returns file mtime (when content of the file were last modified)
func statMtime(fp string) (time.Time, error) {
	file, err := os.Stat(fp)
	if err != nil {
		return time.Time{}, err
	}
	return file.ModTime(), nil
}
