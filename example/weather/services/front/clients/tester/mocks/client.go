// Code generated by Clue Mock Generator v0.18.2, DO NOT EDIT.
//
// Command:
// $ cmg gen goa.design/clue/example/weather/services/front/clients/tester

package mocktester

import (
	"context"
	"testing"

	"goa.design/clue/mock"

	"goa.design/clue/example/weather/services/front/clients/tester"
	"goa.design/clue/example/weather/services/front/gen/front"
)

type (
	Client struct {
		m *mock.Mock
		t *testing.T
	}

	ClientTestAllFunc   func(ctx context.Context, included, excluded []string) (*front.TestResults, error)
	ClientTestSmokeFunc func(ctx context.Context) (*front.TestResults, error)
)

func NewClient(t *testing.T) *Client {
	var (
		m               = &Client{mock.New(), t}
		_ tester.Client = m
	)
	return m
}

func (m *Client) AddTestAll(f ClientTestAllFunc) {
	m.m.Add("TestAll", f)
}

func (m *Client) SetTestAll(f ClientTestAllFunc) {
	m.m.Set("TestAll", f)
}

func (m *Client) TestAll(ctx context.Context, included, excluded []string) (*front.TestResults, error) {
	if f := m.m.Next("TestAll"); f != nil {
		return f.(ClientTestAllFunc)(ctx, included, excluded)
	}
	m.t.Helper()
	m.t.Error("unexpected TestAll call")
	return nil, nil
}

func (m *Client) AddTestSmoke(f ClientTestSmokeFunc) {
	m.m.Add("TestSmoke", f)
}

func (m *Client) SetTestSmoke(f ClientTestSmokeFunc) {
	m.m.Set("TestSmoke", f)
}

func (m *Client) TestSmoke(ctx context.Context) (*front.TestResults, error) {
	if f := m.m.Next("TestSmoke"); f != nil {
		return f.(ClientTestSmokeFunc)(ctx)
	}
	m.t.Helper()
	m.t.Error("unexpected TestSmoke call")
	return nil, nil
}

func (m *Client) HasMore() bool {
	return m.m.HasMore()
}
