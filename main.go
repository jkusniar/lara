package main

import (
	"flag"
	"github.com/jkusniar/lara/api"
	"github.com/jkusniar/lara/app"
	"github.com/jkusniar/lara/auth"
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
			app.Log.Debugf("Error dumping http request: %s\n\n", err)
		} else {
			app.Log.Debugf("Processing HTTP request:\n\n %s \n\n", req)
		}
		next.ServeHTTP(w, r)

		app.Log.Debugf("Request processing took: %s\n", time.Since(start))
	})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := auth.Validate(w, r); err != nil {
			app.Log.Errorf("Authorization failed: %s\n", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	msgHandlerFn, err := dispatcher.Dispatch(r.Body)
	if err != nil {
		app.Log.Errorf("Error parsing request body: %s\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := msgHandlerFn(w); err != nil {
		app.Log.Errorf("Error processing request: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	auth.Login(w, r)
}

func regHandler(w http.ResponseWriter, r *http.Request) {
	auth.Register(w, r)
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
	// cleanup (log, db...)
	app.ShutdownLog()
	os.Exit(0)
}

func main() {
	flag.Parse()

	go signalHandler()

	app.StartLog(logfile, loglevel)
	defer app.ShutdownLog()

	dispatcher = new(msg.Dispatcher)
	dispatcher.Registry = api.RegisterHandlers()

	// TODO, allow only post requests to api/login/register
	// TODO validate content type prior to processing of requests
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
	http.Handle("/login", logMiddleware(http.HandlerFunc(loginHandler)))
	http.Handle("/register", logMiddleware(http.HandlerFunc(regHandler)))
	http.Handle("/api", logMiddleware(authMiddleware(http.HandlerFunc(apiHandler))))

	// TODO refactor logger to logger.Log.Fatal....
	app.Log.Fatal(http.ListenAndServe(":8080", nil))
}
