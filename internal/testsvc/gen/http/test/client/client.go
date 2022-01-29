// Code generated by goa v3.5.4, DO NOT EDIT.
//
// test client HTTP transport
//
// Command:
// $ goa gen goa.design/clue/internal/testsvc/design

package client

import (
	"context"
	"net/http"

	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// Client lists the test service endpoint HTTP clients.
type Client struct {
	// HTTPMethod Doer is the HTTP client used to make requests to the http_method
	// endpoint.
	HTTPMethodDoer goahttp.Doer

	// RestoreResponseBody controls whether the response bodies are reset after
	// decoding so they can be read again.
	RestoreResponseBody bool

	scheme  string
	host    string
	encoder func(*http.Request) goahttp.Encoder
	decoder func(*http.Response) goahttp.Decoder
}

// NewClient instantiates HTTP clients for all the test service servers.
func NewClient(
	scheme string,
	host string,
	doer goahttp.Doer,
	enc func(*http.Request) goahttp.Encoder,
	dec func(*http.Response) goahttp.Decoder,
	restoreBody bool,
) *Client {
	return &Client{
		HTTPMethodDoer:      doer,
		RestoreResponseBody: restoreBody,
		scheme:              scheme,
		host:                host,
		decoder:             dec,
		encoder:             enc,
	}
}

// HTTPMethod returns an endpoint that makes HTTP requests to the test service
// http_method server.
func (c *Client) HTTPMethod() goa.Endpoint {
	var (
		encodeRequest  = EncodeHTTPMethodRequest(c.encoder)
		decodeResponse = DecodeHTTPMethodResponse(c.decoder, c.RestoreResponseBody)
	)
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		req, err := c.BuildHTTPMethodRequest(ctx, v)
		if err != nil {
			return nil, err
		}
		err = encodeRequest(req, v)
		if err != nil {
			return nil, err
		}
		resp, err := c.HTTPMethodDoer.Do(req)
		if err != nil {
			return nil, goahttp.ErrRequestError("test", "http_method", err)
		}
		return decodeResponse(resp)
	}
}
