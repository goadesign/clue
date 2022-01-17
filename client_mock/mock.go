/*
Package mock makes it possible to implement in-memory service client mocks.

Mock service implementations make use of the package to capture and replay
sequences of service method calls. Tests can record a sequence or set permanent
mocks. When using sequences tests can verify that the entire sequence has been
consumed. Sequences and permanent mocks can be used at the same time in which
case mock method calls first consume the sequence then call the permanent mock.

Sequences are created using the Add method while permanent mocks are created
using Set. The method HasMore returns true if a sequence has been set and has
not been entirely consumed.

Example mock client implementation:

    type mock struct {
	    *mock.Mock
	    t *testing.T
    }

    func newMock(t *testing.T) *mock {
        return &mock{mock.New(), t}
    }

    func (m *mock) GetPrices(ctx context.Context, first, last time.Time, nodeID string) ([]*Price, error) {
        if f := m.Next("GetPrices"); f != nil {
            return f.(func(ctx, time.Time, time.Time, string))(ctx, first, last, nodeID)
	}
	m.t.Error("unexpected GetPrices call")
    }

Example usage in tests:

	// Create a new mock client.
	mock := newMock(t)

	// Add a mock for the GetPrices method.
	mock.Add("GetPrices", getPricesFunc)

	// Call the mock.
	prices, err := mock.GetPrices(ctx, firstHour, lastHour, nodeID)

	// Validate prices and err
	if err != nil {
		t.Errorf("GetPrices returned %v", err)
	}

	// Make sure entire sequence has been consumed.
	if mock.HasMore() {
		t.Error("client method sequence not fully consumed")
	}
*/
package mock

type (
	// Mock implementation of a service client.
	Mock struct {
		funcs   map[string]interface{}
		seqs    map[string][]interface{}
		indices []*index
		pos     int
	}

	// index identifies a mock in a sequence.
	index struct {
		name string
		pos  int
	}
)

// New returns a new mock client.
func New() *Mock {
	return &Mock{
		funcs: make(map[string]interface{}),
		seqs:  make(map[string][]interface{}),
	}
}

// Next returns the next mock for the function with the given name. It first
// consumes any sequence then returns any permanent mock. Next returns nil if
// there is no sequence or the sequence has been fully consumed, and there is
// no permanent mock.
func (m *Mock) Next(name string) interface{} {
	if m.pos < len(m.indices) && len(m.seqs[name]) > 0 {
		idx := m.indices[m.pos]
		if idx.name != name || idx.pos >= len(m.seqs[name]) {
			// There is a sequence but the wrong method is being called
			return nil
		}
		f := m.seqs[name][idx.pos]
		idx.pos++
		m.pos++
		return f
	}
	// No sequence or sequence fully consumed - look for permanent mock
	if f, ok := m.funcs[name]; ok {
		return f
	}
	return nil
}

// Add adds f to the mock sequence.
func (m *Mock) Add(name string, f interface{}) {
	m.indices = append(m.indices, &index{name, len(m.seqs[name])})
	m.seqs[name] = append(m.seqs[name], f)
}

// Set a permanent mock for the function with the given name.
func (m *Mock) Set(name string, f interface{}) {
	m.funcs[name] = f
}

// HasMore returns true if the mock sequence isn't fully consumed.
func (m *Mock) HasMore() bool {
	return m.pos < len(m.indices)
}
