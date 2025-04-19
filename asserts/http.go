package asserts

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Pulled stuff from google and adapted it.

type MockRoundTripper struct {
	Resp    *http.Response
	Err     error
	Checker func(*http.Request)
}

func NewRoundTripper(checker func(*http.Request)) *MockRoundTripper {
	return &MockRoundTripper{
		Checker: checker,
		Resp: &http.Response{
			StatusCode: 200,
			Body:       httptest.NewRecorder().Result().Body,
		},
	}
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	m.Checker(req)
	return m.Resp, m.Err
}

func StatusEq(t *testing.T, target int, actual int, extras ...any) {
	if target != actual {
		msg := fmt.Sprintf("Incorrect status code, expected %d got %d", target, actual)
		throw(t, msg, extras...)
	}
}
