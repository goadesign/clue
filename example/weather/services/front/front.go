package front

import (
	"context"
	"fmt"

	"github.com/crossnokaye/micro/example/weather/services/front/clients/forecast"
	"github.com/crossnokaye/micro/example/weather/services/front/clients/locator"
	genfront "github.com/crossnokaye/micro/example/weather/services/front/gen/front"
)

type (
	// Service is the front service implementation.
	Service struct {
		fc forecast.Client
		lc locator.Client
	}
)

// New instantiates a new front service.
func New(fc forecast.Client, lc locator.Client) *Service {
	return &Service{fc: fc, lc: lc}
}

// Forecast returns the forecast for the location at the given IP.
func (s *Service) Forecast(ctx context.Context, ip string) (*genfront.Forecast2, error) {
	l, err := s.lc.GetLocation(context.Background(), ip)
	if err != nil {
		return nil, err
	}
	if l.Country != "US" {
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
		State: l.RegionCode,
	}
	ps := make([]*genfront.Period, len(f.Periods))
	for i, p := range f.Periods {
		pval := genfront.Period(*p)
		ps[i] = &pval
	}
	return &genfront.Forecast2{Location: &loc, Periods: ps}, nil
}
