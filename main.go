package main

import (
	"flag"
	"github.com/jkusniar/lara/api"
	"github.com/jkusniar/lara/logger"
	"github.com/jkusniar/lara/msg"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	dispatcher *msg.Dispatcher
	logfile    string
	loglevel   string
)

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		if req, err := httputil.DumpRequest(r, true); err != nil {
			logger.Debugf("Error dumping http request: %s\n\n", err)
		} else {
			logger.Debugf("Processing HTTP request:\n\n %s \n\n", req)
		}
		next.ServeHTTP(w, r)

		logger.Debugf("Request processing took: %s\n", time.Since(start))
	})
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	msgHandlerFn, err := dispatcher.Dispatch(r.Body)
	if err != nil {
		logger.Errorf("Error parsing request body: %s\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := msgHandlerFn(w); err != nil {
		logger.Errorf("Error processing request: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func init() {
	flag.StringVar(&logfile, "logfile", "", "log to file")
	flag.StringVar(&loglevel, "loglevel", "info", "minimal log level (debug|info|warning)")
}

func signalHandler() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan,
		syscall.SIGINT,
		syscall.SIGKILL,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGABRT,
	)
	<-signalChan
	// cleanup (logger, db...)
	logger.Shutdown()
	os.Exit(0)
}

func main() {
	flag.Parse()

	go signalHandler()

	logger.Start(logfile, loglevel)
	defer logger.Shutdown()

	dispatcher = new(msg.Dispatcher)
	dispatcher.Registry = api.RegisterHandlers()

	handler := http.HandlerFunc(apiHandler)
	http.Handle("/api", logMiddleware(handler))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	logger.Fatal(http.ListenAndServe(":8080", nil))
}
