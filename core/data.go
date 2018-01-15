package core

import (
	"os"

	"github.com/kpango/glg"
)

func InitializeDataPath() {
	err := os.Mkdir(IndexPath, os.FileMode(0777))
	if err != nil {
		glg.Info("Using existing data path")
	} else {
		glg.Info("Initializing data path")
	}
	err = os.Mkdir(BlockPath, os.FileMode(0777))
	if err != nil {
		glg.Info("Using existing block data path")
	} else {
		glg.Info("Initializing block data path")
	}
}
