package svc

import (
	"fmt"
	"net/http"

	"github.com/invopop/jsonschema"
)

var Reflector = jsonschema.Reflector{
	ExpandedStruct: true,
}

type RpcFunctionInfo struct {
	InputSchema  *jsonschema.Schema `json:"input"`
	OutputSchema *jsonschema.Schema `json:"output"`
	Meta         any                `json:"meta"`
}

type RpcFunction struct {
	Key      string
	Docs     any
	Function func(w http.ResponseWriter, r *http.Request)
}

func NewFunctionInfo(input any, output any, meta any) RpcFunctionInfo {
	info := RpcFunctionInfo{
		Meta: meta,
	}

	if input != nil {
		info.InputSchema = Reflector.Reflect(input)
	}

	if output != nil {
		info.OutputSchema = Reflector.Reflect(output)
	}

	return info
}

func NewFunction(key string, docs any, function func(w http.ResponseWriter, r *http.Request)) RpcFunction {
	return RpcFunction{
		Key:      key,
		Docs:     docs,
		Function: function,
	}
}

func (container *RpcContainer) Add(function RpcFunction) {
	container.functions[function.Key] = function.Docs

	container.mux.HandleFunc(fmt.Sprintf("POST /%s", function.Key), function.Function)
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
