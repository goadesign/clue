package forecaster

import (
	"context"
	"fmt"
	"testing"

	"github.com/goadesign/clue/example/weather/services/forecaster/clients/weathergov"
	genforecaster "github.com/goadesign/clue/example/weather/services/forecaster/gen/forecaster"
)

func TestForecast(t *testing.T) {
	// Create mock call sequence with first successful call returning a
	// forecast then failing.
	wc := weathergov.NewMock(t)
	wc.AddGetForecastFunc(func(ctx context.Context, lat, long float64) (*weathergov.Forecast, error) {
		// Make sure service passes the right values to the client
		if lat != latitude {
			t.Errorf("got latitude %f, expected %f", lat, latitude)
		}
		if long != longitude {
			t.Errorf("got longitude %f, expected %f", long, longitude)
		}
		return mockForecast, nil
	})
	wc.AddGetForecastFunc(func(_ context.Context, _, _ float64) (*weathergov.Forecast, error) {
		return nil, fmt.Errorf("test failure")
	})

	// Create service using mock client.
	s := New(wc)

	// Call service, first call should succeed.
	f, err := s.Forecast(context.Background(), &genforecaster.ForecastPayload{Lat: latitude, Long: longitude})
	if err != nil {
		t.Errorf("got error %v, expected nil", err)
	}
	if f.Location.Lat != latitude {
		t.Errorf("got latitude %f, expected %f", f.Location.Lat, latitude)
	}
	if f.Location.Long != longitude {
		t.Errorf("got longitude %f, expected %f", f.Location.Long, longitude)
	}
	if len(f.Periods) != len(mockForecast.Periods) {
		t.Errorf("got %d periods, expected %d", len(f.Periods), len(mockForecast.Periods))
	}
	if f.Periods[0].Name != mockForecast.Periods[0].Name {
		t.Errorf("got name %s, expected %s", f.Periods[0].Name, mockForecast.Periods[0].Name)
	}
	if f.Periods[0].StartTime != mockForecast.Periods[0].StartTime {
		t.Errorf("got start time %s, expected %s", f.Periods[0].StartTime, mockForecast.Periods[0].StartTime)
	}
	if f.Periods[0].EndTime != mockForecast.Periods[0].EndTime {
		t.Errorf("got end time %s, expected %s", f.Periods[0].EndTime, mockForecast.Periods[0].EndTime)
	}
	if f.Periods[0].Temperature != mockForecast.Periods[0].Temperature {
		t.Errorf("got temperature %d, expected %d", f.Periods[0].Temperature, mockForecast.Periods[0].Temperature)
	}
	if f.Periods[0].TemperatureUnit != mockForecast.Periods[0].TemperatureUnit {
		t.Errorf("got temperature unit %s, expected %s", f.Periods[0].TemperatureUnit, mockForecast.Periods[0].TemperatureUnit)
	}
	if f.Periods[0].Summary != mockForecast.Periods[0].Summary {
		t.Errorf("got summary %s, expected %s", f.Periods[0].Summary, mockForecast.Periods[0].Summary)
	}

	// Second call should fail
	_, err = s.Forecast(context.Background(), &genforecaster.ForecastPayload{Lat: latitude, Long: longitude})
	if err == nil {
		t.Errorf("got nil error, expected non-nil")
	}

	// Make sure all calls were made.
	if wc.HasMore() {
		t.Error("expected all calls to weather service to be made")
	}
}

const (
	latitude  = 37.8267
	longitude = -122.4233
)

var mockForecast = &weathergov.Forecast{
	Location: &weathergov.Location{
		Lat:   latitude,
		Long:  longitude,
		City:  "Mountain View",
		State: "CA",
	},
	Periods: []*weathergov.Period{
		{
			Name:            "Tonight",
			StartTime:       "2019-05-29T21:00:00-07:00",
			EndTime:         "2019-05-29T23:00:00-07:00",
			Temperature:     63.0,
			TemperatureUnit: "F",
			Summary:         "Mostly cloudy starting tonight.",
		},
	},
}
