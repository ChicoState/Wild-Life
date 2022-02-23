package log

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

const (
	name       = "WildLife"
	colorGreen = "\u001B[1;32m"
	colorReset = "\u001B[m"
)

// printLn prints a formatted string to stdout
func printLn(line string) {
	fmt.Printf("%s[%s]%s %s%s\n", colorGreen, name, colorReset, line, colorReset)
}

// Logf will attempt to format a log message or print plain text.
func Logf(format string, args ...interface{}) {
	// If debug tracing is enabled, the file and line number of a log function call will prefix the log message.
	if os.Getenv("trace") == "true" {
		// Get function caller
		_, file, ln, ok := runtime.Caller(1)
		// Add the caller file and line number prefix
		if ok {
			fmt.Printf("%s:%d -> ", filepath.Base(file), ln)
		}
	}
	// Print the log message
	printLn(fmt.Sprintf(format, args...))
}
