// Code generated by goa v3.9.1, DO NOT EDIT.
//
// test HTTP server encoders and decoders
//
// Command:
// $ goa gen goa.design/clue/internal/testsvc/design

package server

import (
	"context"
	"io"
	"net/http"
	"strconv"

	test "goa.design/clue/internal/testsvc/gen/test"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// EncodeHTTPMethodResponse returns an encoder for responses returned by the
// test http_method endpoint.
func EncodeHTTPMethodResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res, _ := v.(*test.Fields)
		enc := encoder(ctx, w)
		body := NewHTTPMethodResponseBody(res)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeHTTPMethodRequest returns a decoder for requests sent to the test
// http_method endpoint.
func DecodeHTTPMethodRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			body HTTPMethodRequestBody
			err  error
		)
		err = decoder(r).Decode(&body)
		if err != nil {
			if err == io.EOF {
				return nil, goa.MissingPayloadError()
			}
			return nil, goa.DecodePayloadError(err.Error())
		}

		var (
			i int

			params = mux.Vars(r)
		)
		{
			iRaw := params["i"]
			v, err2 := strconv.ParseInt(iRaw, 10, strconv.IntSize)
			if err2 != nil {
				err = goa.MergeErrors(err, goa.InvalidFieldTypeError("i", iRaw, "integer"))
			}
			i = int(v)
		}
		if err != nil {
			return nil, err
		}
		payload := NewHTTPMethodFields(&body, i)

		return payload, nil
	}
}
