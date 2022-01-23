package forecast

import (
	"context"
	"testing"

	"github.com/crossnokaye/micro/mock"
)

type (
	// GetForecastFunc mocks the GetForecast method.
	GetForecastFunc func(ctx context.Context, lat, long float64) (*Forecast, error)

	// Mock implementation of the forecast client.
	Mock struct {
		m *mock.Mock
		t *testing.T
	}
)

var _ Client = &Mock{}

// NewMock returns a new mock client.
func NewMock(t *testing.T) *Mock {
	return &Mock{mock.New(), t}
}

// AddGetForecastFunc adds f to the mocked call sequence.
func (m *Mock) AddGetForecastFunc(f GetForecastFunc) {
	m.m.Add("GetForecast", f)
}

// SetGetForecastFunc sets f for all calls to the mocked method.
func (m *Mock) SetGetForecastFunc(f GetForecastFunc) {
	m.m.Set("GetForecast", f)
}

// GetForecast implements the Client interface.
func (m *Mock) GetForecast(ctx context.Context, lat, long float64) (*Forecast, error) {
	if f := m.m.Next("GetForecast"); f != nil {
		return f.(GetForecastFunc)(ctx, lat, long)
	}
	m.t.Error("unexpected call to GetForecast")
	return nil, nil
}

// HasMore returns true if there are more calls to be made.
func (m *Mock) HasMore() bool {
	return m.m.HasMore()
}
