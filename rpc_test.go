package svc_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/whynotavailable/svc"
	"github.com/whynotavailable/svc/asserts"
)

func TestRpcContainer(t *testing.T) {
	rpc := svc.NewRpcContainer()

	rpc.AddFunction(svc.RpcFunction{
		Key: "hi",
		Function: func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("hi"))
		},
	})

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

func ExampleRpcContainer() {
	rpc := svc.NewRpcContainer()

	rpc.AddFunction(svc.RpcFunction{
		Key: "hi",
		Function: func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("hi"))
		},
	})
}
