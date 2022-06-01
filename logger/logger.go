package logger

import (
	"io"

	"gopkg.in/natefinch/lumberjack.v2"
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
}

func GetLogger() *io.Writer {
	if logger == nil {
		InitLogger()
	}
	return &logger
}
