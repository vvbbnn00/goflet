// Package log provides a simple logging utility for the application.
package log

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/vvbbnn00/goflet/config"
)

var (
	enabled bool
	level   Level
)

var basePath string
var sepReplace bool

// Level The log level
type Level int

const (
	// DebugLevel The debug level
	DebugLevel Level = iota
	// InfoLevel The info level
	InfoLevel
	// WarnLevel The warn level
	WarnLevel
	// ErrorLevel The error level
	ErrorLevel
	// FatalLevel The fatal level
	FatalLevel
)

func init() {
	basePath, _ = filepath.Abs(".")
	sepReplace = os.PathSeparator != '/'

	log.SetFlags(0) // Disable the default logger data
	conf := config.GofletCfg.LogConfig
	enabled = conf.Enabled
	switch conf.Level {
	case "debug":
		level = DebugLevel
	case "info":
		level = InfoLevel
	case "warn":
		level = WarnLevel
	case "error":
		level = ErrorLevel
	case "fatal":
		level = FatalLevel
	}
}

func levelString(level Level) string {
	switch level {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	}
	return ""
}

func getCallerInfo() string {
	_, file, line, ok := runtime.Caller(3) // Skip 3 stacks
	if ok {
		relativePath, err := filepath.Rel(basePath, file)
		if err != nil {
			return file
		}
		if sepReplace {
			relativePath = strings.ReplaceAll(
				relativePath,
				string(os.PathSeparator),
				"/")
		}
		return fmt.Sprintf("%s:%d", relativePath, line)
	}
	return "unknown"
}
func printf(lvl Level, format string, v ...interface{}) {
	if !enabled || lvl < level {
		return
	}

	msg := fmt.Sprintf(format, v...)
	callerInfo := getCallerInfo()

	if lvl >= ErrorLevel {
		msg += "\n" + string(debug.Stack())
	}

	log.Printf("%s %s: [%s] %s",
		time.Now().Format("2006/01/02 15:04:05"),
		callerInfo,
		levelString(lvl),
		msg)

	if lvl == FatalLevel {
		os.Exit(1)
	}
}

// RawPrintf prints the message without any formatting
func RawPrintf(format string, v ...interface{}) {
	if !enabled || level > InfoLevel {
		return
	}
	fmt.Printf(format, v...)
}

// Debug prints the debug message
func Debug(msg string) {
	printf(DebugLevel, msg)
}

// Info prints the info message
func Info(msg string) {
	printf(InfoLevel, msg)
}

// Warn prints the warn message
func Warn(msg string) {
	printf(WarnLevel, msg)
}

// Error prints the error message
func Error(msg string) {
	printf(ErrorLevel, msg)
}

// Fatal prints the fatal message
func Fatal(msg string) {
	printf(FatalLevel, msg)
}

// Debugf prints the debug message with the format
func Debugf(format string, v ...interface{}) {
	printf(DebugLevel, format, v...)
}

// Infof prints the info message with the format
func Infof(format string, v ...interface{}) {
	printf(InfoLevel, format, v...)
}

// Warnf prints the warn message with the format
func Warnf(format string, v ...interface{}) {
	printf(WarnLevel, format, v...)
}

// Errorf prints the error message with the format
func Errorf(format string, v ...interface{}) {
	printf(ErrorLevel, format, v...)
}

// Fatalf prints the fatal message with the format
func Fatalf(format string, v ...interface{}) {
	printf(FatalLevel, format, v...)
}

// Printf prints the message with the format
func Printf(format string, v ...interface{}) {
	Infof(format, v...)
}
