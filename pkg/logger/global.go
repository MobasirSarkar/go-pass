package logger

import (
	"io"
	"log"
	"os"
)

var LoggerInstance *Logger

func init() {
	LoggerInstance = NewLogger(os.Stdout)
}

func NewLogger(StdOut io.Writer) *Logger {
	logger := &Logger{
		STDOUT:   os.Stdout,
		STDERR:   os.Stderr,
		debug:    false,
		showTime: true,
		saveLogs: true,
		log:      log.New(StdOut, "", 0),
	}
	if logger.saveLogs {
		err := os.MkdirAll("logs", os.ModePerm)
		if err != nil {
			log.Printf("Error Creating log directory: %s \n", err)
		} else {
			logfile, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				log.Println("Error creating log file:", err) // Properly log the error
			} else {
				logger.logFile = logfile
				logger.log.SetOutput(io.MultiWriter(StdOut, logfile)) // Log to both console and file
			}
		}
	}
	return logger
}

func Info(format string, args ...interface{}) {
	LoggerInstance.Info(format, args...)
}

func Good(format string, args ...interface{}) {
	LoggerInstance.Good(format, args...)
}

func Debug(format string, args ...interface{}) {
	LoggerInstance.Debug(format, args...)
}

func DebugError(format string, args ...interface{}) {
	LoggerInstance.DebugError(format, args...)
}

func Warn(format string, args ...interface{}) {
	LoggerInstance.Warn(format, args...)
}

func Error(format string, args ...interface{}) {
	LoggerInstance.Error(format, args...)
}

func Fatal(format string, args ...interface{}) {
	LoggerInstance.Fatal(format, args...)
}

func Panic(format string, args ...interface{}) {
	LoggerInstance.Panic(format, args...)
}

func SetDebug(enable bool) {
	LoggerInstance.SetDebug(enable)
}

func ShowTime(enable bool) {
	LoggerInstance.ShowTime(enable)
}

func SetStdOut(w io.Writer) {
	LoggerInstance.log.SetOutput(w)
	if LoggerInstance.saveLogs && LoggerInstance.logFile != nil {
		LoggerInstance.log.SetOutput(io.MultiWriter(w, LoggerInstance.logFile))
	}
}
