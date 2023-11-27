package tester

// The func_map.go file contains the maps that map test names to the actual test functions. Test
// functions are defined in the integration.go files for each service. The maps are used in the
// tester service to run the tests. Names are used for filtering which tests to run when TestAll
// is called.
//
// The smokeTestMap is a map that associates smoke test names with their corresponding test
// functions.
//
// Similarly, forecasterTestMap and locatorTestMap are maps that associate forecaster and locator
// test names with their corresponding test functions.
//
// The smokeTestMap is used exclusively when TestSmoke is called in the tester service. It will
// contain a subset of the tests defined in the service-specific test maps (such as those for
// locator and forecaster).

import "context"

var smokeTestMap = make(map[string]func(context.Context, *TestCollection))

func (svc *Service) smokeTestMapInit(ctx context.Context) {
	smokeTestMap["TestForecasterValidLatLong"] = svc.TestForecasterValidLatLong
	smokeTestMap["TestLocatorValidIP"] = svc.TestLocatorValidIP
}

var forecasterTestMap = make(map[string]func(context.Context, *TestCollection))

func (svc *Service) forecasterTestMapInit(ctx context.Context) {
	forecasterTestMap["TestForecasterValidLatLong"] = svc.TestForecasterValidLatLong
	forecasterTestMap["TestForecasterInvalidLat"] = svc.TestForecasterInvalidLat
	forecasterTestMap["TestForecasterInvalidLong"] = svc.TestForecasterInvalidLong
}

var locatorTestMap = make(map[string]func(context.Context, *TestCollection))

func (svc *Service) locatorTestMapInit(ctx context.Context) {
	locatorTestMap["TestLocatorValidIP"] = svc.TestLocatorValidIP
	locatorTestMap["TestLocatorInvalidIP"] = svc.TestLocatorInvalidIP
}
