package svc

import (
	"fmt"
	"net/http"

	"github.com/invopop/jsonschema"
)

type RpcFunctionDocs struct {
	Input any
	// Output is an object representing the output type
	Output any

	Meta any
}

type RpcFunctionInfo struct {
	InputSchema  *jsonschema.Schema `json:"input"`
	OutputSchema *jsonschema.Schema `json:"output"`
	Meta         any                `json:"meta"`
}

func (f *RpcFunctionDocs) Info() RpcFunctionInfo {
	info := RpcFunctionInfo{
		Meta: f.Meta,
	}

	reflector := jsonschema.Reflector{
		ExpandedStruct: true,
	}

	if f.Input != nil {
		info.InputSchema = reflector.Reflect(f.Input)
	}

	if f.Output != nil {
		info.OutputSchema = reflector.Reflect(f.Output)
	}

	return info
}

func (container *RpcContainer) Add(key string, docs RpcFunctionDocs, function func(w http.ResponseWriter, r *http.Request)) {
	container.functions[key] = docs.Info()

	container.mux.HandleFunc(fmt.Sprintf("POST /%s", key), function)
}

// RpcContainer is the handler for RPC style functions.
type RpcContainer struct {
	NotFoundHandler http.HandlerFunc
	functions       map[string]any
	mux             http.ServeMux
}

func NewRpcContainer() *RpcContainer {
	return &RpcContainer{
		functions: map[string]any{},
		mux:       http.ServeMux{},
	}
}

func (container *RpcContainer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		WriteJson(w, container.functions)
		return
	}

	handler, pattern := container.mux.Handler(r)

	if pattern == "" && container.NotFoundHandler != nil {
		container.NotFoundHandler.ServeHTTP(w, r)
	} else {
		handler.ServeHTTP(w, r)
	}
}
