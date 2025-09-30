package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

type Logger struct {
	level  LogLevel
	output io.Writer
	logger *log.Logger
}

type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	File      string    `json:"file,omitempty"`
	Line      int       `json:"line,omitempty"`
	Function  string    `json:"function,omitempty"`
	Data      any       `json:"data,omitempty"`
}

func NewLogger(level LogLevel, output io.Writer) *Logger {
	if output == nil {
		output = os.Stdout
	}

	return &Logger{
		level:  level,
		output: output,
		logger: log.New(output, "", 0),
	}
}

func (l *Logger) Debug(message string, data ...any) {
	l.log(DEBUG, message, data...)
}

func (l *Logger) Info(message string, data ...any) {
	l.log(INFO, message, data...)
}

func (l *Logger) Warn(message string, data ...any) {
	l.log(WARN, message, data...)
}

func (l *Logger) Error(message string, data ...any) {
	l.log(ERROR, message, data...)
}

func (l *Logger) Fatal(message string, data ...any) {
	l.log(FATAL, message, data...)
	os.Exit(1)
}

func (l *Logger) log(level LogLevel, message string, data ...any) {
	if level < l.level {
		return
	}

	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     l.getLevelString(level),
		Message:   message,
	}

	if level >= WARN {
		if pc, file, line, ok := runtime.Caller(2); ok {
			entry.File = file
			entry.Line = line
			if fn := runtime.FuncForPC(pc); fn != nil {
				entry.Function = fn.Name()
			}
		}
	}

	if len(data) > 0 {
		entry.Data = data
	}

	jsonData, _ := json.Marshal(entry)
	l.logger.Println(string(jsonData))
}

func (l *Logger) getLevelString(level LogLevel) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

func (l *Logger) SetOutput(output io.Writer) {
	l.output = output
	l.logger = log.New(output, "", 0)
}

func (l *Logger) Clear() {
	// This is a placeholder method for clearing logs
	// In a real implementation, this might clear log files or reset log buffers
	l.Info("Logs cleared")
}

var defaultLogger = NewLogger(INFO, os.Stdout)

func SetDefaultLogger(logger *Logger) {
	defaultLogger = logger
}

func Debug(message string, data ...any) {
	defaultLogger.Debug(message, data...)
}

func Info(message string, data ...any) {
	defaultLogger.Info(message, data...)
}

func Warn(message string, data ...any) {
	defaultLogger.Warn(message, data...)
}

func Error(message string, data ...any) {
	defaultLogger.Error(message, data...)
}

func Fatal(message string, data ...any) {
	defaultLogger.Fatal(message, data...)
}

func GetLogger() *Logger {
	return defaultLogger
}

func NewFileLogger(filename string, level LogLevel) (*Logger, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return NewLogger(level, file), nil
}

func NewRotatingLogger(baseFilename string, level LogLevel, maxSize int64) (*Logger, error) {
	filename := fmt.Sprintf("%s.%s.log", baseFilename, time.Now().Format("2006-01-02"))

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	logger := NewLogger(level, file)

	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			file.Close()

			filename = fmt.Sprintf("%s.%s.log", baseFilename, time.Now().Format("2006-01-02"))
			file, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				Error("Failed to rotate log file", "error", err)
				continue
			}

			logger.SetOutput(file)
		}
	}()

	return logger, nil
}
