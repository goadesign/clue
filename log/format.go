package log

import (
	"bytes"
	"encoding/json"
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

// TimestampFormatLayout is used to set the layout for our TimestampKey
// (default "time") values.  Default format is time.RFC3339.
var TimestampFormatLayout = time.RFC3339

// FormatText is the default log formatter when not running in a terminal, it
// prints entries using the logfmt format:
//
//	time=TIME level=SEVERITY KEY=VAL KEY=VAL ...
//
// Where TIME is the UTC timestamp in RFC3339 format, SEVERITY is one of
// "debug", "info" or "error", and KEY=VAL are the entry key/value pairs.
// Values are quoted and escaped according to the logfmt specification.
//
// Output can be customised with log.TimestampKey, log.TimestampFormatLayout,
// and log.SeverityKey.
func FormatText(e *Entry) []byte {
	var b bytes.Buffer
	enc := logfmt.NewEncoder(&b)
	enc.EncodeKeyval(TimestampKey, e.Time.Format(TimestampFormatLayout)) // nolint: errcheck
	enc.EncodeKeyval(SeverityKey, e.Severity)                            // nolint: errcheck
	for _, kv := range e.KeyVals {
		// Make logfmt format slices
		v := kv.V
		switch kv.V.(type) {
		case []int, []int32, []int64, []uint, []uint32, []uint64, []float32, []float64, []string, []bool, []interface{}:
			var buf bytes.Buffer
			writeJSON(kv.V, &buf)
			v = buf.String()
		}
		enc.EncodeKeyval(kv.K, v) // nolint: errcheck
	}
	b.WriteByte('\n')
	return b.Bytes()
}

// FormatJSON is a log formatter that prints entries using JSON. Entries are
// formatted as follows:
//
//	{
//	  "time": "TIMESTAMP", // UTC timestamp in RFC3339 format
//	  "level": "SEVERITY", // one of DEBUG, INFO or ERROR
//	  "key1": "val1",      // entry key/value pairs
//	  "key2": "val2",
//	  ...
//	}
//
// note: the implementation avoids using reflection (and thus the json package)
// for efficiency.
//
// Output can be customised with log.TimestampKey, log.TimestampFormatLayout,
// and log.SeverityKey.
func FormatJSON(e *Entry) []byte {
	var b bytes.Buffer
	b.WriteString(`{"`)
	b.WriteString(TimestampKey)
	b.WriteString(`":"`)
	b.WriteString(e.Time.Format(TimestampFormatLayout))
	b.WriteString(`","`)
	b.WriteString(SeverityKey)
	b.WriteString(`":"`)
	b.WriteString(e.Severity.String())
	b.WriteByte('"')
	if len(e.KeyVals) > 0 {
		b.WriteByte(',')
		for i, kv := range e.KeyVals {
			b.WriteString(`"`)
			b.WriteString(kv.K)
			b.WriteString(`":`)
			writeJSON(kv.V, &b)
			if i < len(e.KeyVals)-1 {
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
//	SEVERITY[seconds] key=val key=val ...
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
	if len(e.KeyVals) > 0 {
		b.WriteByte(' ')
		for i, kv := range e.KeyVals {
			b.WriteString(e.Severity.Color())
			b.WriteString(kv.K)
			b.WriteString(reset)
			b.WriteString(fmt.Sprintf("=%v", kv.V))
			if i < len(e.KeyVals)-1 {
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
		res, _ := json.Marshal(v)
		b.Write(res)
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
			res, _ := json.Marshal(v[j])
			b.Write(res)
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
		res, _ := json.Marshal(fmt.Sprintf("%v", v))
		b.Write(res)
	}
}
