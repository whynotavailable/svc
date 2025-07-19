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

	Meta map[string]string
}

type RpcFunctionInfo struct {
	InputSchema  *jsonschema.Schema `json:"input"`
	OutputSchema *jsonschema.Schema `json:"output"`
	Meta         map[string]string  `json:"meta"`
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

	if info.Meta == nil {
		info.Meta = map[string]string{}
	}

	return info
}

func (container *RpcContainer) Add(key string, docs *RpcFunctionDocs, function func(w http.ResponseWriter, r *http.Request)) {
	if docs != nil {
		container.functions[key] = docs.Info()
	} else {
		container.functions[key] = nil
	}

	container.mux.HandleFunc(fmt.Sprintf("POST /%s", key), function)
}

func (container *RpcContainer) SetupDocs() {
	container.mux.HandleFunc("GET /_info", func(w http.ResponseWriter, r *http.Request) {
		WriteJson(w, container.functions)
	})
}

// RpcContainer is the handler for RPC style functions.
type RpcContainer struct {
	functions map[string]any
	mux       http.ServeMux
}

func NewRpcContainer() *RpcContainer {
	return &RpcContainer{
		functions: map[string]any{},
		mux:       http.ServeMux{},
	}
}

func (container *RpcContainer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	container.mux.ServeHTTP(w, r)
}
