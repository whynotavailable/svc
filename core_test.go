package svc_test

import (
	"net/http"

	"github.com/whynotavailable/svc"
)

func ExampleLoggingMiddleware() {
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		svc.WriteText(w, "ok")
	})

	var handler http.Handler = svc.NewLoggingMiddleware(http.DefaultServeMux)

	http.ListenAndServe("0.0.0.0:3456", handler)
}

func ExampleGenericMiddleware() {
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		svc.WriteText(w, "ok")
	})

	var handler http.Handler = http.DefaultServeMux

	handler = svc.NewMiddleware(handler, func(r *http.Request) {
		// After Logging
	})

	handler = svc.NewLoggingMiddleware(handler)

	handler = svc.NewMiddleware(handler, func(r *http.Request) {
		// Before Logging
		r.Header["Test"] = []string{"Val"}
	})

	http.ListenAndServe("0.0.0.0:3456", handler)
}
