package debug

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	goahttp "goa.design/goa/v3/http"
)

func TestAdapt(t *testing.T) {
	mux := goahttp.NewMuxer()
	adapted := Adapt(mux)
	adapted.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	adapted.Handle("/foo", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("FOO"))
	}))
	req := &http.Request{Method: "GET", URL: &url.URL{Path: "/"}}
	w := httptest.NewRecorder()
	adapted.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("got status %d, expected %d", w.Code, http.StatusOK)
	}
	if w.Body.String() != "OK" {
		t.Errorf("got body %q, expected %q", w.Body.String(), "OK")
	}
	req = &http.Request{Method: "GET", URL: &url.URL{Path: "/foo"}}
	w = httptest.NewRecorder()
	adapted.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("got status %d, expected %d", w.Code, http.StatusOK)
	}
	if w.Body.String() != "FOO" {
		t.Errorf("got body %q, expected %q", w.Body.String(), "FOO")
	}
}
