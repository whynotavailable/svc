package svc_test

import (
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
	svc.SetupMux(testMux, "/rpc", &rpcContainer)

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

	svc.SetupMux(http.DefaultServeMux, "/rpc", &rpcContainer)

	http.ListenAndServe("0.0.0.0:3456", nil)
}
