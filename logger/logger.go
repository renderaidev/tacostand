/* -----------------------------------------------------------------------------
 * Copyright (c) Somus OÜ. All rights reserved.
 * This software is licensed under the MIT license.
 * See the LICENSE file for further information.
 * -------------------------------------------------------------------------- */

package logger

import (
	"log"

	"github.com/gookit/color"
)

// Logger is a wrapper around a log.Logger that adds some additional
// functionality.
type Logger struct {
	*log.Logger

	debugMode bool
}

// NewLogger creates a new Logger object with a provided Go logger.
func NewLogger(logger *log.Logger, debugMode bool) *Logger {
	l := &Logger{
		Logger:    logger,
		debugMode: debugMode,
	}

	return l
}

func (l *Logger) withPrefix(prefix string, v ...interface{}) []interface{} {
	params := append([]interface{}{prefix}, v...)
	return params
}

func (l *Logger) withColor(c color.Color, v ...interface{}) string {
	return c.Render(v...)
}

// Debug logs a message at the debug level. These logs will only appear in the
// console if the `DEBUG_MODE` environment variable is set to `true`.
func (l *Logger) Debug(v ...interface{}) {
	if l.debugMode {
		l.Println(l.withPrefix("[DEBUG] ", v...)...)
	}
}

// Info logs an informational message. It will be colored blue.
func (l *Logger) Info(v ...interface{}) {
	l.Println(l.withColor(color.Blue, l.withPrefix("[INFO] ", v...)...))
}

// Warn logs a warning message. It will be colored yellow.
func (l *Logger) Warn(v ...interface{}) {
	l.Println(l.withColor(color.Yellow, l.withPrefix("[WARN] ", v...)...))
}

// Danger logs an error message. It will be colored red.
func (l *Logger) Danger(v ...interface{}) {
	l.Println(l.withColor(color.Red, l.withPrefix("[DANGER] ", v...)...))
}

// Success logs a success message. It will be colored green.
func (l *Logger) Success(v ...interface{}) {
	l.Println(l.withColor(color.Green, l.withPrefix("[SUCCESS] ", v...)...))
}
