package log

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"
)

func TestFormat(t *testing.T) {
	now, epoc := timeNow, epoch
	timeNow = func() time.Time { return time.Date(2022, time.January, 9, 20, 29, 45, 0, time.UTC) }
	epoch = timeNow()
	defer func() { timeNow = now; epoch = epoc }()

	keyVals := []KV{
		{"string", "val"},
		{"stringWithQuotes", `example "val"`},
		{"int", 1},
		{"int32", int32(2)},
		{"int64", int64(3)},
		{"uint", uint(4)},
		{"uint32", uint32(5)},
		{"uint64", uint64(6)},
		{"float32", float32(7)},
		{"float64", float64(8.1)},
		{"bool", true},
		{"nil", nil},
		{"dur", 123 * time.Millisecond},
		{"sliceString", []string{"a", "b", "c"}},
		{"sliceStringWithQuotes", []string{"example \"a\""}},
		{"sliceInt", []int{1, 1}},
		{"sliceInt32", []int32{2, 2}},
		{"sliceInt64", []int64{3, 3}},
		{"sliceUint", []uint{4, 4}},
		{"sliceUint32", []uint32{5, 5}},
		{"sliceUint64", []uint64{6, 6}},
		{"sliceFloat32", []float32{7, 7}},
		{"sliceFloat64", []float64{8.1, 8.1}},
		{"sliceBool", []bool{true, false, true}},
		{"sliceNil", []interface{}{nil, nil, nil}},
		{"sliceMix", []interface{}{"a", 1, true, nil}},
	}

	formattedKeyVals := "string=val " +
		`stringWithQuotes="example \"val\"" ` +
		"int=1 " +
		"int32=2 " +
		"int64=3 " +
		"uint=4 " +
		"uint32=5 " +
		"uint64=6 " +
		"float32=7 " +
		"float64=8.1 " +
		"bool=true " +
		"nil=null " +
		"dur=123ms " +
		`sliceString="[\"a\",\"b\",\"c\"]" ` +
		`sliceStringWithQuotes="[\"example \\\"a\\\"\"]" ` +
		"sliceInt=[1,1] " +
		"sliceInt32=[2,2] " +
		"sliceInt64=[3,3] " +
		"sliceUint=[4,4] " +
		"sliceUint32=[5,5] " +
		"sliceUint64=[6,6] " +
		"sliceFloat32=[7,7] " +
		"sliceFloat64=[8.1,8.1] " +
		"sliceBool=[true,false,true] " +
		"sliceNil=[null,null,null] " +
		`sliceMix="[\"a\",1,true,null]"`

	coloredKeyVals := func(col string) string {
		return col + "string\033[0m=val " +
			col + "stringWithQuotes\033[0m=example \"val\" " +
			col + "int\033[0m=1 " +
			col + "int32\033[0m=2 " +
			col + "int64\033[0m=3 " +
			col + "uint\033[0m=4 " +
			col + "uint32\033[0m=5 " +
			col + "uint64\033[0m=6 " +
			col + "float32\033[0m=7 " +
			col + "float64\033[0m=8.1 " +
			col + "bool\033[0m=true " +
			col + "nil\033[0m=<nil> " +
			col + "dur\033[0m=123ms " +
			col + "sliceString\033[0m=[a b c] " +
			col + "sliceStringWithQuotes\033[0m=[example \"a\"] " +
			col + "sliceInt\033[0m=[1 1] " +
			col + "sliceInt32\033[0m=[2 2] " +
			col + "sliceInt64\033[0m=[3 3] " +
			col + "sliceUint\033[0m=[4 4] " +
			col + "sliceUint32\033[0m=[5 5] " +
			col + "sliceUint64\033[0m=[6 6] " +
			col + "sliceFloat32\033[0m=[7 7] " +
			col + "sliceFloat64\033[0m=[8.1 8.1] " +
			col + "sliceBool\033[0m=[true false true] " +
			col + "sliceNil\033[0m=[<nil> <nil> <nil>] " +
			col + "sliceMix\033[0m=[a 1 true <nil>]"
	}

	jsonKeyVals := `"string":"val",` +
		`"stringWithQuotes":"example \"val\"",` +
		`"int":1,` +
		`"int32":2,` +
		`"int64":3,` +
		`"uint":4,` +
		`"uint32":5,` +
		`"uint64":6,` +
		`"float32":7,` +
		`"float64":8.1,` +
		`"bool":true,` +
		`"nil":null,` +
		`"dur":"123ms",` +
		`"sliceString":["a","b","c"],` +
		`"sliceStringWithQuotes":["example \"a\""],` +
		`"sliceInt":[1,1],` +
		`"sliceInt32":[2,2],` +
		`"sliceInt64":[3,3],` +
		`"sliceUint":[4,4],` +
		`"sliceUint32":[5,5],` +
		`"sliceUint64":[6,6],` +
		`"sliceFloat32":[7,7],` +
		`"sliceFloat64":[8.1,8.1],` +
		`"sliceBool":[true,false,true],` +
		`"sliceNil":[null,null,null],` +
		`"sliceMix":["a",1,true,null]`

	cases := []struct {
		name    string
		logfn   func(ctx context.Context, keyvals ...Fielder)
		format  FormatFunc
		keyVals []KV
		want    string
	}{
		{
			name:    "default debug",
			logfn:   Debug,
			format:  FormatText,
			keyVals: keyVals,
			want:    "time=2022-01-09T20:29:45Z level=debug " + formattedKeyVals + "\n",
		},
		{
			name:    "default info",
			logfn:   Info,
			format:  FormatText,
			keyVals: keyVals,
			want:    "time=2022-01-09T20:29:45Z level=info " + formattedKeyVals + "\n",
		},
		{
			name:    "default print",
			logfn:   Print,
			format:  FormatText,
			keyVals: keyVals,
			want:    "time=2022-01-09T20:29:45Z level=info " + formattedKeyVals + "\n",
		},
		{
			name:    "colored debug",
			logfn:   Debug,
			format:  FormatTerminal,
			keyVals: keyVals,
			want:    "\033[37mDEBG\033[0m[0000] " + coloredKeyVals("\033[37m") + "\n",
		},
		{
			name:    "colored info",
			logfn:   Info,
			format:  FormatTerminal,
			keyVals: keyVals,
			want:    "\033[34mINFO\033[0m[0000] " + coloredKeyVals("\033[34m") + "\n",
		},
		{
			name:    "colored print",
			logfn:   Print,
			format:  FormatTerminal,
			keyVals: keyVals,
			want:    "\033[34mINFO\033[0m[0000] " + coloredKeyVals("\033[34m") + "\n",
		},
		{
			name:    "json debug",
			logfn:   Debug,
			format:  FormatJSON,
			keyVals: keyVals,
			want:    `{"time":"2022-01-09T20:29:45Z","level":"debug",` + jsonKeyVals + "}\n",
		},
		{
			name:    "json info",
			logfn:   Info,
			format:  FormatJSON,
			keyVals: keyVals,
			want:    `{"time":"2022-01-09T20:29:45Z","level":"info",` + jsonKeyVals + "}\n",
		},
		{
			name:    "json print",
			logfn:   Print,
			format:  FormatJSON,
			keyVals: keyVals,
			want:    `{"time":"2022-01-09T20:29:45Z","level":"info",` + jsonKeyVals + "}\n",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			ctx := Context(context.Background(), WithOutput(&buf), WithFormat(tc.format), WithDebug())
			tc.logfn(ctx, kvList(tc.keyVals))
			got := buf.String()
			if got != tc.want {
				t.Errorf("got %s, want %s", got, tc.want)
			}
		})
	}

	errorCases := []struct {
		name    string
		logfn   func(ctx context.Context, err error, keyvals ...Fielder)
		format  FormatFunc
		keyVals []KV
		want    string
	}{
		{
			name:    "default error",
			logfn:   Error,
			format:  FormatText,
			keyVals: keyVals,
			want:    "time=2022-01-09T20:29:45Z level=error " + formattedKeyVals + " err=error\n",
		},
		{
			name:    "colored error",
			logfn:   Error,
			format:  FormatTerminal,
			keyVals: keyVals,
			want:    "\033[1;31mERRO\033[0m[0000] " + coloredKeyVals("\033[1;31m") + " \033[1;31merr\033[0m=error\n",
		},
		{
			name:    "json info",
			logfn:   Error,
			format:  FormatJSON,
			keyVals: keyVals,
			want:    `{"time":"2022-01-09T20:29:45Z","level":"error",` + jsonKeyVals + `,"err":"error"}` + "\n",
		},
	}
	for _, tc := range errorCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			ctx := Context(context.Background(), WithOutput(&buf), WithFormat(tc.format), WithDebug())
			tc.logfn(ctx, errors.New("error"), kvList(tc.keyVals))
			got := buf.String()
			if got != tc.want {
				t.Errorf("got %s, want %s", got, tc.want)
			}
		})
	}
}
