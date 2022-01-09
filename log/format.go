package log

import (
	"bytes"
	"fmt"
	"strconv"
)

// DefaultFormat is the default log formatter, it prints entries in the
// following format:
//
//    [LEVEL] [key=val key=val ...] message
func DefaultFormat(e *Entry) []byte {
	var b bytes.Buffer
	b.WriteByte('[')
	b.WriteString(e.Level.String())
	b.WriteByte(']')
	if len(e.KeyVals) > 0 {
		b.WriteByte(' ')
		b.WriteByte('[')
		keys, vals := e.KeyVals.Parse()
		for i := 0; i < len(keys); i++ {
			b.WriteString(keys[i])
			b.WriteByte('=')
			b.WriteString(fmt.Sprintf("%v", vals[i]))
			if i < len(keys)-1 {
				b.WriteByte(' ')
			}
		}
		b.WriteByte(']')
	}
	b.WriteByte(' ')
	b.WriteString(e.Message + "\n")
	return b.Bytes()
}

// JSONFormat is a log formatter that prints entries using JSON.
//
// note: the implementation avoids using reflection (and thus the json package)
// for efficiency.
func JSONFormat(e *Entry) []byte {
	var b bytes.Buffer
	b.WriteByte('{')
	b.WriteString(`"level":`)
	b.WriteString(`"`)
	b.WriteString(e.Level.String())
	b.WriteString(`",`)
	b.WriteString(`"message":`)
	b.WriteString(`"`)
	b.WriteString(e.Message)
	b.WriteString(`"`)
	if len(e.KeyVals) > 0 {
		b.WriteByte(',')
		keys, vals := e.KeyVals.Parse()
		for i := 0; i < len(keys); i++ {
			b.WriteString(`"`)
			b.WriteString(keys[i])
			b.WriteString(`":`)
			writeJSON(vals[i], &b)
			if i < len(keys)-1 {
				b.WriteByte(',')
			}
		}
	}
	b.WriteByte('}')
	return b.Bytes()
}

// ColoredFormat is a log formatter that prints entries using colored text.
func ColoredFormat(e *Entry) []byte {
	var b bytes.Buffer
	b.WriteByte('[')
	b.WriteString(e.Level.Color())
	b.WriteString(e.Level.String())
	b.WriteString(reset)
	b.WriteByte(']')
	if len(e.KeyVals) > 0 {
		b.WriteByte(' ')
		b.WriteByte('[')
		keys, vals := e.KeyVals.Parse()
		for i := 0; i < len(keys); i++ {
			b.WriteString(blue)
			b.WriteString(keys[i])
			b.WriteString(reset)
			b.WriteByte('=')
			b.WriteString(green)
			b.WriteString(fmt.Sprintf("%v", vals[i]))
			b.WriteString(reset)
			if i < len(keys)-1 {
				b.WriteByte(' ')
			}
		}
		b.WriteByte(']')
	}
	b.WriteByte(' ')
	b.WriteString(e.Level.Color())
	b.WriteString(e.Message)
	b.WriteString(reset)
	b.WriteByte('\n')
	return b.Bytes()
}

func writeJSON(val interface{}, b *bytes.Buffer) {
	switch v := val.(type) {
	case nil:
		b.WriteString("null")
	case string:
		b.WriteByte('"')
		b.WriteString(v)
		b.WriteByte('"')
	case int, int32, int64, uint, uint32, uint64:
		b.WriteString(fmt.Sprintf("%d", v))
	case float32:
		b.WriteString(strconv.FormatFloat(float64(v), 'g', -1, 64))
	case float64:
		b.WriteString(strconv.FormatFloat(v, 'g', -1, 64))
	case bool:
		b.WriteString(fmt.Sprintf("%t", v))
	case []string:
		b.WriteByte('[')
		for j := 0; j < len(v); j++ {
			b.WriteByte('"')
			b.WriteString(v[j])
			b.WriteByte('"')
			if j < len(v)-1 {
				b.WriteByte(',')
			}
		}
		b.WriteByte(']')
	case []int:
		b.WriteByte('[')
		for j := 0; j < len(v); j++ {
			b.WriteString(fmt.Sprintf("%d", v[j]))
			if j < len(v)-1 {
				b.WriteByte(',')
			}
		}
		b.WriteByte(']')
	case []int32:
		b.WriteByte('[')
		for j := 0; j < len(v); j++ {
			b.WriteString(fmt.Sprintf("%d", v[j]))
			if j < len(v)-1 {
				b.WriteByte(',')
			}
		}
		b.WriteByte(']')
	case []int64:
		b.WriteByte('[')
		for j := 0; j < len(v); j++ {
			b.WriteString(fmt.Sprintf("%d", v[j]))
			if j < len(v)-1 {
				b.WriteByte(',')
			}
		}
		b.WriteByte(']')
	case []uint:
		b.WriteByte('[')
		for j := 0; j < len(v); j++ {
			b.WriteString(fmt.Sprintf("%d", v[j]))
			if j < len(v)-1 {
				b.WriteByte(',')
			}
		}
		b.WriteByte(']')
	case []uint32:
		b.WriteByte('[')
		for j := 0; j < len(v); j++ {
			b.WriteString(fmt.Sprintf("%d", v[j]))
			if j < len(v)-1 {
				b.WriteByte(',')
			}
		}
		b.WriteByte(']')
	case []uint64:
		b.WriteByte('[')
		for j := 0; j < len(v); j++ {
			b.WriteString(fmt.Sprintf("%d", v[j]))
			if j < len(v)-1 {
				b.WriteByte(',')
			}
		}
		b.WriteByte(']')
	case []float32:
		b.WriteByte('[')
		for j := 0; j < len(v); j++ {
			b.WriteString(strconv.FormatFloat(float64(v[j]), 'g', -1, 64))
			if j < len(v)-1 {
				b.WriteByte(',')
			}
		}
		b.WriteByte(']')
	case []float64:
		b.WriteByte('[')
		for j := 0; j < len(v); j++ {
			b.WriteString(strconv.FormatFloat(v[j], 'g', -1, 64))
			if j < len(v)-1 {
				b.WriteByte(',')
			}
		}
		b.WriteByte(']')
	case []bool:
		b.WriteByte('[')
		for j := 0; j < len(v); j++ {
			b.WriteString(fmt.Sprintf("%t", v[j]))
			if j < len(v)-1 {
				b.WriteByte(',')
			}
		}
		b.WriteByte(']')
	case []interface{}:
		b.WriteByte('[')
		for j := 0; j < len(v); j++ {
			writeJSON(v[j], b)
			if j < len(v)-1 {
				b.WriteByte(',')
			}
		}
		b.WriteByte(']')
	default:
		b.WriteString(fmt.Sprintf("%v", v))
	}
}

const (
	reset = "\033[0m"
	green = "\033[32m"
	blue  = "\033[34m"
)
