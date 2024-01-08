package tester

import (
	"context"

	"goa.design/clue/log"

	"goa.design/clue/example/weather/services/tester/clients/forecaster"
	"goa.design/clue/example/weather/services/tester/clients/locator"
)

type (
	Service struct {
		lc locator.Client
		fc forecaster.Client
		// Define test func maps for each collection tested (except TestAll)
		smokeTestMap      map[string]func(context.Context, *TestCollection)
		forecasterTestMap map[string]func(context.Context, *TestCollection)
		locatorTestMap    map[string]func(context.Context, *TestCollection)
	}
)

// New instantiates a new tester service.
func New(lc locator.Client, fc forecaster.Client) *Service {
	svc := &Service{
		lc: lc,
		fc: fc,
	}
	// initalize the test func maps with test function names as keys to funcs to be run.
	svc.smokeTestMap = map[string]func(context.Context, *TestCollection){
		"TestForecasterValidLatLong": svc.TestForecasterValidLatLong,
		"TestLocatorValidIP":         svc.TestLocatorValidIP,
	}
	svc.forecasterTestMap = map[string]func(context.Context, *TestCollection){
		"TestForecasterValidLatLong": svc.TestForecasterValidLatLong,
		"TestForecasterInvalidLat":   svc.TestForecasterInvalidLat,
		"TestForecasterInvalidLong":  svc.TestForecasterInvalidLong,
	}
	svc.locatorTestMap = map[string]func(context.Context, *TestCollection){
		"TestLocatorValidIP":   svc.TestLocatorValidIP,
		"TestLocatorInvalidIP": svc.TestLocatorInvalidIP,
	}

	return svc
}

func logError(ctx context.Context, err error) error {
	log.Error(ctx, err)
	return err
}
