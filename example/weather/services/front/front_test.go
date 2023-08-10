package front

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"goa.design/clue/example/weather/services/front/clients/forecaster"
	mockforecaster "goa.design/clue/example/weather/services/front/clients/forecaster/mocks"
	"goa.design/clue/example/weather/services/front/clients/locator"
	mocklocator "goa.design/clue/example/weather/services/front/clients/locator/mocks"
	genfront "goa.design/clue/example/weather/services/front/gen/front"
)

func TestForecast(t *testing.T) {
	cases := []struct {
		name         string
		locationFunc mocklocator.ClientGetLocationFunc
		forecastFunc mockforecaster.ClientGetForecastFunc

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
			lmock := mocklocator.NewClient(t)
			lmock.AddGetLocation(c.locationFunc)
			fmock := mockforecaster.NewClient(t)
			fmock.AddGetForecast(c.forecastFunc)
			s := New(fmock, lmock)
			result, err := s.Forecast(context.Background(), testIP)
			if c.expectedError != nil {
				assert.Nil(t, result)
				require.NotNil(t, err)
				assert.Equal(t, c.expectedError.Error(), err.Error())
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, c.expectedResult, result)
			assert.False(t, lmock.HasMore())
			assert.False(t, fmock.HasMore())
		})
	}
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
		Periods: []*genfront.Period{{
			Name:            "morning",
			StartTime:       "2022-01-22T21:57:40+00:00",
			EndTime:         "2022-01-22T21:57:40+00:00",
			Temperature:     10,
			TemperatureUnit: "C",
			Summary:         "cool",
		}},
	}

	errForecast = fmt.Errorf("test forecast error")
	errLocation = fmt.Errorf("test location error")
)

func getLocationInUSFunc(t *testing.T) mocklocator.ClientGetLocationFunc {
	return func(ctx context.Context, ip string) (*locator.WorldLocation, error) {
		assert.Equal(t, testIP, ip)
		return &locator.WorldLocation{
			Lat:     testLocation.Lat,
			Long:    testLocation.Long,
			City:    testLocation.City,
			Region:  testLocation.State,
			Country: "United States",
		}, nil
	}
}

func getLocationNotInUSFunc(t *testing.T) mocklocator.ClientGetLocationFunc {
	return func(ctx context.Context, ip string) (*locator.WorldLocation, error) {
		assert.Equal(t, testIP, ip)
		return &locator.WorldLocation{
			Lat:     testLocation.Lat,
			Long:    testLocation.Long,
			City:    testLocation.City,
			Region:  testLocation.State,
			Country: "NOT US",
		}, nil
	}
}

func getLocationErrorFunc(t *testing.T) mocklocator.ClientGetLocationFunc {
	return func(ctx context.Context, ip string) (*locator.WorldLocation, error) {
		assert.NotEmpty(t, ip)
		return nil, errLocation
	}
}

func getForecastFunc(t *testing.T) mockforecaster.ClientGetForecastFunc {
	return func(ctx context.Context, lat float64, long float64) (*forecaster.Forecast, error) {
		assert.Equal(t, testLocation.Lat, lat)
		assert.Equal(t, testLocation.Long, long)
		lval := forecaster.Location(*testForecast.Location)
		ps := make([]*forecaster.Period, len(testForecast.Periods))
		for i, p := range testForecast.Periods {
			pval := forecaster.Period(*p)
			ps[i] = &pval
		}
		return &forecaster.Forecast{Location: &lval, Periods: ps}, nil
	}
}

func getForecastErrorFunc(t *testing.T) mockforecaster.ClientGetForecastFunc {
	return func(ctx context.Context, lat float64, long float64) (*forecaster.Forecast, error) {
		return nil, errForecast
	}
}
