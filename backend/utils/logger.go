package utils

// implement logger with info, error, warn levels

import (
	"log"
	"os"
)

type Logger struct {
	info  *log.Logger
	error *log.Logger
	warn  *log.Logger
}

func NewLogger() *Logger {
	return &Logger{
		info:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		error: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		warn:  log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (logger *Logger) Info(message string) {
	logger.info.Println(message)
}

func (logger *Logger) Error(message string) {
	logger.error.Println(message)
}

func (logger *Logger) Warn(message string) {
	logger.warn.Println(message)
}
