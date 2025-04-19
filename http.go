package svc

import (
	"log/slog"
	"net/http"
)

// Real simple http container that acts as a regular mux but with middleware support

type HttpContainer struct {
	middlewars []Middleware
	Mux        http.ServeMux
}

func (c *HttpContainer) AddMiddleware(f Middleware) {
	c.middlewars = append(c.middlewars, f)
}

func NewHttpContainer() *HttpContainer {
	return &HttpContainer{
		middlewars: []Middleware{},
		Mux:        http.ServeMux{},
	}
}

func (container *HttpContainer) Handle(pattern string, handler http.Handler) {
	container.Mux.Handle(pattern, handler)
}

func (container *HttpContainer) HandleFunc(pattern string, handler HandlerFunc) {
	container.Mux.HandleFunc(pattern, handler)
}

func (container *HttpContainer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := ExecuteMiddleware(container.middlewars, r)
	if err != nil {
		slog.Error("middleware error", slog.String("err", err.Error()))
		WriteError(w, err)
		return
	}

	container.Mux.ServeHTTP(w, r)
}
