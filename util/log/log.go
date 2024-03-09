package log

import (
	"fmt"
	"goflet/config"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

var (
	enabled bool
	level   Level
)

var basePath string
var sepReplace bool

type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
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
			return "unknown"
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

func RawPrintf(format string, v ...interface{}) {
	if !enabled || level > InfoLevel {
		return
	}
	fmt.Printf(format, v...)
}

func Debug(msg string) func(string, ...interface{}) {
	return func(s string, i ...interface{}) {
		printf(DebugLevel, s, i)
	}
}

func Info(msg string) {
	printf(InfoLevel, msg)
}

func Warn(msg string) {
	printf(WarnLevel, msg)
}

func Error(msg string) {
	printf(ErrorLevel, msg)
}

func Fatal(msg string) {
	printf(FatalLevel, msg)
}

func Debugf(format string, v ...interface{}) {
	printf(DebugLevel, format, v...)
}

func Infof(format string, v ...interface{}) {
	printf(InfoLevel, format, v...)
}

func Warnf(format string, v ...interface{}) {
	printf(WarnLevel, format, v...)
}

func Errorf(format string, v ...interface{}) {
	printf(ErrorLevel, format, v...)
}

func Fatalf(format string, v ...interface{}) {
	printf(FatalLevel, format, v...)
}

func Printf(format string, v ...interface{}) {
	Infof(format, v...)
}
