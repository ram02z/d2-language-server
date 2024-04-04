package log

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type Logger struct {
	*log.Logger
}

func NewLogger(prefix string) *Logger {
	var logDir string
	if runtime.GOOS == "windows" {
		logDir = filepath.Join(os.Getenv("APPDATA"), prefix, "logs")
	} else {
		logDir = filepath.Join(os.Getenv("HOME"), ".local", "state", prefix,
			"logs")
	}

	err := os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	currentTime := time.Now()
	fileName := filepath.Join(logDir, fmt.Sprintf("%s-debug.log",
		currentTime.Format("2006-01-02-15-04-05")))

	logfile, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("file does not exist")
	}

	return &Logger{
		Logger: log.New(logfile, fmt.Sprintf("[%s] ", prefix),
			log.Ldate|log.Ltime|log.Lshortfile),
	}
}
