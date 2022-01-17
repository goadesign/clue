package instrument

import (
	"bytes"
	"context"
	"testing"

	"github.com/crossnokaye/micro/log"
)

func TestLogger(t *testing.T) {
	data := struct {
		foo string
		bar int
	}{"foo", 1}
	cases := []struct {
		name     string
		vals     []interface{}
		expected string
	}{
		{"empty", []interface{}{}, "\n"},
		{"one", []interface{}{"one"}, "one\n"},
		{"many", []interface{}{"one", 2, data}, "one 2 {foo 1}\n"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var buf bytes.Buffer
			ctx := log.Context(context.Background(),
				log.WithOutput(&buf),
				log.WithFormat(func(e *log.Entry) []byte { return []byte(e.Message) }))
			l := logger{ctx}
			l.Println(c.vals...)
			if buf.String() != c.expected {
				t.Errorf("expected %q, got %q", c.expected, buf.String())
			}
		})
	}
}
