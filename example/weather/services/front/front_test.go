package front

import (
	"context"
	"fmt"
	"testing"

	"github.com/goadesign/clue/example/weather/services/front/clients/forecaster"
	"github.com/goadesign/clue/example/weather/services/front/clients/locator"
	genfront "github.com/goadesign/clue/example/weather/services/front/gen/front"
)

func TestForecast(t *testing.T) {
	cases := []struct {
		name         string
		locationFunc locator.GetLocationFunc
		forecastFunc forecaster.GetForecastFunc

		expectedResult *genfront.Forecast2
		expectedError  error
	}{
		{"success", getLocationInUSFunc(t), getForecastFunc(t), testForecast, nil},
		{"forecast error", getLocationInUSFunc(t), getForecastErrorFunc(t), nil, errForecast},
		{"location not in US", getLocationNotInUSFunc(t), nil, nil, genfront.MakeNotUsa(fmt.Errorf("IP not in the US (NOT US)"))},
		{"location error", getLocationErrorFunc(t), nil, nil, errLocation},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			lmock := locator.NewMock(t)
			lmock.AddGetLocationFunc(c.locationFunc)
			fmock := forecaster.NewMock(t)
			fmock.AddGetForecastFunc(c.forecastFunc)
			s := New(fmock, lmock)
			result, err := s.Forecast(context.Background(), testIP)
			if (c.expectedError != nil) && (err.Error() != c.expectedError.Error()) {
				t.Errorf("Forecast: got error %s, expected %s", err, c.expectedError)
			}
			if !equal(result, c.expectedResult) {
				t.Errorf("Forecast: got %#v, expected %v", result, c.expectedResult)
			}
		})
	}
}

func equal(a, b *genfront.Forecast2) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if a.Location.Lat != b.Location.Lat {
		return false
	}
	if a.Location.Long != b.Location.Long {
		return false
	}
	if a.Location.City != b.Location.City {
		return false
	}
	if a.Location.State != b.Location.State {
		return false
	}
	if len(a.Periods) != len(b.Periods) {
		return false
	}
	for i, p := range a.Periods {
		if p.Name != b.Periods[i].Name {
			return false
		}
		if p.StartTime != b.Periods[i].StartTime {
			return false
		}
		if p.EndTime != b.Periods[i].EndTime {
			return false
		}
		if p.Temperature != b.Periods[i].Temperature {
			return false
		}
		if p.TemperatureUnit != b.Periods[i].TemperatureUnit {
			return false
		}
		if p.Summary != b.Periods[i].Summary {
			return false
		}
	}
	return true
}

var (
	testIP       = "8.8.8.8"
	testLocation = &genfront.Location{
		Lat:   23.0,
		Long:  -32.0,
		City:  "Test City",
		State: "Test State",
	}

	testForecast = &genfront.Forecast2{
		Location: testLocation,
		Periods: []*genfront.Period{
			{"morning", "2022-01-22T21:57:40+00:00", "2022-01-22T21:57:40+00:00", 10, "C", "cool"},
		},
	}

	errForecast = fmt.Errorf("test forecast error")
	errLocation = fmt.Errorf("test location error")
)

func getLocationInUSFunc(t *testing.T) locator.GetLocationFunc {
	return func(ctx context.Context, ip string) (*locator.WorldLocation, error) {
		if ip != testIP {
			t.Errorf("GetLocation: got %s, expected %s", ip, testIP)
			return nil, nil
		}
		return &locator.WorldLocation{
			Lat:     testLocation.Lat,
			Long:    testLocation.Long,
			City:    testLocation.City,
			Region:  testLocation.State,
			Country: "United States",
		}, nil
	}
}

func getLocationNotInUSFunc(t *testing.T) locator.GetLocationFunc {
	return func(ctx context.Context, ip string) (*locator.WorldLocation, error) {
		if ip != testIP {
			t.Errorf("GetLocation: got %s, expected %s", ip, testIP)
			return nil, nil
		}
		return &locator.WorldLocation{
			Lat:     testLocation.Lat,
			Long:    testLocation.Long,
			City:    testLocation.City,
			Region:  testLocation.State,
			Country: "NOT US",
		}, nil
	}
}

func getLocationErrorFunc(t *testing.T) locator.GetLocationFunc {
	return func(ctx context.Context, ip string) (*locator.WorldLocation, error) {
		if ip == "" {
			t.Error("GetLocation: expected non-empty IP")
			return nil, nil
		}
		return nil, errLocation
	}
}

func getForecastFunc(t *testing.T) forecaster.GetForecastFunc {
	return func(ctx context.Context, lat float64, long float64) (*forecaster.Forecast, error) {
		if lat != testLocation.Lat || long != testLocation.Long {
			t.Errorf("GetForecast: expected (%f, %f), got (%f, %f)", testLocation.Lat, testLocation.Long, lat, long)
			return nil, nil
		}
		lval := forecaster.Location(*testForecast.Location)
		ps := make([]*forecaster.Period, len(testForecast.Periods))
		for i, p := range testForecast.Periods {
			pval := forecaster.Period(*p)
			ps[i] = &pval
		}
		return &forecaster.Forecast{Location: &lval, Periods: ps}, nil
	}
}

func getForecastErrorFunc(t *testing.T) forecaster.GetForecastFunc {
	return func(ctx context.Context, lat float64, long float64) (*forecaster.Forecast, error) {
		return nil, errForecast
	}
}
