package svc_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/whynotavailable/svc"
	"github.com/whynotavailable/svc/asserts"
)

func TestProxySimple(t *testing.T) {
	svc.ProxyInit()

	counter := 0

	tripper := asserts.NewRoundTripper(func(r *http.Request) {
		asserts.Eq(t, r.URL.String(), "https://google.com/hi")
		counter++
	})

	proxy := svc.ProxyContainer{
		Target: "https://google.com",
		Client: &http.Client{
			Transport: tripper,
		},
	}

	rr := httptest.ResponseRecorder{}
	r := httptest.NewRequest("GET", "/proxy/hi", nil)

	router := http.NewServeMux()

	svc.SetupContainer(router, "/proxy", &proxy)
	router.ServeHTTP(&rr, r)
	asserts.Eq(t, http.StatusOK, rr.Code, "proxy status")

	asserts.Eq(t, 1, counter, "Counter")
}

func ExampleProxyContainer() {
	container := svc.NewProxyContainer("https://target.com")

	svc.SetupContainer(http.DefaultServeMux, "/proxy", container)

	http.ListenAndServe("0.0.0.0:3456", nil)
}
