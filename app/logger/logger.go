package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

type FileLogger struct {
	out    io.WriteCloser
	logger *log.Logger
}

var logFlags = log.Ldate | log.Ltime | log.Lmicroseconds

const AppName = "xBrowser"

var (
	AppLogDir     string
	AppLogDirErr  error
	AppLogName    string
	AppErrLogName string
)

func init() {
	logName := AppName + "-" + time.Now().Local().Format("2006-01-02")
	AppLogName = logName + ".log"
	AppErrLogName = logName + "_error.log"
	AppLogDir, AppLogDirErr = getLogDir()
}

// NewDefaultLogger creates a new Logger.
func NewAppLogger() (*FileLogger, error) {
	return NewLogger(AppLogName)
}

func NewAppErrLogger() (*FileLogger, error) {
	return NewLogger(AppErrLogName)
}

func NewLogger(logName string) (*FileLogger, error) {

	f, err := os.Create(AppLogDir + logName)
	if err != nil {
		return nil, fmt.Errorf("Failed to open log file %s error: %s", AppLogDir+logName, err)
	}
	return &FileLogger{
		out:    f,
		logger: log.New(f, "", logFlags),
	}, nil
}

func getLogDir() (logDir string, err error) {
	home, _ := os.UserHomeDir()
	switch runtime.GOOS {
	case "darwin":
		logDir = fmt.Sprintf("%s/Library/Logs/%s/", home, AppName)
	case "windows":
		logDir = fmt.Sprintf("%s\\%s\\logs\\", home, AppName)
	case "linux":
		logDir = fmt.Sprintf("%s/%s/logs/", home, AppName)
	default:
		return "", fmt.Errorf("not supported on this platform: %s", runtime.GOOS)
	}
	if err = os.MkdirAll(logDir, 0755); err != nil {
		return "", fmt.Errorf("Failed to create dir %s. error: %s", logDir, err)
	}
	fmt.Println("log dir:", logDir)
	return logDir, nil
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
