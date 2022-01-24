package front

import (
	"context"
	"fmt"

	"github.com/crossnokaye/micro/example/weather/services/front/clients/forecaster"
	"github.com/crossnokaye/micro/example/weather/services/front/clients/locator"
	genfront "github.com/crossnokaye/micro/example/weather/services/front/gen/front"
)

type (
	// Service is the front service implementation.
	Service struct {
		fc forecaster.Client
		lc locator.Client
	}
)

// New instantiates a new front service.
func New(fc forecaster.Client, lc locator.Client) *Service {
	return &Service{fc: fc, lc: lc}
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
