// Code generated by goa v3.14.6, DO NOT EDIT.
//
// Forecaster gRPC client CLI support package
//
// Command:
// $ goa gen goa.design/clue/example/weather/services/forecaster/design -o
// services/forecaster

package client

import (
	"encoding/json"
	"fmt"

	forecaster "goa.design/clue/example/weather/services/forecaster/gen/forecaster"
	forecasterpb "goa.design/clue/example/weather/services/forecaster/gen/grpc/forecaster/pb"
)

// BuildForecastPayload builds the payload for the Forecaster forecast endpoint
// from CLI flags.
func BuildForecastPayload(forecasterForecastMessage string) (*forecaster.ForecastPayload, error) {
	var err error
	var message forecasterpb.ForecastRequest
	{
		if forecasterForecastMessage != "" {
			err = json.Unmarshal([]byte(forecasterForecastMessage), &message)
			if err != nil {
				return nil, fmt.Errorf("invalid JSON for message, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"lat\": 37.8267,\n      \"long\": -122.4233\n   }'")
			}
		}
	}
	v := &forecaster.ForecastPayload{
		Lat:  message.Lat,
		Long: message.Long,
	}

	return v, nil
}
