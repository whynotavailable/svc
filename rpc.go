package svc

import (
	"fmt"
	"net/http"
	"reflect"
)

type RcpHandler = func(w http.ResponseWriter, r *http.Request)

type RpcFunction struct {
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
		meta: map[string]string{},
	}
	container.functions[key] = function

	container.mux.HandleFunc(fmt.Sprintf("POST /%s", key), handler)

	return &function
}

type RpcContainer struct {
	functions  map[string]RpcFunction
	middlewars []Middleware
	mux        http.ServeMux
}

func NewRpcContainer() *RpcContainer {
	return &RpcContainer{
		functions:  map[string]RpcFunction{},
		middlewars: []Middleware{},
		mux:        http.ServeMux{},
	}
}

func (container *RpcContainer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, middleware := range container.middlewars {
		err := middleware(r)
		if err != nil {
			WriteError(w, err)
			return
		}
	}

	container.mux.ServeHTTP(w, r)
}
