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
