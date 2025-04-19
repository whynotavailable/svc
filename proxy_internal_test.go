package svc

import (
	"net/http"
	"testing"

	"github.com/whynotavailable/svc/asserts"
)

func TestCopyHeaders(t *testing.T) {
	initHopByHopHeaders()

	var map1 http.Header = map[string][]string{
		"Auth": {"hi"},
		"Te":   {"any"},
	}

	var map2 http.Header = map[string][]string{}

	copyHeaders(map1, map2)

	asserts.Eq(t, map2["Auth"][0], "hi")

	asserts.NoKey(t, map2, "Te")
}
