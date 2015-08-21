package main

import (
	"flag"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jkusniar/lara/api"
	"github.com/jkusniar/lara/app"
	"github.com/jkusniar/lara/auth"
	"github.com/jkusniar/lara/msg"
)

var (
	dispatcher *msg.Dispatcher
	logfile    string
	loglevel   string
)

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		if app.Log.DebugEnabled() {
			d, err := httputil.DumpRequest(r, true)
			if err != nil {
				app.Log.Errorf("Error dumping HTTP request: %s",
					err)
			} else {
				app.Log.Debugf("Processing HTTP request: %s", d)
			}
		}

		next.ServeHTTP(w, r)

		app.Log.Debugf("Request processing took: %s", time.Since(start))
	})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := auth.Validate(w, r); err != nil {
			app.Log.Errorf("Authorization failed: %s", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fn, err := dispatcher.Dispatch(r.Body)
	if err != nil {
		app.Log.Errorf("Error parsing request body: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := fn(w); err != nil {
		app.Log.Errorf("Error processing request: %s", err)
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
	flag.StringVar(&loglevel, "loglevel", "info",
		"minimal log level (debug|info|warning)")
}

// signaleHandler is goroutine handling termination signals. It shuts down
// application properly
func signalHandler() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		syscall.SIGINT,
		syscall.SIGKILL,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGABRT,
	)
	<-sigChan

	app.ShutdownLog()
	os.Exit(0)
}

func main() {
	flag.Parse()

	go signalHandler()

	app.StartLog(logfile, loglevel)
	defer app.ShutdownLog()

	dispatcher = msg.NewDispatcher(api.RegisterHandlers())

	// TODO allow only post requests to api/login/register
	// TODO validate content type prior to processing of requests
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
	http.Handle("/login",
		logMiddleware(http.HandlerFunc(loginHandler)))
	http.Handle("/register",
		logMiddleware(http.HandlerFunc(regHandler)))
	http.Handle("/api",
		logMiddleware(authMiddleware(http.HandlerFunc(apiHandler))))

	app.Log.Fatal(http.ListenAndServe(":8080", nil))
}
