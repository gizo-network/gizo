package helpers

import "os"

func FileExists(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}
