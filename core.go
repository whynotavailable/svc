package svc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Middleware = func(r *http.Request) error

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
