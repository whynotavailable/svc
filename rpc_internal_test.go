package svc

import (
	"net/http"
	"testing"

	"github.com/whynotavailable/svc/asserts"
)

func TestBodyType(t *testing.T) {
	container := NewRpcContainer()

	container.AddFunction("hi", func(w http.ResponseWriter, r *http.Request) {
		WriteText(w, "hi")
	}).BodyType(SchemaObject{})

	container.GenerateDocs()

	asserts.NotNil(t, container.functions["hi"].bodyType)
}
