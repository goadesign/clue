package log

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	"github.com/go-logfmt/logfmt"
)

// reset escape sequence for color unset
const reset = "\033[0m"

// epoch used to compute relative timestamps
var epoch time.Time

func init() {
	epoch = time.Now()
}

// FormatText is the default log formatter when not running in a terminal, it
// prints entries using the logfmt format:
//
//    time=TIME level=SEVERITY msg=MESSAGE KEY=VAL KEY=VAL ...
//
// Where TIME is the UTC timestamp in RFC3339 format, SEVERITY is one of
// "debug", "info" or "error", MESSAGE is the log message, and KEY=VAL are the
// entry key/value pairs.
func FormatText(e *Entry) []byte {
	kvs := []interface{}{"time", e.Time.Format(time.RFC3339), "level", e.Severity}
	if len(e.Message) > 0 {
		kvs = append(kvs, "msg", e.Message)
	}
	kvs = append(kvs, e.KeyVals...)
	var b bytes.Buffer
	logfmt.NewEncoder(&b).EncodeKeyvals(kvs...)
	b.WriteByte('\n')
	return b.Bytes()
}

// FormatJSON is a log formatter that prints entries using JSON. Entries are
// formatted as follows:
//
//   {
//     "time": "TIMESTAMP", // UTC timestamp in RFC3339 format
//     "level": "SEVERITY", // one of DEBUG, INFO or ERROR
//     "msg": "MESSAGE",    // log message
//     "key1": "val1",      // entry key/value pairs
//     "key2": "val2",
//     ...
//   }
//
// note: the implementation avoids using reflection (and thus the json package)
// for efficiency.
func FormatJSON(e *Entry) []byte {
	var b bytes.Buffer
	b.WriteString(`{"time":"`)
	b.WriteString(e.Time.Format(time.RFC3339))
	b.WriteString(`","level":"`)
	b.WriteString(e.Severity.String())
	b.WriteByte('"')
	if len(e.Message) > 0 {
		b.WriteString(`,"msg":"`)
		b.WriteString(e.Message)
		b.WriteString(`"`)
	}
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
	b.WriteString("}\n")
	return b.Bytes()
}

// FormatTerminal is a log formatter that prints entries suitable for terminal
// that supports colors. It prints entries in the following format:
//
//    SEVERITY[seconds] message key=val key=val ...
//
// Where SEVERITY is one of DEBG, INFO or ERRO, seconds is the number of seconds
// since the application started, message is the log message, and key=val are
// the entry key/value pairs. The severity and keys are colored according to the
// severity (gray for debug entries, blue for info entries and red for errors).
func FormatTerminal(e *Entry) []byte {
	var b bytes.Buffer
	b.WriteString(e.Severity.Color())
	b.WriteString(e.Severity.Code())
	b.WriteString(reset)
	b.WriteString(fmt.Sprintf("[%04d]", int(e.Time.Sub(epoch)/time.Second)))
	if len(e.Message) > 0 {
		b.WriteByte(' ')
		b.WriteString(e.Message)
	}
	if len(e.KeyVals) > 0 {
		b.WriteByte(' ')
		keys, vals := e.KeyVals.Parse()
		for i := 0; i < len(keys); i++ {
			b.WriteString(e.Severity.Color())
			b.WriteString(keys[i])
			b.WriteString(reset)
			b.WriteString(fmt.Sprintf("=%v", vals[i]))
			if i < len(keys)-1 {
				b.WriteByte(' ')
			}
		}
	}
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
		b.WriteByte('"')
		b.WriteString(fmt.Sprintf("%v", v))
		b.WriteByte('"')
	}
}
