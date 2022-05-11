package log

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

const (
	name       = "WildLife"
	colorGreen = "\u001B[1;32m"
	colorRed   = "\u001B[1;31m"
	colorBlue  = "\u001B[1;34m"
	colorReset = "\u001B[m"
)

// printLn prints a formatted string to stdout
func printLn(line string) {
	fmt.Printf("%s[%s]%s %s%s\n", colorGreen, name, colorReset, line, colorReset)
}

func errorLn(line string) {
	fmt.Fprintf(os.Stderr, "%s[Error]%s %s%s\n", colorRed, colorReset, line, colorReset)
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

func Errf(format string, args ...interface{}) {
	// If debug tracing is enabled, the file and line number of a log function call will prefix the log message.
	if os.Getenv("trace") == "true" {
		// Get function caller
		_, _, ln, ok := runtime.Caller(1)
		// Add the caller file and line number prefix
		if ok {
			fmt.Printf("%s[Error]%s:%d -> ", colorRed, colorReset, ln)
		}
	}
	// Print the log message
	errorLn(fmt.Sprintf(format, args...))
}

// Middleware provides an interface for logging http requests
func Middleware(next http.Handler) http.Handler {
	// Create http Handler function
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Start timing the request
		t1 := time.Now()
		// Log after next.ServeHTTP exits
		defer func() {
			fmt.Printf("%s[%s]%s %s %s (%s) %s\n", colorGreen, r.Method, colorReset, r.RemoteAddr, r.URL.Path,
				time.Since(t1), colorReset)
		}()
		// Serve request
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
