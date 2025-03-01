package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/MobasirSarkar/pass-manage/pkg/colors"
)

type Logger struct {
	STDOUT   *os.File
	STDERR   *os.File
	log      *log.Logger
	showTime bool
	debug    bool
	saveLogs bool
	logFile  *os.File
}

func FunctionTrace() (string, int) {
	caller := make([]uintptr, 15)
	callNums := runtime.Callers(2, caller)
	frames := runtime.CallersFrames(caller[:callNums])

	for {
		frame, more := frames.Next()
		if !strings.Contains(frame.File, "/logger/") {
			fileParts := strings.Split(frame.File, "/")
			fileName := fileParts[len(fileParts)-1]
			return fileName, frame.Line
		}
		if !more {
			break
		}
	}
	return "unknown", 0
}

func (logger *Logger) formatMessage(level string, format string, args ...interface{}) (string, string) {
	loc, _ := time.LoadLocation("Local")
	timestamp := time.Now().In(loc).Format("02 Jan 2006, 03:04:05 PM IST")
	filePath, line := FunctionTrace()
	fileParts := strings.Split(filePath, "/")
	fileName := fileParts[len(fileParts)-1]

	coloredMessage := fmt.Sprintf("[%s] [%s] [%s:%d] - %s",
		colors.Green(timestamp), colors.Blue(level), colors.Yellow(fileName), line, fmt.Sprintf(format, args...))

	plainMessage := fmt.Sprintf("[%s] [%s] [%s:%d] - %s",
		timestamp, level, fileName, line, fmt.Sprintf(format, args...))

	return coloredMessage, plainMessage
}

func (logger *Logger) logMessage(level string, format string, args ...interface{}) {
	coloredMsg, plainMsg := logger.formatMessage(level, format, args...)

	// Print colored logs to console
	fmt.Fprintln(logger.STDOUT, coloredMsg)

	// Save plain logs to file
	if logger.saveLogs && logger.logFile != nil {
		logger.logFile.WriteString(plainMsg + "\n")
	}
}

func (logger *Logger) Info(format string, args ...interface{}) {
	logger.logMessage("INFO", format, args...)
}

func (logger *Logger) Good(format string, args ...interface{}) {
	logger.logMessage("GOOD", format, args...)
}

func (logger *Logger) Debug(format string, args ...interface{}) {
	if logger.debug {
		logger.logMessage("DBUG", format, args...)
	}
}

func (logger *Logger) DebugError(format string, args ...interface{}) {
	if logger.debug {
		logger.logMessage("DBER", format, args...)
	}
}

func (logger *Logger) Warn(format string, args ...interface{}) {
	logger.logMessage("WARN", format, args...)
}

func (logger *Logger) Error(format string, args ...interface{}) {
	logger.logMessage("ERROR", format, args...)
}

func (logger *Logger) Fatal(format string, args ...interface{}) {
	logger.logMessage("FATAL", format, args...)
	os.Exit(1)
}

func (logger *Logger) Panic(format string, args ...interface{}) {
	message, plainMsg := logger.formatMessage("PANIC", format, args...)
	logger.log.Println(message)
	if logger.saveLogs && logger.logFile != nil {
		logger.logFile.WriteString(plainMsg + "\n")
	}
	panic(message)
}

func (logger *Logger) SetDebug(enable bool) {
	logger.debug = enable
}

func (logger *Logger) ShowTime(enable bool) {
	logger.showTime = enable
}
