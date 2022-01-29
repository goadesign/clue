package locator

import (
	"context"
	"testing"

	"goa.design/clue/mock"
)

type (
	// GetLocationFunc mocks the GetIPLocation method.
	GetLocationFunc func(ctx context.Context, ip string) (*WorldLocation, error)

	// Mock implementation of the locator client.
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

// AddGetIPLocationFunc adds f to the mocked call sequence.
func (m *Mock) AddGetLocationFunc(f GetLocationFunc) {
	m.m.Add("GetLocation", f)
}

// SetGetIPLocationFunc sets f for all calls to the mocked method.
func (m *Mock) SetGetLocationFunc(f GetLocationFunc) {
	m.m.Set("GetLocation", f)
}

// GetLocation implements the Client interface.
func (m *Mock) GetLocation(ctx context.Context, ip string) (*WorldLocation, error) {
	if f := m.m.Next("GetLocation"); f != nil {
		return f.(GetLocationFunc)(ctx, ip)
	}
	m.t.Error("unexpected call to GetLocation")
	return nil, nil
}

// HasMore returns true if there are more calls to be made.
func (m *Mock) HasMore() bool {
	return m.m.HasMore()
}
