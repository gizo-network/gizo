package helpers

import "os"

//FileExists checks if file exists on the disk
func FileExists(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}
