package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = API("itest", func() {
	Description("instrument test service")
})

var _ = Service("test", func() {
	Method("http_method", func() {
		Payload(Fields)
		Result(Fields)
		HTTP(func() {
			POST("/{i}")
		})
	})

	Method("grpc_method", func() {
		Payload(Fields)
		Result(Fields)
		GRPC(func() {})
	})

	Method("grpc_stream", func() {
		StreamingPayload(Fields)
		StreamingResult(Fields)
		GRPC(func() {})
	})
})

var Fields = Type("Fields", func() {
	Field(1, "s", String, "String operand")
	Field(2, "i", Int, "Int operand")
})
