package debug

import (
	"net/http"

	goahttp "goa.design/goa/v3/http"
)

// muxAdapter is a debug.Muxer adapter for goahttp.Muxer.
type muxAdapter struct {
	muxer goahttp.Muxer
}

// HTTP methods supported by the adapter.
var httpMethods = []string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "CONNECT", "OPTIONS", "TRACE"}

// Adapt returns a debug.Muxer adapter for the given goahttp.Muxer.
func Adapt(m goahttp.Muxer) Muxer {
	return muxAdapter{muxer: m}
}

func (m muxAdapter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.muxer.ServeHTTP(w, r)
}

func (m muxAdapter) Handle(path string, handler http.Handler) {
	for _, method := range httpMethods {
		m.muxer.Handle(method, path, handler.ServeHTTP)
	}
}

func (m muxAdapter) HandleFunc(path string, handler func(http.ResponseWriter, *http.Request)) {
	for _, method := range httpMethods {
		m.muxer.Handle(method, path, handler)
	}
}
