package log

type (
	// KV represents a key/value pair. Values must be strings, numbers,
	// booleans, nil or a slice of these types.
	KV struct {
		K string
		V interface{}
	}

	// Fielder is an interface that will return a slice of KV
	Fielder interface {
		LogFields() []KV
	}

	// Fields allows to quickly define fields for cases where you are OK with
	// non-deterministic order of the fields
	Fields map[string]interface{}

	kvList []KV
)

func (kv KV) LogFields() []KV {
	return []KV{kv}
}

func (f Fields) LogFields() []KV {
	fields := make([]KV, 0, len(f))
	for k, v := range f {
		fields = append(fields, KV{k, v})
	}
	return fields
}

func (kvs kvList) merge(fielders []Fielder) kvList {
	for _, fielder := range fielders {
		switch fielder := fielder.(type) {
		case KV:
			// Avoid unnecessary allocation the slice from KV.LogFields
			kvs = append(kvs, fielder)
		default:
			kvs = append(kvs, fielder.LogFields()...)
		}
	}
	return kvs
}

func (kvs kvList) LogFields() []KV {
	return kvs
}
