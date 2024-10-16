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
	totalLen := len(kvs)
	cachedFields := make([][]KV, len(fielders))
	for i, fielder := range fielders {
		if _, ok := fielder.(KV); ok {
			totalLen++
		} else {
			fields := fielder.LogFields()
			cachedFields[i] = fields
			totalLen += len(fields)
		}
	}
	result := make(kvList, len(kvs), totalLen)
	copy(result, kvs)
	for i, fielder := range fielders {
		if kv, ok := fielder.(KV); ok {
			result = append(result, kv)
		} else {
			result = append(result, cachedFields[i]...)
		}
	}
	return result
}

func (kvs kvList) LogFields() []KV {
	return kvs
}
