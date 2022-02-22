
# mock: Downstream Mocking

[![Build Status](https://github.com/goadesign/clue/workflows/CI/badge.svg?branch=main&event=push)](https://github.com/goadesign/clue/actions?query=branch%3Amain+event%3Apush)
[![Go Reference](https://pkg.go.dev/badge/goa.design/clue/mock.svg)](https://pkg.go.dev/goa.design/clue/mock)

## Overview

Package `mock` makes it possible to implement downstream service
dependency mocks that run in-memory.

Conceptually a test may want to verify that a service method is called multiple
times with different values each time. The mock package `Add` method can be used
by the test to record each method call and verify that the correct values are
used. 

Alternatively (or complementarily) the `Set` method defines a mock that
is used for all method calls by the test. If both `Add` and `Set` are used by
the test then the mocks recorded by `Add` are used first.

## Usage

### Creating a Dependency Mock

```go
// mock implementation of the `prices` service
type mock struct {
        *mock.Mock    // Embed the mock package Mock struct which provides the mock API
        t *testing.T
}

func newMock(t *testing.T) *mock {
        return &mock{mock.New(), t}
}

// mock implementation of the GetPrices method. The implementation leverages
// the mock package to replay a sequence of calls.
func (m *mock) GetPrices(ctx context.Context, first, last time.Time, nodeID string) ([]*Price, error) {
        if f := m.Next("GetPrices"); f != nil { // Get the next mock in the sequence (or the permanent mock)
                return f.(func(ctx, time.Time, time.Time, string) ([]*Price, error))(ctx, first, last, nodeID)
        }
        m.t.Error("unexpected GetPrices call")
}
```

### Using a Mock in a Test

Sequences are created using the `Add` method while permanent mocks are created
using `Set`. The method `HasMore` returns true if a sequence has been set and
has not been entirely consumed.

```go
// Create a new mock client (defined above).
mock := newMock(t)

// Add a mock call for the GetPrices method.
mock.Add("GetPrices", getPricesFunc)

// Call the mock.
prices, err := mock.GetPrices(ctx, firstHour, lastHour, nodeID)

// Validate prices and err
if err != nil {
	t.Errorf("GetPrices returned %v", err)
}

// Make sure entire sequence has been consumed (in this example there is only
// one call).
if mock.HasMore() {
	t.Error("GetPrices was not called")
}
```
