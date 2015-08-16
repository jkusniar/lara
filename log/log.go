package log

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

const (
	DEBUG = "DEBUG"
	INFO  = "INFO"
	WARN  = "WARN"
	ERROR = "ERROR"
)

const (
	LEVEL_ERROR uint8 = iota
	LEVEL_WARN
	LEVEL_INFO
	LEVEL_DEBUG
)

type Logger struct {
	mutex   sync.Mutex
	entries chan logEntry
	done    chan bool
	writer  io.WriteCloser
	level   uint8
	// TODO implement automatic file rotation. https://github.com/natefinch/lumberjack
}

func NewLogger(w io.WriteCloser, level uint8) *Logger {
	l := &Logger{
		entries: make(chan logEntry, 10),
		done:    make(chan bool),
		writer:  w,
		level:   level,
	}

	go l.listen()

	return l
}

func (l *Logger) listen() {
	for {
		entry, more := <-l.entries
		if more {
			fmt.Fprintln(l.writer, entry)
			// FIXME better use: l.Writer.Write()
		} else {
			l.done <- true
			return
		}
	}
}

func (l *Logger) Shutdown() {
	fmt.Println("closing channel")
	close(l.entries)
	fmt.Println("waiting for channel to process logs")
	<-l.done

	if l.writer != nil {
		fmt.Println("closing logfile")
		l.writer.Close()
	}
}

func (l *Logger) CanLog(level uint8) bool {
	// TODO should lock with mutex before reading?
	return l.level >= level && l.writer != nil
}

func (l *Logger) SetLevel(level uint8) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.level = level
}

func (l *Logger) Debug(v ...interface{}) {
	if l.CanLog(LEVEL_DEBUG) {
		l.entries <- createLogEntry(DEBUG, v...)
	}
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.CanLog(LEVEL_DEBUG) {
		l.entries <- createLogEntryf(DEBUG, format, v...)
	}
}

func (l *Logger) Info(v ...interface{}) {
	if l.CanLog(LEVEL_INFO) {
		l.entries <- createLogEntry(INFO, v...)
	}
}

func (l *Logger) Infof(format string, v ...interface{}) {
	if l.CanLog(LEVEL_INFO) {
		l.entries <- createLogEntryf(INFO, format, v...)
	}
}

func (l *Logger) Warn(v ...interface{}) {
	if l.CanLog(LEVEL_WARN) {
		l.entries <- createLogEntry(WARN, v...)
	}
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	if l.CanLog(LEVEL_WARN) {
		l.entries <- createLogEntryf(WARN, format, v...)
	}
}

func (l *Logger) Error(v ...interface{}) {
	if l.CanLog(LEVEL_ERROR) {
		l.entries <- createLogEntry(ERROR, v...)
	}
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	if l.CanLog(LEVEL_ERROR) {
		l.entries <- createLogEntryf(ERROR, format, v...)
	}
}

func (l *Logger) Panic(v ...interface{}) {
	entry := createLogEntry(ERROR, v...)
	if l.CanLog(LEVEL_ERROR) {
		l.entries <- entry
	}

	panic(entry.Message)

}

func (l *Logger) Panicf(format string, v ...interface{}) {
	entry := createLogEntryf(ERROR, format, v...)
	if l.CanLog(LEVEL_ERROR) {
		l.entries <- entry
	}

	panic(entry.Message)
}

func (l *Logger) Fatal(v ...interface{}) {
	if l.CanLog(LEVEL_ERROR) {
		l.entries <- createLogEntry(ERROR, v...)
	}

	os.Exit(1)

}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	if l.CanLog(LEVEL_ERROR) {
		l.entries <- createLogEntryf(ERROR, format, v...)
	}

	os.Exit(1)
}

type logEntry struct {
	Level    string
	Message  string
	Filename string
	Line     int
	Time     time.Time
}

func (e logEntry) String() string {
	return fmt.Sprintf("[%v] [%v:%v] [%v] %v",
		e.Time, e.Filename, e.Line, e.Level, e.Message)
}

func createLogEntryf(level, format string, v ...interface{}) logEntry {
	now := time.Now()
	file, line := callerInfo()
	return logEntry{
		Level:    level,
		Message:  fmt.Sprintf(format, v...),
		Filename: file,
		Line:     line,
		Time:     now,
	}
}

func createLogEntry(level string, v ...interface{}) logEntry {
	now := time.Now()
	file, line := callerInfo()
	return logEntry{
		Level:    level,
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
		file = "???.go"
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
