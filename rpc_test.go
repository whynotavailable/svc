package svc_test

import (
	"net/http"

	"github.com/whynotavailable/svc"
)

func ExampleRpcContainer_AddFunction() {
	rpcContainer := svc.NewRpcContainer()

	rpcContainer.AddFunction("hello", func(w http.ResponseWriter, r *http.Request) {
		svc.WriteJson(w, svc.SimpleMessage{
			Message: "dave",
		})
	})

	rpcContainer.AddFunction("hello-body", func(w http.ResponseWriter, r *http.Request) {
		body, err := svc.ReadJson[svc.SimpleMessage](r)
		if err != nil {
			svc.WriteErrorWithCode(w, err, http.StatusBadRequest)
			return
		}

		svc.WriteJson(w, svc.SimpleMessage{
			Message: body.Message,
		})
	}).BodyType(svc.SimpleMessage{}) // Only needed if using doc generation

	rpcContainer.GenerateDocs()
	rpcContainer.SetupMux(http.DefaultServeMux, "/rpc")

	http.ListenAndServe("0.0.0.0:3456", nil)
}
