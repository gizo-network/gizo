package core

import (
	"os"

	"github.com/kpango/glg"
)

func InitializeDataPath() {
	os.Mkdir(IndexPath, os.FileMode(0777))
	os.Mkdir(BlockPath, os.FileMode(0777))

}

func RemoveDataPath() {
	err := os.RemoveAll(IndexPath)
	if err != nil {
		glg.Fatal(err)
	}
}
