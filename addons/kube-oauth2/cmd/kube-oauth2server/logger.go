package main

import (
	"log"
	"net/http"
	"time"
)

type trackedWriter struct {
	http.ResponseWriter
	status int
}

func (w *trackedWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *trackedWriter) Write(data []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	return w.ResponseWriter.Write(data)
}

type requestLogger struct {
	next http.Handler
}

func (l requestLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t1 := time.Now()
	tw := &trackedWriter{ResponseWriter: w}
	l.next.ServeHTTP(tw, r)
	log.Printf("%s %d %s   %v", r.Method, tw.status, r.URL.String(), time.Since(t1))
}
