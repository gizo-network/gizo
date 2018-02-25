package core

import (
	"os"

	"github.com/kpango/glg"
)

//InitializeDataPath creates .gizo folder and block subfolder
func InitializeDataPath(path string) {
	glg.Info("Core: Initializing Data Path")
	os.Mkdir(path, os.FileMode(0777))
	os.Mkdir(path, os.FileMode(0777))
}

//RemoveDataPath delete's .gizo folder
func RemoveDataPath(path string) {
	glg.Info("Core: Removing Data Path")
	err := os.RemoveAll(path)
	if err != nil {
		glg.Fatal(err)
	}
}
