package svc_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/whynotavailable/svc"
)

func Assert(t *testing.T, cond bool, args ...any) {
	if !cond {
		t.Error(args...)
		t.FailNow()
	}
}

func NoError(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func CallTest(mux *http.ServeMux, r *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, r)

	return recorder
}

func TestSimple(t *testing.T) {
	testMux := http.NewServeMux()
	rpcContainer := svc.NewRpcContainer()

	rpcContainer.AddFunction("hello", func(w http.ResponseWriter, r *http.Request) {
		svc.WriteJson(w, svc.SimpleMessage{
			Message: "dave",
		})
	})

	rpcContainer.GenerateDocs()
	svc.SetupContainer(testMux, "/rpc", rpcContainer)

	{
		resp := CallTest(testMux, httptest.NewRequest("POST", "/rpc/hello", nil))

		Assert(t, resp.Code == http.StatusOK, "RPC call not ok", resp.Code)
		rawData, err := io.ReadAll(resp.Body)
		NoError(t, err)

		var data svc.SimpleMessage
		err = json.Unmarshal(rawData, &data)
		NoError(t, err)

		Assert(t, data.Message == "dave", "Message not dave", data.Message)
	}

	{
		resp := CallTest(testMux, httptest.NewRequest("GET", "/rpc/_info", nil))

		Assert(t, resp.Code == http.StatusOK, "docs call not ok", resp.Code)
	}
}

func TestBody(t *testing.T) {
	testMux := http.NewServeMux()
	rpcContainer := svc.NewRpcContainer()

	rpcContainer.AddFunction("hello", func(w http.ResponseWriter, r *http.Request) {
		body, err := svc.ReadJson[svc.SimpleMessage](r)
		if err != nil {
			svc.WriteErrorWithCode(w, err, http.StatusBadRequest)
			return
		}

		svc.WriteJson(w, svc.SimpleMessage{
			Message: body.Message,
		})
	}).BodyType(svc.SimpleMessage{})

	rpcContainer.GenerateDocs()
	svc.SetupContainer(testMux, "/rpc", rpcContainer)

	{
		// Good body test
		requestBody := svc.SimpleMessage{
			Message: "Cali",
		}
		requestBytes, err := json.Marshal(requestBody)
		NoError(t, err)

		resp := CallTest(testMux, httptest.NewRequest("POST", "/rpc/hello", bytes.NewBuffer(requestBytes)))

		Assert(t, resp.Code == http.StatusOK, "RPC call not ok", resp.Code)
		rawData, err := io.ReadAll(resp.Body)
		NoError(t, err)

		var data svc.SimpleMessage
		err = json.Unmarshal(rawData, &data)
		NoError(t, err)

		Assert(t, data.Message == requestBody.Message, "Message not correct", data.Message)
	}

	{
		// Bad body test
		requestBody := "string"
		requestBytes, err := json.Marshal(requestBody)
		NoError(t, err)

		resp := CallTest(testMux, httptest.NewRequest("POST", "/rpc/hello", bytes.NewBuffer(requestBytes)))

		Assert(t, resp.Code == http.StatusBadRequest, "RPC call not bad request", resp.Code)
	}
}

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

	svc.SetupContainer(http.DefaultServeMux, "/rpc", rpcContainer)

	http.ListenAndServe("0.0.0.0:3456", nil)
}
