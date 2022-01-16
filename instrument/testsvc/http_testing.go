package testsvc

import (
	"context"
	"net"
	"net/http"
	"testing"

	goahttp "goa.design/goa/v3/http"

	"github.com/crossnokaye/micro/instrument/testsvc/gen/http/test/client"
	"github.com/crossnokaye/micro/instrument/testsvc/gen/http/test/server"
	"github.com/crossnokaye/micro/instrument/testsvc/gen/test"
)

type (
	// HTTPClient is a test service HTTP client.
	HTTPClient interface {
		HTTPMethod(ctx context.Context, req *Fields) (res *Fields, err error)
	}

	// HTTPOption is a function that can be used to configure the HTTP server.
	HTTPOption func(*httpOptions)

	httpc struct {
		genc *client.Client
	}

	httpOptions struct {
		fn         UnaryFunc
		middleware func(http.Handler) http.Handler
	}
)

// SetupHTTP instantiates the test gRPC service with the given options.
// It returns a ready-to-use client and a function used to stop the server.
func SetupHTTP(t *testing.T, opts ...HTTPOption) (c HTTPClient, stop func()) {
	t.Helper()

	// Configure options
	var options httpOptions
	for _, opt := range opts {
		opt(&options)
	}

	// Create test HTTP server
	svc := &svc{httpfn: options.fn}
	endpoints := test.NewEndpoints(svc)
	mux := goahttp.NewMuxer()
	svr := server.New(endpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, nil, nil)
	server.Mount(mux, svr)
	var handler http.Handler = mux
	if options.middleware != nil {
		handler = options.middleware(handler)
	}

	// Start listen loop
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	go http.Serve(l, handler)

	// Connect to in-memory server
	c = httpc{client.NewClient("http", l.Addr().String(), http.DefaultClient, goahttp.RequestEncoder, goahttp.ResponseDecoder, false)}

	// Cleanup
	stop = func() {
		l.Close()
	}

	return
}

func WithHTTPFunc(fn UnaryFunc) HTTPOption {
	return func(opt *httpOptions) {
		opt.fn = fn
	}
}

func WithHTTPMiddleware(fn func(http.Handler) http.Handler) HTTPOption {
	return func(opt *httpOptions) {
		opt.middleware = fn
	}
}

// HTTPMethod implements HTTPClient.HTTPMethod using the generated client.
func (c httpc) HTTPMethod(ctx context.Context, req *Fields) (res *Fields, err error) {
	rq := &test.Fields{}
	if req != nil {
		*rq = test.Fields(*req)
	}
	resp, err := c.genc.HTTPMethod()(ctx, rq)
	if err != nil {
		return nil, err
	}
	if resp != nil {
		res = &Fields{}
		*res = Fields(*(resp.(*test.Fields)))
	}
	return
}
