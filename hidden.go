// +build !windows

package main

import (
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
