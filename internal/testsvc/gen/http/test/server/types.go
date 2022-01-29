// Code generated by goa v3.5.4, DO NOT EDIT.
//
// test HTTP server types
//
// Command:
// $ goa gen goa.design/clue/internal/testsvc/design

package server

import (
	test "goa.design/clue/internal/testsvc/gen/test"
)

// HTTPMethodRequestBody is the type of the "test" service "http_method"
// endpoint HTTP request body.
type HTTPMethodRequestBody struct {
	// String operand
	S *string `form:"s,omitempty" json:"s,omitempty" xml:"s,omitempty"`
}

// HTTPMethodResponseBody is the type of the "test" service "http_method"
// endpoint HTTP response body.
type HTTPMethodResponseBody struct {
	// String operand
	S *string `form:"s,omitempty" json:"s,omitempty" xml:"s,omitempty"`
	// Int operand
	I *int `form:"i,omitempty" json:"i,omitempty" xml:"i,omitempty"`
}

// NewHTTPMethodResponseBody builds the HTTP response body from the result of
// the "http_method" endpoint of the "test" service.
func NewHTTPMethodResponseBody(res *test.Fields) *HTTPMethodResponseBody {
	body := &HTTPMethodResponseBody{
		S: res.S,
		I: res.I,
	}
	return body
}

// NewHTTPMethodFields builds a test service http_method endpoint payload.
func NewHTTPMethodFields(body *HTTPMethodRequestBody, i int) *test.Fields {
	v := &test.Fields{
		S: body.S,
	}
	v.I = &i

	return v
}
