package svc_test

import (
	"net/http"

	"github.com/whynotavailable/svc"
)

func ExampleHttpContainer() {
	container := svc.NewHttpContainer()
	container.Mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		svc.WriteText(w, "ok")
	})

	http.ListenAndServe("0.0.0.0:3456", container)
}
