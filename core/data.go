package core

import (
	"os"

	"github.com/kpango/glg"
)

//InitializeDataPath creates .gizo folder and block subfolder
func InitializeDataPath() {
	if os.Getenv("ENV") == "dev" {
		glg.Info("Core: Initializing Dev Data Path")
		os.Mkdir(BlockPathDev, os.FileMode(0777))
		os.Mkdir(IndexPathDev, os.FileMode(0777))
	} else {
		glg.Info("Core: Initializing Prod Data Path")
		os.Mkdir(BlockPathProd, os.FileMode(0777))
		os.Mkdir(IndexPathProd, os.FileMode(0777))
	}
}

//RemoveDataPath delete's .gizo folder
func RemoveDataPath() {
	if os.Getenv("ENV") == "dev" {
		glg.Info("Core: Removing Dev Data Path")
		err := os.RemoveAll(IndexPathDev)
		if err != nil {
			glg.Fatal(err)
		}
	} else {
		glg.Info("Core: Removing Prod Data Path")
		err := os.RemoveAll(IndexPathProd)
		if err != nil {
			glg.Fatal(err)
		}
	}
}
