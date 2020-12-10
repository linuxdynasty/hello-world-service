package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/netdata/statsd"
)

// statusCodeReporter is a compatible `http.ResponseWriter`
// which stores the `statusCode` for further reporting.
type statusCodeReporter struct {
	http.ResponseWriter
	written    bool
	statusCode int
}

func (w *statusCodeReporter) WriteHeader(statusCode int) {
	if w.written {
		return
	}

	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *statusCodeReporter) Write(b []byte) (int, error) {
	w.written = true
	return w.ResponseWriter.Write(b)
}

func main() {
	statsdExporter := os.Getenv("STATSD_EXPORTER")
	statsWriter, err := statsd.UDP(statsdExporter)
	if err != nil {
		panic(err)
	}

	statsD := statsd.NewClient(statsWriter, "helloworld.")
	statsD.FlushEvery(5 * time.Second)

	statsDMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			if len(path) == 1 {
				path = "index" // for root.
			} else if path == "/favicon.ico" {
				next.ServeHTTP(w, r)
				return
			} else {
				path = path[1:]
				path = strings.Replace(path, "/", ".", -1)
			}

			statsD.Increment(fmt.Sprintf("%s.request", path))

			newResponseWriter := &statusCodeReporter{ResponseWriter: w, statusCode: http.StatusOK}

			stop := statsD.Record(fmt.Sprintf("%s.time", path), 1)
			next.ServeHTTP(newResponseWriter, r)
			stop()

			statsD.Increment(fmt.Sprintf("%s.response.%d", path, newResponseWriter.statusCode))
		})
	}

	mux := http.DefaultServeMux

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello world from index")
	})

	mux.HandleFunc("/other", func(w http.ResponseWriter, r *http.Request) {
		statsD.Unique("other.unique", 1)
		fmt.Fprintln(w, "Hello from other page")
	})

	http.ListenAndServe(":80", statsDMiddleware(mux))
}
