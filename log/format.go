package log

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"unicode/utf8"

	"google.golang.org/grpc/codes"
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
	b := make([]byte, 0, 256)

	b = appendKeyValue(b, TimestampKey, e.Time.Format(TimestampFormatLayout))
	b = append(b, ' ')
	b = appendKeyValue(b, SeverityKey, e.Severity)

	for _, kv := range e.KeyVals {
		b = append(b, ' ')
		b = appendKeyValue(b, kv.K, kv.V)
	}

	b = append(b, '\n')
	return b
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
	b := make([]byte, 0, 256)

	b = append(b, '{')
	b = appendJSONKeyValue(b, TimestampKey, e.Time.Format(TimestampFormatLayout))
	b = append(b, ',')
	b = appendJSONKeyValue(b, SeverityKey, e.Severity.String())

	for _, kv := range e.KeyVals {
		b = append(b, ',')
		b = appendJSONKeyValue(b, kv.K, kv.V)
	}

	b = append(b, "}\n"...)
	return b
}

func appendKeyValue(b []byte, key string, value any) []byte {
	b = append(b, key...)
	b = append(b, '=')
	return appendTextValue(b, value)
}

func appendTextValue(b []byte, value any) []byte {
	switch v := value.(type) {
	case string:
		return appendEscapedString(b, v)
	case int, int32, int64, uint, uint32, uint64, float32, float64, bool:
		return appendJSONValue(b, v)
	case []any:
		return appendTextArray(b, v)
	default:
		// Fallback to fmt.Sprintf for complex types
		return append(b, fmt.Sprintf("%v", v)...)
	}
}

func appendEscapedString(b []byte, s string) []byte {
	if needsQuoting(s) {
		b = append(b, '"')
		for i := 0; i < len(s); i++ {
			if s[i] == '"' || s[i] == '\\' {
				b = append(b, '\\')
			}
			b = append(b, s[i])
		}
		b = append(b, '"')
	} else {
		b = append(b, s...)
	}
	return b
}

func needsQuoting(s string) bool {
	if len(s) == 0 {
		return true
	}
	for i := 0; i < len(s); i++ {
		if s[i] <= ' ' || s[i] == '=' || s[i] == '"' || s[i] == '\\' {
			return true
		}
	}
	return false
}

func appendTextArray(b []byte, arr []any) []byte {
	b = append(b, '[')
	for i, v := range arr {
		if i > 0 {
			b = append(b, ' ')
		}
		b = appendTextValue(b, v)
	}
	return append(b, ']')
}

func appendJSONKeyValue(b []byte, key string, value any) []byte {
	b = append(b, '"')
	b = append(b, key...)
	b = append(b, `":`...)
	return appendJSONValue(b, value)
}

func appendJSONValue(b []byte, value any) []byte {
	switch v := value.(type) {
	case string:
		return appendJSONString(b, v)
	case int:
		return strconv.AppendInt(b, int64(v), 10)
	case int32:
		return strconv.AppendInt(b, int64(v), 10)
	case int64:
		return strconv.AppendInt(b, v, 10)
	case uint:
		return strconv.AppendUint(b, uint64(v), 10)
	case uint32:
		return strconv.AppendUint(b, uint64(v), 10)
	case uint64:
		return strconv.AppendUint(b, v, 10)
	case float32:
		return strconv.AppendFloat(b, float64(v), 'g', -1, 32)
	case float64:
		return strconv.AppendFloat(b, v, 'g', -1, 64)
	case bool:
		return strconv.AppendBool(b, v)
	case []any:
		return appendJSONArray(b, v)
	case time.Duration:
		return appendJSONString(b, v.String())
	case codes.Code:
		return appendJSONString(b, v.String())
	default:
		jsonValue, _ := json.Marshal(v)
		return append(b, jsonValue...)
	}
}

func appendJSONString(b []byte, s string) []byte {
	b = append(b, '"')
	for i := 0; i < len(s); i++ {
		if s[i] < utf8.RuneSelf {
			switch s[i] {
			case '"', '\\':
				b = append(b, '\\', s[i])
			case '\b':
				b = append(b, '\\', 'b')
			case '\f':
				b = append(b, '\\', 'f')
			case '\n':
				b = append(b, '\\', 'n')
			case '\r':
				b = append(b, '\\', 'r')
			case '\t':
				b = append(b, '\\', 't')
			default:
				b = append(b, s[i])
			}
		} else {
			b = append(b, s[i])
		}
	}
	return append(b, '"')
}

func appendJSONArray(b []byte, arr []any) []byte {
	b = append(b, '[')
	for i, v := range arr {
		if i > 0 {
			b = append(b, ',')
		}
		b = appendJSONValue(b, v)
	}
	return append(b, ']')
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
	b := make([]byte, 0, 256)
	b = append(b, e.Severity.Color()...)
	b = append(b, e.Severity.Code()...)
	b = append(b, reset...)
	b = fmt.Appendf(b, "[%04d]", int(e.Time.Sub(epoch)/time.Second))
	if len(e.KeyVals) > 0 {
		b = append(b, ' ')
		for i, kv := range e.KeyVals {
			b = append(b, e.Severity.Color()...)
			b = append(b, kv.K...)
			b = append(b, reset...)
			b = fmt.Appendf(b, "=%v", kv.V)
			if i < len(e.KeyVals)-1 {
				b = append(b, ' ')
			}
		}
	}
	b = append(b, '\n')
	return b
}
