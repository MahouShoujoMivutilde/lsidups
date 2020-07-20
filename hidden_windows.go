package main

import (
	"syscall"
)

func IsHidden(fp string) bool {
	ptr, err := syscall.UTF16PtrFromString(fp)
	if err != nil {
		return false
	}
	attrs, err := syscall.GetFileAttributes(ptr)
	if err != nil {
		return false
	}
	return attrs&syscall.FILE_ATTRIBUTE_HIDDEN == syscall.FILE_ATTRIBUTE_HIDDEN
}
