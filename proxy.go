package svc

import (
	"fmt"
	"io"
	"net/http"
)

type ProxyContainer struct {
	Client     http.Client
	Target     string
	middlewars []Middleware
}

var hopByHopHeaders map[string]bool = nil

func (container *ProxyContainer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s%s", container.Target, r.URL.String())
	req, err := http.NewRequest(r.Method, url, r.Body)
	if err != nil {
		WriteError(w, err)
		return
	}

	copyHeaders(r.Header, req.Header)

	resp, err := container.Client.Do(req)
	if err != nil {
		WriteError(w, err)
		return
	}

	copyHeaders(resp.Header, w.Header())

	w.WriteHeader(resp.StatusCode)

	if resp.Body != nil {
		defer resp.Body.Close()

		io.Copy(w, resp.Body)
	}
}

func copyHeaders(src http.Header, dst http.Header) {
	for key, header := range src {
		if hopByHopHeaders[key] {
			continue
		}

		for _, val := range header {
			dst.Add(key, val)
		}
	}
}

// ProxyInit should be called once if you are using proxies.
// It won't hurt to call it more though
func ProxyInit() {
	if hopByHopHeaders == nil {
		initHopByHopHeaders()
	}
}

func initHopByHopHeaders() {
	hopByHopHeaders = map[string]bool{
		"Connection":          true,
		"Keep-Alive":          true,
		"Proxy-Authenticate":  true,
		"Proxy-Authorization": true,
		"Te":                  true,
		"Trailers":            true,
		"Transfer-Encoding":   true,
		"Upgrade":             true,
	}
}
