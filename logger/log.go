package logger

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

var (
	debug      *Log
	info       *Log
	warn       *Log
	error      *Log
	loggerfile *os.File
	level      uint8
	done       = make(chan bool)
)

const (
	DEBUG = "DEBUG: "
	INFO  = "INFO: "
	WARN  = "WARN: "
	ERROR = "ERROR: "
)

const (
	LEVEL_ERROR uint8 = iota
	LEVEL_WARN
	LEVEL_INFO
	LEVEL_DEBUG
)

type Log struct {
	LogChan chan string
	log     *log.Logger
	// TODO implement automatic file rotation using another channel
}

func newLog(w io.Writer, prefix string) *Log {
	return &Log{
		LogChan: make(chan string, 10),
		log:     log.New(w, prefix, log.LstdFlags|log.Lshortfile),
	}
}

func (l *Log) listen() {
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

func (l *Log) shutdown() {
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
		debug.LogChan <- fmt.Sprint(v...)
	}
}

func Debugf(format string, v ...interface{}) {
	if level >= LEVEL_DEBUG {
		debug.LogChan <- fmt.Sprintf(format, v...)
	}
}

func Info(v ...interface{}) {
	if level >= LEVEL_INFO {
		info.LogChan <- fmt.Sprint(v...)
	}
}

func Infof(format string, v ...interface{}) {
	if level >= LEVEL_INFO {
		info.LogChan <- fmt.Sprintf(format, v...)
	}
}

func Warn(v ...interface{}) {
	if level >= LEVEL_WARN {
		warn.LogChan <- fmt.Sprint(v...)
	}
}

func Warnf(format string, v ...interface{}) {
	if level >= LEVEL_WARN {
		warn.LogChan <- fmt.Sprintf(format, v...)
	}
}

func Error(v ...interface{}) {
	error.LogChan <- fmt.Sprint(v...)
}

func Errorf(format string, v ...interface{}) {
	error.LogChan <- fmt.Sprintf(format, v...)
}

func Panic(v ...interface{}) {
	msg := fmt.Sprint(v...)
	error.LogChan <- msg
	panic(msg)

}

func Panicf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	error.LogChan <- msg
	panic(msg)
}

func Fatal(v ...interface{}) {
	msg := fmt.Sprint(v...)
	error.LogChan <- msg
	os.Exit(1)

}

func Fatalf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	error.LogChan <- msg
	os.Exit(1)
}
