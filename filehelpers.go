package main

import (
	"os"
	"path/filepath"
	"strings"
)

// IsHidden checks if path fp has hidden elements, unix only
func IsHidden(fp string) bool {
	for _, element := range strings.Split(fp, string(filepath.Separator)) {
		if strings.HasPrefix(element, ".") && element != "." {
			return true
		}
	}
	return false
}

// GetFiles recursively walks dir tree and returns all files inside (absolute paths)
func GetFiles(dir string) ([]string, error) {
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

// FilterFiles takes a slice of files and function to evaluate each file with
func FilterFiles(slice []string, condition func(string) bool) []string {
	var newSlice []string
	for _, element := range slice {
		if condition(element) {
			newSlice = append(newSlice, element)
		}
	}
	return newSlice
}

// FilterExt takes a slice of files and slice of extensions (with dots)
// to search (searchExt), returns only files with extensions from searchExt
func FilterExt(files []string, searchExt []string) []string {
	return FilterFiles(files, func(fp string) bool {
		for _, ext := range searchExt {
			fpext := strings.ToLower(filepath.Ext(fp))
			if fpext == ext {
				return true
			}
		}
		return false
	})
}
