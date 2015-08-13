package logger

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"time"
)

var (
	debug      Log
	info       Log
	warn       Log
	error      Log
	loggerfile *os.File
	level      uint8
	done       = make(chan bool)
)

const (
	DEBUG = "DEBUG "
	INFO  = "INFO "
	WARN  = "WARN "
	ERROR = "ERROR "
)

const (
	LEVEL_ERROR uint8 = iota
	LEVEL_WARN
	LEVEL_INFO
	LEVEL_DEBUG
)

type Log struct {
	LogChan chan logEntry
	log     *log.Logger
	// TODO implement automatic file rotation using another channel
}

type logEntry struct {
	Message  string
	Filename string
	Line     int
	Time     time.Time
}

func (entry logEntry) String() string {
	return fmt.Sprintf("%v %v:%v: %v", entry.Time, entry.Filename, entry.Line, entry.Message)
}

func newLog(w io.Writer, prefix string) Log {
	return Log{
		LogChan: make(chan logEntry, 10),
		log:     log.New(w, prefix, 0),
	}
}

func (l Log) listen() {
	for {
		msg, more := <-l.LogChan
		if more {
			l.log.Println(msg)
		} else {
			done <- true
			return
		}
	}
}

func (l Log) shutdown() {
	close(l.LogChan)
}

func Start(logfile, loglevel string) {
	stdout := io.MultiWriter(ioutil.Discard, os.Stdout)
	stderr := io.MultiWriter(ioutil.Discard, os.Stderr)

	if logfile != "" {
		file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalln("Failed to open log file ", logfile, ":", err)
		}

		stdout = io.MultiWriter(file, os.Stdout)
		stderr = io.MultiWriter(file, os.Stderr)
		loggerfile = file
	}

	switch loglevel {
	case "debug":
		level = LEVEL_DEBUG
	case "info":
		level = LEVEL_INFO
	case "warn":
		level = LEVEL_WARN
	default:
		level = LEVEL_ERROR
	}

	debug = newLog(stdout, DEBUG)
	info = newLog(stdout, INFO)
	warn = newLog(stdout, WARN)
	error = newLog(stderr, ERROR)

	go debug.listen()
	go info.listen()
	go warn.listen()
	go error.listen()
}

func Shutdown() {
	debug.shutdown()
	<-done

	info.shutdown()
	<-done

	warn.shutdown()
	<-done

	error.shutdown()
	<-done

	if loggerfile != nil {
		loggerfile.Close()
	}
}

func Debug(v ...interface{}) {
	if level >= LEVEL_DEBUG {
		debug.LogChan <- createLogEntry(v...)
	}
}

func Debugf(format string, v ...interface{}) {
	if level >= LEVEL_DEBUG {
		debug.LogChan <- createLogEntryf(format, v...)
	}
}

func Info(v ...interface{}) {
	if level >= LEVEL_INFO {
		info.LogChan <- createLogEntry(v...)
	}
}

func Infof(format string, v ...interface{}) {
	if level >= LEVEL_INFO {
		info.LogChan <- createLogEntryf(format, v...)
	}
}

func Warn(v ...interface{}) {
	if level >= LEVEL_WARN {
		warn.LogChan <- createLogEntry(v...)
	}
}

func Warnf(format string, v ...interface{}) {
	if level >= LEVEL_WARN {
		warn.LogChan <- createLogEntryf(format, v...)
	}
}

func Error(v ...interface{}) {
	error.LogChan <- createLogEntry(v...)
}

func Errorf(format string, v ...interface{}) {
	error.LogChan <- createLogEntryf(format, v...)
}

func Panic(v ...interface{}) {
	entry := createLogEntry(v...)
	error.LogChan <- entry
	panic(entry.Message)

}

func Panicf(format string, v ...interface{}) {
	entry := createLogEntryf(format, v...)
	error.LogChan <- entry
	panic(entry.Message)
}

func Fatal(v ...interface{}) {
	error.LogChan <- createLogEntry(v...)
	os.Exit(1)

}

func Fatalf(format string, v ...interface{}) {
	error.LogChan <- createLogEntryf(format, v...)
	os.Exit(1)
}

func createLogEntryf(format string, v ...interface{}) logEntry {
	now := time.Now()
	file, line := callerInfo()
	return logEntry{
		Message:  fmt.Sprintf(format, v...),
		Filename: file,
		Line:     line,
		Time:     now,
	}
}

func createLogEntry(v ...interface{}) logEntry {
	now := time.Now()
	file, line := callerInfo()
	return logEntry{
		Message:  fmt.Sprint(v...),
		Filename: file,
		Line:     line,
		Time:     now,
	}
}

func callerInfo() (file string, line int) {
	var ok bool
	_, file, line, ok = runtime.Caller(3)
	if !ok {
		file = "???"
		line = 0
	}

	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	file = short

	return
}
