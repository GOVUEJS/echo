package logger

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log"
	"os"
)

var logger io.Writer

func InitLogger() {
	logger = &lumberjack.Logger{
		Filename:   "/var/log/HWISECHO/echo.log",
		MaxSize:    500, // megabytes
		MaxBackups: 28,
		MaxAge:     28,    //days
		Compress:   false, // disabled by default
		LocalTime:  true,
	}
	log.SetOutput(io.MultiWriter(os.Stdout, logger))
}

func GetLogger() *io.Writer {
	if logger == nil {
		InitLogger()
	}
	return &logger
}
