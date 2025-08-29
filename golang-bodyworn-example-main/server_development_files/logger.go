package server

import (
	"fmt"
	"sync"
)

// ServerLogger interface defines logging methods
type ServerLogger interface {
	Error(v ...interface{}) error
	Warning(v ...interface{}) error
	Info(v ...interface{}) error

	Errorf(format string, a ...interface{}) error
	Warningf(format string, a ...interface{}) error
	Infof(format string, a ...interface{}) error
}

var (
	logger      ServerLogger
	loggerMutex sync.Mutex
)

func InitLogger() {
	logger = &DefaultLogger{}
}



// SetLogger allows setting a custom logger
func SetLogger(l ServerLogger) {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()
	logger = l
}

// DefaultLogger provides a default implementation for logging
type DefaultLogger struct{}

type Level int

const (
	Error   Level = 3
	Warning Level = 4
	Info    Level = 6
)

// String converts the Level to a string
func (l Level) String() string {
	switch l {
	case Error:
		return "Error"
	case Warning:
		return "Warning"
	case Info:
		return "Info"
	}
	return "Unknown log level"
}

// log prints the log message with the specified level
func (l DefaultLogger) log(level Level, v ...interface{}) {
	fmt.Printf("%s: %s", level, fmt.Sprintln(v...))
}

// logf prints a formatted log message with the specified level
func (l DefaultLogger) logf(level Level, format string, v ...interface{}) {
	fmt.Printf("%s: %s\n", level, fmt.Sprintf(format, v...))
}

// Error logs an error message
func (l DefaultLogger) Error(v ...interface{}) error {
	l.log(Error, v...)
	return nil
}

// Warning logs a warning message
func (l DefaultLogger) Warning(v ...interface{}) error {
	l.log(Warning, v...)
	return nil
}

// Info logs an info message
func (l DefaultLogger) Info(v ...interface{}) error {
	l.log(Info, v...)
	return nil
}

// Errorf logs an error message with formatting
func (l DefaultLogger) Errorf(format string, a ...interface{}) error {
	l.logf(Error, format, a...)
	return nil
}

// Warningf logs a warning message with formatting
func (l DefaultLogger) Warningf(format string, a ...interface{}) error {
	l.logf(Warning, format, a...)
	return nil
}

// Infof logs an info message with formatting
func (l DefaultLogger) Infof(format string, a ...interface{}) error {
	l.logf(Info, format, a...)
	return nil
}
