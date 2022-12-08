package log

import (
	"log"
	"os"
)

type Logger struct {
	info  *log.Logger
	error *log.Logger
	debug *log.Logger
}

func NewLogger() Logger {
	return Logger{
		info:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime),
		error: log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime),
		debug: log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime),
	}
}

func (l *Logger) Error(file string, method string, message string) {
	l.error.Println(file+": "+method, ": ", message)
}

func (l *Logger) Info(file string, method string, message string) {
	l.info.Println(file + ": " + method + ": " + message)
}

func (l *Logger) Debug(message string) {
	l.debug.Println(message)
}
