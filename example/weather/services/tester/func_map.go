package tester

import "context"

var smokeTestMap = make(map[string]func(context.Context, *TestCollection))

func (svc *Service) smokeTestMapInit(ctx context.Context) {
	smokeTestMap["TestValidLatLong"] = svc.TestForecasterValidLatLong
	smokeTestMap["TestLocatorValidIP"] = svc.TestLocatorValidIP
}

var forecasterTestMap = make(map[string]func(context.Context, *TestCollection))

func (svc *Service) forecasterTestMapInit(ctx context.Context) {
	forecasterTestMap["TestValidLatLong"] = svc.TestForecasterValidLatLong
	forecasterTestMap["TestInvalidLat"] = svc.TestForecasterInvalidLat
	forecasterTestMap["TestInvalidLong"] = svc.TestForecasterInvalidLong
}

var locatorTestMap = make(map[string]func(context.Context, *TestCollection))

func (svc *Service) locatorTestMapInit(ctx context.Context) {
	locatorTestMap["TestLocatorValidIP"] = svc.TestLocatorValidIP
	locatorTestMap["TestLocatorInvalidIP"] = svc.TestLocatorInvalidIP
}
