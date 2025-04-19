package svc

import (
	"fmt"
	"log/slog"
	"net/http"
	"reflect"
)

// RpcFunction contains building blocks for documentation for functions.
type RpcFunction struct {
	bodyType reflect.Type
	meta     map[string]string
}

// BodyType allows you to set the type of a function body.
// Accepts an object to reflect.
func (f *RpcFunction) BodyType(obj any) *RpcFunction {
	f.bodyType = reflect.TypeOf(obj)

	return f
}

// Meta sets a key and a value for functions. This is included in the documentation endpoint.
func (f *RpcFunction) Meta(key string, value string) *RpcFunction {
	f.meta[key] = value

	return f
}

// AddFunction adds a new function to the container. Returns a pointer to the function for chaining.
func (container *RpcContainer) AddFunction(key string, handler HandlerFunc) *RpcFunction {
	function := RpcFunction{
		meta: map[string]string{},
	}
	container.functions[key] = function

	container.mux.HandleFunc(fmt.Sprintf("POST /%s", key), handler)

	return &function
}

// RpcContainer is the handler for RPC style functions.
type RpcContainer struct {
	functions  map[string]RpcFunction
	middlewars []Middleware
	mux        http.ServeMux
}

func (c *RpcContainer) AddMiddleware(f Middleware) {
	c.middlewars = append(c.middlewars, f)
}

func NewRpcContainer() *RpcContainer {
	return &RpcContainer{
		functions:  map[string]RpcFunction{},
		middlewars: []Middleware{},
		mux:        http.ServeMux{},
	}
}

func (container *RpcContainer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := ExecuteMiddleware(container.middlewars, r)
	if err != nil {
		slog.Error("middleware error", slog.String("err", err.Error()))
		WriteError(w, err)
		return
	}

	container.mux.ServeHTTP(w, r)
}
