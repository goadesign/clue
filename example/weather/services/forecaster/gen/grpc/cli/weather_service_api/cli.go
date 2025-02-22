// Code generated by goa v3.20.0, DO NOT EDIT.
//
// Weather Service API gRPC client CLI support package
//
// Command:
// $ goa gen goa.design/clue/example/weather/services/forecaster/design -o
// services/forecaster

package cli

import (
	"flag"
	"fmt"
	"os"

	forecasterc "goa.design/clue/example/weather/services/forecaster/gen/grpc/forecaster/client"
	goa "goa.design/goa/v3/pkg"
	grpc "google.golang.org/grpc"
)

// UsageCommands returns the set of commands and sub-commands using the format
//
//	command (subcommand1|subcommand2|...)
func UsageCommands() string {
	return `forecaster forecast
`
}

// UsageExamples produces an example of a valid invocation of the CLI tool.
func UsageExamples() string {
	return os.Args[0] + ` forecaster forecast --message '{
      "lat": 37.8267,
      "long": -122.4233
   }'` + "\n" +
		""
}

// ParseEndpoint returns the endpoint and payload as specified on the command
// line.
func ParseEndpoint(
	cc *grpc.ClientConn,
	opts ...grpc.CallOption,
) (goa.Endpoint, any, error) {
	var (
		forecasterFlags = flag.NewFlagSet("forecaster", flag.ContinueOnError)

		forecasterForecastFlags       = flag.NewFlagSet("forecast", flag.ExitOnError)
		forecasterForecastMessageFlag = forecasterForecastFlags.String("message", "", "")
	)
	forecasterFlags.Usage = forecasterUsage
	forecasterForecastFlags.Usage = forecasterForecastUsage

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
		case "forecaster":
			svcf = forecasterFlags
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
		case "forecaster":
			switch epn {
			case "forecast":
				epf = forecasterForecastFlags

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
		data     any
		endpoint goa.Endpoint
		err      error
	)
	{
		switch svcn {
		case "forecaster":
			c := forecasterc.NewClient(cc, opts...)
			switch epn {

			case "forecast":
				endpoint = c.Forecast()
				data, err = forecasterc.BuildForecastPayload(*forecasterForecastMessageFlag)
			}
		}
	}
	if err != nil {
		return nil, nil, err
	}

	return endpoint, data, nil
}

// forecasterUsage displays the usage of the forecaster command and its
// subcommands.
func forecasterUsage() {
	fmt.Fprintf(os.Stderr, `Service that provides weather forecasts
Usage:
    %[1]s [globalflags] forecaster COMMAND [flags]

COMMAND:
    forecast: Retrieve weather forecast for a given location

Additional help:
    %[1]s forecaster COMMAND --help
`, os.Args[0])
}
func forecasterForecastUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] forecaster forecast -message JSON

Retrieve weather forecast for a given location
    -message JSON: 

Example:
    %[1]s forecaster forecast --message '{
      "lat": 37.8267,
      "long": -122.4233
   }'
`, os.Args[0])
}
