// Code generated by goa v3.9.1, DO NOT EDIT.
//
// itest gRPC client CLI support package
//
// Command:
// $ goa gen goa.design/clue/internal/testsvc/design

package cli

import (
	"flag"
	"fmt"
	"os"

	testc "goa.design/clue/internal/testsvc/gen/grpc/test/client"
	goa "goa.design/goa/v3/pkg"
	grpc "google.golang.org/grpc"
)

// UsageCommands returns the set of commands and sub-commands using the format
//
//	command (subcommand1|subcommand2|...)
func UsageCommands() string {
	return `test (grpc-method|grpc-stream)
`
}

// UsageExamples produces an example of a valid invocation of the CLI tool.
func UsageExamples() string {
	return os.Args[0] + ` test grpc-method --message '{
      "i": 4434351785751264939,
      "s": "Corrupti repellat autem sit architecto ut."
   }'` + "\n" +
		""
}

// ParseEndpoint returns the endpoint and payload as specified on the command
// line.
func ParseEndpoint(cc *grpc.ClientConn, opts ...grpc.CallOption) (goa.Endpoint, interface{}, error) {
	var (
		testFlags = flag.NewFlagSet("test", flag.ContinueOnError)

		testGrpcMethodFlags       = flag.NewFlagSet("grpc-method", flag.ExitOnError)
		testGrpcMethodMessageFlag = testGrpcMethodFlags.String("message", "", "")

		testGrpcStreamFlags = flag.NewFlagSet("grpc-stream", flag.ExitOnError)
	)
	testFlags.Usage = testUsage
	testGrpcMethodFlags.Usage = testGrpcMethodUsage
	testGrpcStreamFlags.Usage = testGrpcStreamUsage

	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		return nil, nil, err
	}

	if flag.NArg() < 2 { // two non flag args are required: SERVICE and ENDPOINT (aka COMMAND)
		return nil, nil, fmt.Errorf("not enough arguments")
	}

	var (
		svcn string
		svcf *flag.FlagSet
	)
	{
		svcn = flag.Arg(0)
		switch svcn {
		case "test":
			svcf = testFlags
		default:
			return nil, nil, fmt.Errorf("unknown service %q", svcn)
		}
	}
	if err := svcf.Parse(flag.Args()[1:]); err != nil {
		return nil, nil, err
	}

	var (
		epn string
		epf *flag.FlagSet
	)
	{
		epn = svcf.Arg(0)
		switch svcn {
		case "test":
			switch epn {
			case "grpc-method":
				epf = testGrpcMethodFlags

			case "grpc-stream":
				epf = testGrpcStreamFlags

			}

		}
	}
	if epf == nil {
		return nil, nil, fmt.Errorf("unknown %q endpoint %q", svcn, epn)
	}

	// Parse endpoint flags if any
	if svcf.NArg() > 1 {
		if err := epf.Parse(svcf.Args()[1:]); err != nil {
			return nil, nil, err
		}
	}

	var (
		data     interface{}
		endpoint goa.Endpoint
		err      error
	)
	{
		switch svcn {
		case "test":
			c := testc.NewClient(cc, opts...)
			switch epn {
			case "grpc-method":
				endpoint = c.GrpcMethod()
				data, err = testc.BuildGrpcMethodPayload(*testGrpcMethodMessageFlag)
			case "grpc-stream":
				endpoint = c.GrpcStream()
				data = nil
			}
		}
	}
	if err != nil {
		return nil, nil, err
	}

	return endpoint, data, nil
}

// testUsage displays the usage of the test command and its subcommands.
func testUsage() {
	fmt.Fprintf(os.Stderr, `Service is the test service interface.
Usage:
    %[1]s [globalflags] test COMMAND [flags]

COMMAND:
    grpc-method: GrpcMethod implements grpc_method.
    grpc-stream: GrpcStream implements grpc_stream.

Additional help:
    %[1]s test COMMAND --help
`, os.Args[0])
}
func testGrpcMethodUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] test grpc-method -message JSON

GrpcMethod implements grpc_method.
    -message JSON: 

Example:
    %[1]s test grpc-method --message '{
      "i": 4434351785751264939,
      "s": "Corrupti repellat autem sit architecto ut."
   }'
`, os.Args[0])
}

func testGrpcStreamUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] test grpc-stream

GrpcStream implements grpc_stream.

Example:
    %[1]s test grpc-stream
`, os.Args[0])
}
