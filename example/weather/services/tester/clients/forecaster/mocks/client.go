// Code generated by Clue Mock Generator v0.18.2, DO NOT EDIT.
//
// Command:
// $ cmg gen goa.design/clue/example/weather/services/tester/clients/forecaster

package mockforecaster

import (
	"context"
	"testing"

	"goa.design/clue/mock"

	"goa.design/clue/example/weather/services/tester/clients/forecaster"
)

type (
	Client struct {
		m *mock.Mock
		t *testing.T
	}

	ClientGetForecastFunc func(ctx context.Context, lat, long float64) (*forecaster.Forecast, error)
)

func NewClient(t *testing.T) *Client {
	var (
		m                   = &Client{mock.New(), t}
		_ forecaster.Client = m
	)
	return m
}

func (m *Client) AddGetForecast(f ClientGetForecastFunc) {
	m.m.Add("GetForecast", f)
}

func (m *Client) SetGetForecast(f ClientGetForecastFunc) {
	m.m.Set("GetForecast", f)
}

func (m *Client) GetForecast(ctx context.Context, lat, long float64) (*forecaster.Forecast, error) {
	if f := m.m.Next("GetForecast"); f != nil {
		return f.(ClientGetForecastFunc)(ctx, lat, long)
	}
	m.t.Helper()
	m.t.Error("unexpected GetForecast call")
	return nil, nil
}

func (m *Client) HasMore() bool {
	return m.m.HasMore()
}
