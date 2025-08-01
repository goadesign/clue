package mock

import "sync"

type (
	// Mock implementation of a service client.
	Mock struct {
		funcs   map[string]any
		seqs    map[string][]any
		indices []*index
		pos     int
		lock    sync.Mutex
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
		funcs: make(map[string]any),
		seqs:  make(map[string][]any),
	}
}

// If there is no mock left in the sequence then Next returns the permanent mock
// for name if any, nil otherwise.  If there are mocks left in the sequence then
// Next returns the next mock if its name is name, nil otherwise.
func (m *Mock) Next(name string) any {
	m.lock.Lock()
	defer m.lock.Unlock()

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
func (m *Mock) Add(name string, f any) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.indices = append(m.indices, &index{name, len(m.seqs[name])})
	m.seqs[name] = append(m.seqs[name], f)
}

// Set a permanent mock for the function with the given name.
func (m *Mock) Set(name string, f any) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.funcs[name] = f
}

// HasMore returns true if the mock sequence isn't fully consumed.
func (m *Mock) HasMore() bool {
	m.lock.Lock()
	defer m.lock.Unlock()

	return m.pos < len(m.indices)
}
