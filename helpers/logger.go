package helpers

import "github.com/kpango/glg"

type Log struct {
	Logger *glg.Glg
}

func NewLogger() *Log {
	l := &Log{}
	l.Logger = glg.Get().SetMode(glg.WRITER).AddLevelWriter(glg.INFO, glg.FileWriter("/tmp/errors.log", 0666))
	return l
}
