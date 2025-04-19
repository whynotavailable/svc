package svc_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/whynotavailable/svc"
	"github.com/whynotavailable/svc/asserts"
)

func TestBasicsForMiddleware(t *testing.T) {
	container := svc.NewHttpContainer()

	container.AddMiddleware(svc.LoggingMiddleware)

	container.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		svc.WriteText(w, "ok")
	})

	rr := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/hi", nil)

	container.ServeHTTP(rr, r)
	asserts.StatusEq(t, http.StatusOK, rr.Code)
}

func TestFailureInMiddleware(t *testing.T) {
	container := svc.NewHttpContainer()

	container.AddMiddleware(func(*http.Request) error {
		return errors.New("what")
	})

	container.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		svc.WriteText(w, "ok")
	})

	rr := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/hi", nil)

	container.ServeHTTP(rr, r)
	asserts.StatusEq(t, http.StatusInternalServerError, rr.Code)
}

func ExampleHttpContainer() {
	container := svc.NewHttpContainer()

	container.AddMiddleware(svc.LoggingMiddleware)

	container.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		svc.WriteText(w, "ok")
	})

	http.ListenAndServe("0.0.0.0:3456", container)
}
