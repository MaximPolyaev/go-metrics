package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

func New(f *os.File) *Logger {
	lg := logrus.New()
	lg.SetOutput(f)
	lg.SetLevel(logrus.InfoLevel)

	return &Logger{Logger: lg}
}
