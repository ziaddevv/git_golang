package utils

import "os"

func Exists(path string) bool {
	info, err := os.Stat(path)
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}

	return info != nil
}
