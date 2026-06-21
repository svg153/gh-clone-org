package cmd

import (
	"fmt"
	"os"
	"time"
)

// LogLevel represents the verbosity level
type LogLevel int

const (
	LogSilent LogLevel = iota // 0: no output
	LogInfo                   // 1: clone status per repo
	LogVerbose                // 2: git output
	LogDebug                  // 3: full debug info (API calls, timing)
)

// Logger provides structured logging with levels
type Logger struct {
	level LogLevel
}

// NewLogger creates a logger with the given verbosity level
func NewLogger(level LogLevel) *Logger {
	return &Logger{level: level}
}

// Info logs at level 1+ (clone status)
func (l *Logger) Info(msg string) {
	if l.level >= LogInfo {
		fmt.Println(msg)
	}
}

// Verbose logs at level 2+ (git output)
func (l *Logger) Verbose(msg string) {
	if l.level >= LogVerbose {
		fmt.Println(msg)
	}
}

// Debug logs at level 3+ (API calls, timing)
func (l *Logger) Debug(msg string) {
	if l.level >= LogDebug {
		fmt.Println(msg)
	}
}

// Println always prints (used for critical output)
func (l *Logger) Println(msg string) {
	fmt.Println(msg)
}

// PrintError always prints errors
func (l *Logger) PrintError(msg string) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", msg)
}

// Timer provides simple timing
type Timer struct {
	start time.Time
	name  string
}

// NewTimer creates a new timer
func NewTimer(name string) *Timer {
	return &Timer{start: time.Now(), name: name}
}

// Stop logs the elapsed time at debug level
func (t *Timer) Stop(l *Logger) {
	elapsed := time.Since(t.start)
	l.Debug(fmt.Sprintf("Timer %s: %v", t.name, elapsed))
}
