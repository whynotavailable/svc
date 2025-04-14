package svc

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

type RcpHandler = func(w http.ResponseWriter, r *http.Request)

type RpcFunction struct {
	Handler  RcpHandler
	bodyType reflect.Type
	meta     map[string]string
}

func (f *RpcFunction) BodyType(obj any) *RpcFunction {
	f.bodyType = reflect.TypeOf(obj)

	return f
}

func (f *RpcFunction) Meta(key string, value string) *RpcFunction {
	f.meta[key] = value

	return f
}

func (container *RpcContainer) AddFunction(key string, handler RcpHandler) *RpcFunction {
	function := RpcFunction{
		Handler: handler,
		meta:    map[string]string{},
	}
	container.functions[key] = &function

	return &function
}

type RpcContainer struct {
	functions  map[string]*RpcFunction
	docs       map[string]FunctionDoc
	middlewars []Middleware
}

func NewRpcContainer() RpcContainer {
	return RpcContainer{
		functions:  map[string]*RpcFunction{},
		middlewars: []Middleware{},
	}
}

func (container *RpcContainer) SetupMux(mux *http.ServeMux, prefix string) error {
	mux.Handle(fmt.Sprintf("%s/", prefix), http.StripPrefix(prefix, container))
	return nil
}

func (container *RpcContainer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if container.docs != nil && r.URL.Path == "/_info" {
			WriteJson(w, container.docs)
			return
		} else {
			WriteErrorNotFound(w)
			return
		}
	}

	if r.Method != http.MethodPost {
		WriteErrorBadRequest(w)
		return
	}

	functionKey := strings.TrimLeft(r.URL.Path, "/")

	f, ok := container.functions[functionKey]

	if !ok {
		WriteErrorNotFound(w)
		return
	}

	f.Handler(w, r)
}
