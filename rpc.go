package svc

import (
	"fmt"
	"net/http"

	"github.com/invopop/jsonschema"
)

// RpcFunction contains building blocks for documentation for functions.
type RpcFunction struct {
	Name         string
	InputObject  any
	OutputObject any
	Meta         map[string]string
	Function     HandlerFunc
}

type RpcFunctionInfo struct {
	Name         string
	InputSchema  *jsonschema.Schema
	OutputSchema *jsonschema.Schema
	Meta         map[string]string
}

// AddFunction adds a new function to the container. Returns a pointer to the function for chaining.
func (container *RpcContainer) AddFunction(function RpcFunction) *RpcFunction {
	info := RpcFunctionInfo{
		Name: function.Name,
		Meta: function.Meta,
	}

	if function.InputObject != nil {
		info.InputSchema = jsonschema.Reflect(function.InputObject)
	}

	if function.OutputObject != nil {
		info.OutputSchema = jsonschema.Reflect(function.OutputObject)
	}

	container.functions[function.Name] = info

	container.mux.HandleFunc(fmt.Sprintf("POST /%s", function.Name), function.Function)

	return &function
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
