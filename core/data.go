package core

import (
	"os"

	"github.com/kpango/glg"
)

//InitializeDataPath creates .gizo folder and block subfolder
func InitializeDataPath() {
	glg.Info("Core: Initializing Data Path")
	os.Mkdir(IndexPath, os.FileMode(0777))
	os.Mkdir(BlockPath, os.FileMode(0777))
}

//RemoveDataPath delete's .gizo folder
func RemoveDataPath() {
	glg.Info("Core: Removing Data Path")
	err := os.RemoveAll(IndexPath)
	if err != nil {
		glg.Fatal(err)
	}
}
