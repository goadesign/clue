package forecaster

import (
	"context"

	"github.com/crossnokaye/micro/example/weather/services/forecaster/clients/weathergov"
	genforecaster "github.com/crossnokaye/micro/example/weather/services/forecaster/gen/forecaster"
)

type (
	// Service is the forecast service implementation.
	Service struct {
		wc weathergov.Client
	}
)

// New instantiates a new forecast service.
func New(wc weathergov.Client) *Service {
	return &Service{wc: wc}
}

// Forecast returns the forecast for the given location.
func (s *Service) Forecast(ctx context.Context, p *genforecaster.ForecastPayload) (*genforecaster.Forecast2, error) {
	f, err := s.wc.GetForecast(ctx, p.Lat, p.Long)
	if err != nil {
		return nil, err
	}
	location := &genforecaster.Location{
		Lat:   f.Location.Lat,
		Long:  f.Location.Long,
		City:  f.Location.City,
		State: f.Location.State,
	}
	periods := make([]*genforecaster.Period, len(f.Periods))
	for i, p := range f.Periods {
		periods[i] = &genforecaster.Period{
			Name:            p.Name,
			StartTime:       p.StartTime,
			EndTime:         p.EndTime,
			Temperature:     p.Temperature,
			TemperatureUnit: p.TemperatureUnit,
			Summary:         p.Summary,
		}
	}
	return &genforecaster.Forecast2{Location: location, Periods: periods}, nil
}
