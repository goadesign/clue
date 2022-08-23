// Code generated by goa v3.8.3, DO NOT EDIT.
//
// IP Location API gRPC client CLI support package
//
// Command:
// $ goa gen goa.design/clue/example/weather/services/locator/design -o
// services/locator

package cli

import (
	"flag"
	"fmt"
	"os"

	locatorc "goa.design/clue/example/weather/services/locator/gen/grpc/locator/client"
	goa "goa.design/goa/v3/pkg"
	grpc "google.golang.org/grpc"
)

// UsageCommands returns the set of commands and sub-commands using the format
//
//	command (subcommand1|subcommand2|...)
func UsageCommands() string {
	return `locator get-location
`
}

// UsageExamples produces an example of a valid invocation of the CLI tool.
func UsageExamples() string {
	return os.Args[0] + ` locator get-location --message '{
      "field": "Totam delectus odio consequatur."
   }'` + "\n" +
		""
}

// ParseEndpoint returns the endpoint and payload as specified on the command
// line.
func ParseEndpoint(cc *grpc.ClientConn, opts ...grpc.CallOption) (goa.Endpoint, interface{}, error) {
	var (
		locatorFlags = flag.NewFlagSet("locator", flag.ContinueOnError)

		locatorGetLocationFlags       = flag.NewFlagSet("get-location", flag.ExitOnError)
		locatorGetLocationMessageFlag = locatorGetLocationFlags.String("message", "", "")
	)
	locatorFlags.Usage = locatorUsage
	locatorGetLocationFlags.Usage = locatorGetLocationUsage

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
		case "locator":
			svcf = locatorFlags
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
		case "locator":
			switch epn {
			case "get-location":
				epf = locatorGetLocationFlags

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
		case "locator":
			c := locatorc.NewClient(cc, opts...)
			switch epn {
			case "get-location":
				endpoint = c.GetLocation()
				data, err = locatorc.BuildGetLocationPayload(*locatorGetLocationMessageFlag)
			}
		}
	}
	if err != nil {
		return nil, nil, err
	}

	return endpoint, data, nil
}

// locatorUsage displays the usage of the locator command and its subcommands.
func locatorUsage() {
	fmt.Fprintf(os.Stderr, `Public HTTP frontend
Usage:
    %[1]s [globalflags] locator COMMAND [flags]

COMMAND:
    get-location: Retrieve location information for a given IP address

Additional help:
    %[1]s locator COMMAND --help
`, os.Args[0])
}
func locatorGetLocationUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] locator get-location -message JSON

Retrieve location information for a given IP address
    -message JSON: 

Example:
    %[1]s locator get-location --message '{
      "field": "Totam delectus odio consequatur."
   }'
`, os.Args[0])
}
