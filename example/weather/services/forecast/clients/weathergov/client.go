/*
Package weathergov provides a client for the Weather.gov API described at
https://www.weather.gov/documentation/services-web-api#
*/
package weathergov

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type (
	// Client is a client for the Weather.gov API.
	Client interface {
		// GetForecast gets the forecast for the given location.
		GetForecast(ctx context.Context, lat, long float64) (*Forecast, error)
		// Name provides a client name used to report health check issues.
		Name() string
		// Ping checks the client is healthy.
		Ping(ctx context.Context) error
	}

	// Forecast represents the forecast for a given location.
	Forecast struct {
		// Location is the location of the forecast.
		Location *Location
		// Periods is the forecast for the location.
		Periods []*Period
	}

	// Location represents the geographical location of a forecast.
	Location struct {
		// Lat is the latitude of the location.
		Lat float64
		// Long is the longitude of the location.
		Long float64
		// City is the city of the location.
		City string
		// State is the state of the location.
		State string
	}

	// Period represents a forecast period.
	Period struct {
		// Name is the name of the forecast period.
		Name string
		// StartTime is the start time of the forecast period in RFC3339 format.
		StartTime string
		// EndTime is the end time of the forecast period in RFC3339 format.
		EndTime string
		// Temperature is the temperature of the forecast period.
		Temperature int
		// TemperatureUnit is the temperature unit of the forecast period.
		TemperatureUnit string
		// Summary is the summary of the forecast period.
		Summary string
	}

	// client implements Client.
	client struct {
		c *http.Client
	}

	// point represents the response from the Weather.gov API /point endpoint.
	point struct {
		ID         string `json:"id"`
		Properties struct {
			ForecastURL string `json:"forecast"`
			Location    struct {
				Properties struct {
					City  string `json:"city"`
					State string `json:"state"`
				} `json:"properties"`
			} `json:"relativeLocation"`
		}
	}

	// forecast represents the response from the Weather.gov API for a forecast URL.
	forecast struct {
		Properties struct {
			Periods []struct {
				Name            string `json:"name"`
				StartTime       string `json:"startTime"`
				EndTime         string `json:"endTime"`
				Temperature     int    `json:"temperature"`
				TemperatureUnit string `json:"temperatureUnit"`
				ShortForecast   string `json:"shortForecast"`
			} `json:"periods"`
		} `json:"properties"`
	}
)

// baseURL is the base URL for the Weather.gov API.
const baseURL = "https://api.weather.gov"

// New returns a new client for the Weather.gov API.
func New(c *http.Client) Client {
	return &client{c: c}
}

// GetForecast gets the forecast for the given location.
func (c *client) GetForecast(ctx context.Context, lat, long float64) (*Forecast, error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/points/%f,%f", baseURL, lat, long), nil)
	resp, err := c.doWithRetries(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var p point
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		return nil, err
	}

	req, err = http.NewRequest("GET", p.Properties.ForecastURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to parse forecastURL %q: %s", p.Properties.ForecastURL, err)
	}
	resp, err = c.doWithRetries(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var f forecast
	if err := json.NewDecoder(resp.Body).Decode(&f); err != nil {
		return nil, err
	}
	return &Forecast{
		Location: &Location{
			Lat:   lat,
			Long:  long,
			City:  p.Properties.Location.Properties.City,
			State: p.Properties.Location.Properties.State,
		},
		Periods: toPeriods(f.Properties.Periods),
	}, nil

}

// Name provides a client name used to report health check issues.
func (c *client) Name() string {
	return "weathergov"
}

// Ping checks the client is healthy.
func (c *client) Ping(ctx context.Context) error {
	req, _ := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/points/%f,%f", baseURL, 34.4239263, -119.7068831), nil)
	resp, err := c.doWithRetries(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (c *client) doWithRetries(req *http.Request) (*http.Response, error) {
	resp, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}
	retries := 0
	for resp.StatusCode != http.StatusOK && retries < 10 {
		resp.Body.Close()
		time.Sleep(time.Second)
		resp, err = c.c.Do(req)
		if err != nil {
			return nil, err
		}
		retries++
	}
	if resp.StatusCode != http.StatusOK {
		msg, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			msg = []byte("unknown error")
		}
		resp.Body.Close()
		return nil, fmt.Errorf("unexpected response status code %d (%s)", resp.StatusCode, string(msg))
	}
	return resp, nil
}

func toPeriods(ps []struct {
	Name            string `json:"name"`
	StartTime       string `json:"startTime"`
	EndTime         string `json:"endTime"`
	Temperature     int    `json:"temperature"`
	TemperatureUnit string `json:"temperatureUnit"`
	ShortForecast   string `json:"shortForecast"`
}) []*Period {
	var periods []*Period
	for _, p := range ps {
		periods = append(periods, &Period{
			Name:            p.Name,
			StartTime:       p.StartTime,
			EndTime:         p.EndTime,
			Temperature:     p.Temperature,
			TemperatureUnit: p.TemperatureUnit,
			Summary:         p.ShortForecast,
		})
	}
	return periods
}
