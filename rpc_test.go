package svc_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/whynotavailable/svc"
	"github.com/whynotavailable/svc/asserts"
)

var example = svc.RpcFunction{
	Key: "hi",
	Function: func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	},
}

func TestRpcContainer(t *testing.T) {
	rpc := svc.NewRpcContainer()

	rpc.Add(example)

	{
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("POST", "/hi", nil)

		rpc.ServeHTTP(recorder, request)

		asserts.Eq(t, recorder.Code, http.StatusOK)

		data, err := io.ReadAll(recorder.Body)

		asserts.NoError(t, err)

		asserts.Eq(t, "hi", string(data))
	}

	{
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("POST", "/bye", nil)

		rpc.ServeHTTP(recorder, request)

		asserts.Eq(t, recorder.Code, http.StatusNotFound)
	}
}

type additionInput struct {
	A int
	B int
}

type simpleResponse[T any] struct {
	Value T `json:"value"`
}

var addition = svc.RpcFunction{
	Key: "addition",
	Docs: svc.NewFunctionInfo(
		additionInput{},
		simpleResponse[int]{},
		nil,
	),
	Function: func(w http.ResponseWriter, r *http.Request) {
		body, err := svc.ReadJson[additionInput](r)
		if err != nil {
			svc.WriteErrorBadRequest(w)
		}

		svc.WriteJson(w, simpleResponse[int]{
			Value: body.A + body.B,
		})
	},
}

func ExampleRpcContainer() {
	rpc := svc.NewRpcContainer()

	rpc.Add(addition)

	err := http.ListenAndServe("127.0.0.1:3333", rpc)
	if err != nil {
		fmt.Println(err)
	}
}
