package front

import (
	"context"
	"fmt"

	"goa.design/clue/example/weather/services/front/clients/forecaster"
	"goa.design/clue/example/weather/services/front/clients/locator"
	"goa.design/clue/example/weather/services/front/clients/tester"
	genfront "goa.design/clue/example/weather/services/front/gen/front"
)

type (
	// Service is the front service implementation.
	Service struct {
		fc forecaster.Client
		lc locator.Client
		tc tester.Client
	}
)

// New instantiates a new front service.
func New(fc forecaster.Client, lc locator.Client, tc tester.Client) *Service {
	return &Service{fc: fc, lc: lc, tc: tc}
}

// Forecast returns the forecast for the location at the given IP.
func (s *Service) Forecast(ctx context.Context, ip string) (*genfront.Forecast2, error) {
	l, err := s.lc.GetLocation(ctx, ip)
	if err != nil {
		return nil, err
	}
	if l.Country != "United States" {
		return nil, genfront.MakeNotUsa(fmt.Errorf("IP not in the US (%s)", l.Country))
	}
	f, err := s.fc.GetForecast(ctx, l.Lat, l.Long)
	if err != nil {
		return nil, err
	}
	loc := genfront.Location{
		Lat:   l.Lat,
		Long:  l.Long,
		City:  l.City,
		State: l.Region,
	}
	ps := make([]*genfront.Period, len(f.Periods))
	for i, p := range f.Periods {
		pval := genfront.Period(*p)
		ps[i] = &pval
	}
	return &genfront.Forecast2{Location: &loc, Periods: ps}, nil
}

func (s *Service) TestAll(ctx context.Context, payload *genfront.TestAllPayload) (*genfront.TestResults, error) {
	tcPayload := &tester.TestAllPayload{Include: payload.Include, Exclude: payload.Exclude}
	return s.tc.TestAll(ctx, tcPayload)
}

func (s *Service) TestSmoke(ctx context.Context) (*genfront.TestResults, error) {
	return s.tc.TestSmoke(ctx)
}
