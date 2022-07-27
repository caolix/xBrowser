package logger

import (
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
)

type FileLogger struct {
	out    io.WriteCloser
	logger *log.Logger
}

var logFlags = log.Ldate | log.Ltime | log.Lmicroseconds

// NewDefaultLogger creates a new Logger.
func NewFileLogger(path string) logger.Logger {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic("Failed to open log file " + path)
	}
	l := &FileLogger{
		out:    f,
		logger: log.New(f, "", logFlags),
	}
	return l
}

func getCaller(skipCallDepth int) string {
	_, fullPath, line, ok := runtime.Caller(skipCallDepth)
	if !ok {
		return ""
	}
	fileParts := strings.Split(fullPath, "/")
	file := fileParts[len(fileParts)-1]
	return fmt.Sprintf("%s:%d", file, line)
}

func (l *FileLogger) prefixArray() []interface{} {
	array := make([]interface{}, 0, 3)
	array = append(array, getCaller(3))
	return array
}

// Print works like Sprintf.
func (l *FileLogger) Print(message string) {
	prefixArray := l.prefixArray()
	prefixArray = append(prefixArray, "PRN")
	l.logger.Println(prefixArray, message)
}

// Trace level logging. Works like Sprintf.
func (l *FileLogger) Trace(message string) {
	prefixArray := l.prefixArray()
	prefixArray = append(prefixArray, "TRA")
	l.logger.Println(prefixArray, message)
}

// Debug level logging. Works like Sprintf.
func (l *FileLogger) Debug(message string) {
	prefixArray := l.prefixArray()
	prefixArray = append(prefixArray, "DEB")
	l.logger.Println(prefixArray, message)
}

// Info level logging. Works like Sprintf.
func (l *FileLogger) Info(message string) {
	prefixArray := l.prefixArray()
	prefixArray = append(prefixArray, "INF")
	l.logger.Println(prefixArray, message)
}

// Warning level logging. Works like Sprintf.
func (l *FileLogger) Warning(message string) {
	prefixArray := l.prefixArray()
	prefixArray = append(prefixArray, "WAR")
	l.logger.Println(prefixArray, message)
}

// Error level logging. Works like Sprintf.
func (l *FileLogger) Error(message string) {
	prefixArray := l.prefixArray()
	prefixArray = append(prefixArray, "ERR")
	l.logger.Println(prefixArray, message)
}

// Fatal level logging. Works like Sprintf.
func (l *FileLogger) Fatal(message string) {
	prefixArray := l.prefixArray()
	prefixArray = append(prefixArray, "FAT")
	l.logger.Println(prefixArray, message)
	os.Exit(1)
}
