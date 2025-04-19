package svc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type Middleware = func(r *http.Request) error

func ExecuteMiddleware(middlewares []Middleware, r *http.Request) error {
	for _, middleware := range middlewares {
		err := middleware(r)
		if err != nil {
			return err
		}
	}

	return nil
}

// HandlerFunc is exactly the same as a standard function handler, just in type alias form
type HandlerFunc = func(w http.ResponseWriter, r *http.Request)

func ReadJson[T any](r *http.Request) (T, error) {
	var obj T

	if r.Body == nil {
		return obj, errors.New("tried to parse an empty body")
	}

	rawData, err := io.ReadAll(r.Body)
	r.Body = nil

	if err != nil {
		return obj, err
	}

	err = json.Unmarshal(rawData, &obj)
	if err != nil {
		return obj, err
	}

	return obj, nil
}

func WriteJson(w http.ResponseWriter, obj any) error {
	rawData, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(rawData)

	return nil
}

func WriteText(w http.ResponseWriter, text string) error {
	w.Header().Add("Content-Type", "plain/text")
	w.WriteHeader(200)
	w.Write([]byte(text))

	return nil
}

func WriteError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintln(w, err)
}

func WriteErrorWithCode(w http.ResponseWriter, err error, code int) {
	w.WriteHeader(code)
	fmt.Fprintln(w, err)
}

func WriteErrorNotFound(w http.ResponseWriter) {
	WriteErrorWithCode(w, errors.New("not found"), http.StatusNotFound)
}

func WriteErrorBadRequest(w http.ResponseWriter) {
	WriteErrorWithCode(w, errors.New("bad request"), http.StatusBadRequest)
}

type SimpleMessage struct {
	Message string `json:"message"`
}

func SetupContainer(mux *http.ServeMux, prefix string, handler http.Handler) {
	mux.Handle(fmt.Sprintf("%s/", prefix), http.StripPrefix(prefix, handler))
}

func LoggingMiddleware(r *http.Request) error {
	slog.Info("Request", slog.String("method", r.Method), slog.String("path", r.URL.String()))
	return nil
}
