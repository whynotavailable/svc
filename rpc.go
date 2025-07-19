package svc

import (
	"fmt"
	"net/http"

	"github.com/invopop/jsonschema"
)

// RpcFunction contains building blocks for documentation for functions.
type RpcFunction struct {
	Key string
	// Input is an object representing the input type
	Input any
	// Output is an object representing the output type
	Output   any
	Meta     map[string]string
	Function func(w http.ResponseWriter, r *http.Request)
}

type GenericFunction[T any] struct{}

type RpcFunctionInfo struct {
	InputSchema  *jsonschema.Schema `json:"input"`
	OutputSchema *jsonschema.Schema `json:"output"`
	Meta         map[string]string  `json:"meta"`
}

func (f *RpcFunction) Info() RpcFunctionInfo {
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

// AddFunction adds a new function to the container. Returns a pointer to the function for chaining.
func (container *RpcContainer) AddFunction(function RpcFunction) *RpcFunction {
	container.functions[function.Key] = function.Info()

	container.mux.HandleFunc(fmt.Sprintf("POST /%s", function.Key), function.Function)

	return &function
}

func (container *RpcContainer) SetupDocs() {
	container.mux.HandleFunc("GET /_info", func(w http.ResponseWriter, r *http.Request) {
		WriteJson(w, container.functions)
	})
}

// RpcContainer is the handler for RPC style functions.
type RpcContainer struct {
	functions map[string]RpcFunctionInfo
	mux       http.ServeMux
}

func NewRpcContainer() *RpcContainer {
	return &RpcContainer{
		functions: map[string]RpcFunctionInfo{},
		mux:       http.ServeMux{},
	}
}

func (container *RpcContainer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	container.mux.ServeHTTP(w, r)
}
