package inits

import (
	"io"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func InitLog() *os.File {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})

	logFile, err := os.OpenFile(filepath.Join("log.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Panicf("Error opening log file: %s", err)
	}

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	return logFile // <-- отдаём наверх, чтобы закрыть в main
}
