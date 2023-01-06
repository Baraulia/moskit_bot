package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"runtime"
)

type Logger struct {
	*logrus.Entry
}

func GetLogger() *Logger {
	logger := logrus.New()
	logger.SetReportCaller(true)
	logger.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			filename := path.Base(frame.File)
			return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s:%d", filename, frame.Line)
		},
		FullTimestamp: true,
		ForceColors:   true,
	}

	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.DebugLevel)

	return &Logger{logrus.NewEntry(logger)}
}
