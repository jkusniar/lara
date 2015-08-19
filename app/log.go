// Application wide services and utilities
package app

import (
	"fmt"
	log "github.com/jkusniar/go-log"
	"io"
	"os"
)

var (
	Log *log.Logger
)

//better name
func StartLog(logfile, loglevel string) {
	var writer io.WriteCloser
	if logfile != "" {
		file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open log file %v: %v", logfile, err)
			os.Exit(1)
		}

		writer = file
	}

	var min_level uint8
	switch loglevel {
	case "debug":
		min_level = log.LevelDebug
	case "info":
		min_level = log.LevelInfo
	case "warn":
		min_level = log.LevelWarn
	default:
		min_level = log.LevelError
	}

	Log = log.New(writer, min_level)
}

// better name
func ShutdownLog() {
	Log.Shutdown()
}
