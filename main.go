package main

import (
	"github.com/jkusniar/lara/api"
	"github.com/jkusniar/lara/msg"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

var (
	dispatcher *msg.Dispatcher
)

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		if req, err := httputil.DumpRequest(r, true); err != nil {
			log.Printf("Error dumping http request: %s\n\n", err)
		} else {
			log.Printf("Processing HTTP request:\n\n %s \n\n", req)
		}
		next.ServeHTTP(w, r)

		log.Printf("Request processing took: %s\n", time.Since(start))
	})
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	msgHandlerFn, err := dispatcher.Dispatch(r.Body)
	if err != nil {
		log.Printf("Error parsing request body: %s\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := msgHandlerFn(w); err != nil {
		log.Printf("Error processing request: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	dispatcher = new(msg.Dispatcher)
	dispatcher.Registry = api.RegisterHandlers()
}

func main() {
	handler := http.HandlerFunc(apiHandler)
	http.Handle("/api", logMiddleware(handler))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
