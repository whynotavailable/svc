package svc_test

import (
	"net/http"

	"github.com/whynotavailable/svc"
)

func ExampleRpcContainer() {
	rpc := svc.NewRpcContainer()

	rpc.AddFunction(svc.RpcFunction{
		Name: "hi",
		Function: func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("hi"))
		},
	})
}
