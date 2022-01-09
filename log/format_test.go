package log

import (
	"bytes"
	"context"
	"testing"
)

func TestFormat(t *testing.T) {
	keyVals := []interface{}{
		"string", "val",
		"int", 1,
		"int32", int32(2),
		"int64", int64(3),
		"uint", uint(4),
		"uint32", uint32(5),
		"uint64", uint64(6),
		"float32", float32(7),
		"float64", float64(8.1),
		"bool", true,
		"nil", nil,
		"sliceString", []string{"a", "b", "c"},
		"sliceInt", []int{1, 1},
		"sliceInt32", []int32{2, 2},
		"sliceInt64", []int64{3, 3},
		"sliceUint", []uint{4, 4},
		"sliceUint32", []uint32{5, 5},
		"sliceUint64", []uint64{6, 6},
		"sliceFloat32", []float32{7, 7},
		"sliceFloat64", []float64{8.1, 8.1},
		"sliceBool", []bool{true, false, true},
		"sliceNil", []interface{}{nil, nil, nil},
		"sliceMix", []interface{}{"a", 1, true, nil},
	}
	formattedKeyVals := "string=val " +
		"int=1 " +
		"int32=2 " +
		"int64=3 " +
		"uint=4 " +
		"uint32=5 " +
		"uint64=6 " +
		"float32=7 " +
		"float64=8.1 " +
		"bool=true " +
		"nil=<nil> " +
		"sliceString=[a b c] " +
		"sliceInt=[1 1] " +
		"sliceInt32=[2 2] " +
		"sliceInt64=[3 3] " +
		"sliceUint=[4 4] " +
		"sliceUint32=[5 5] " +
		"sliceUint64=[6 6] " +
		"sliceFloat32=[7 7] " +
		"sliceFloat64=[8.1 8.1] " +
		"sliceBool=[true false true] " +
		"sliceNil=[<nil> <nil> <nil>] " +
		"sliceMix=[a 1 true <nil>]"
	coloredKeyVals := "\033[34mstring\033[0m=\033[32mval\033[0m " +
		"\033[34mint\033[0m=\033[32m1\033[0m " +
		"\033[34mint32\033[0m=\033[32m2\033[0m " +
		"\033[34mint64\033[0m=\033[32m3\033[0m " +
		"\033[34muint\033[0m=\033[32m4\033[0m " +
		"\033[34muint32\033[0m=\033[32m5\033[0m " +
		"\033[34muint64\033[0m=\033[32m6\033[0m " +
		"\033[34mfloat32\033[0m=\033[32m7\033[0m " +
		"\033[34mfloat64\033[0m=\033[32m8.1\033[0m " +
		"\033[34mbool\033[0m=\033[32mtrue\033[0m " +
		"\033[34mnil\033[0m=\033[32m<nil>\033[0m " +
		"\033[34msliceString\033[0m=\033[32m[a b c]\033[0m " +
		"\033[34msliceInt\033[0m=\033[32m[1 1]\033[0m " +
		"\033[34msliceInt32\033[0m=\033[32m[2 2]\033[0m " +
		"\033[34msliceInt64\033[0m=\033[32m[3 3]\033[0m " +
		"\033[34msliceUint\033[0m=\033[32m[4 4]\033[0m " +
		"\033[34msliceUint32\033[0m=\033[32m[5 5]\033[0m " +
		"\033[34msliceUint64\033[0m=\033[32m[6 6]\033[0m " +
		"\033[34msliceFloat32\033[0m=\033[32m[7 7]\033[0m " +
		"\033[34msliceFloat64\033[0m=\033[32m[8.1 8.1]\033[0m " +
		"\033[34msliceBool\033[0m=\033[32m[true false true]\033[0m " +
		"\033[34msliceNil\033[0m=\033[32m[<nil> <nil> <nil>]\033[0m " +
		"\033[34msliceMix\033[0m=\033[32m[a 1 true <nil>]\033[0m"
	jsonKeyVals := `"string":"val",` +
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
		`"sliceString":["a","b","c"],` +
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
		logfn   func(ctx context.Context, msg string, keyvals ...interface{})
		format  FormatFunc
		msg     string
		keyVals []interface{}
		want    string
	}{
		{
			name:    "default debug",
			logfn:   Debug,
			format:  DefaultFormat,
			msg:     "hello",
			keyVals: keyVals,
			want:    "[DEBUG] [" + formattedKeyVals + "] hello\n",
		},
		{
			name:    "default info",
			logfn:   Info,
			format:  DefaultFormat,
			msg:     "hello",
			keyVals: keyVals,
			want:    "[INFO] [" + formattedKeyVals + "] hello\n",
		},
		{
			name:    "default print",
			logfn:   Print,
			format:  DefaultFormat,
			msg:     "hello",
			keyVals: keyVals,
			want:    "[INFO] [" + formattedKeyVals + "] hello\n",
		},
		{
			name:    "default error",
			logfn:   Error,
			format:  DefaultFormat,
			msg:     "hello",
			keyVals: keyVals,
			want:    "[ERROR] [" + formattedKeyVals + "] hello\n",
		},
		{
			name:    "colored debug",
			logfn:   Debug,
			format:  ColoredFormat,
			msg:     "hello",
			keyVals: keyVals,
			want:    "[\033[1;32mDEBUG\033[0m] [" + coloredKeyVals + "] \033[1;32mhello\033[0m\n",
		},
		{
			name:    "colored info",
			logfn:   Info,
			format:  ColoredFormat,
			msg:     "hello",
			keyVals: keyVals,
			want:    "[\033[1;34mINFO\033[0m] [" + coloredKeyVals + "] \033[1;34mhello\033[0m\n",
		},
		{
			name:    "colored print",
			logfn:   Print,
			format:  ColoredFormat,
			msg:     "hello",
			keyVals: keyVals,
			want:    "[\033[1;34mINFO\033[0m] [" + coloredKeyVals + "] \033[1;34mhello\033[0m\n",
		},
		{
			name:    "colored error",
			logfn:   Error,
			format:  ColoredFormat,
			msg:     "hello",
			keyVals: keyVals,
			want:    "[\033[1;31mERROR\033[0m] [" + coloredKeyVals + "] \033[1;31mhello\033[0m\n",
		},
		{
			name:    "json debug",
			logfn:   Debug,
			format:  JSONFormat,
			msg:     "hello",
			keyVals: keyVals,
			want:    `{"severity":"DEBUG","message":"hello",` + jsonKeyVals + "}",
		},
		{
			name:    "json info",
			logfn:   Info,
			format:  JSONFormat,
			msg:     "hello",
			keyVals: keyVals,
			want:    `{"severity":"INFO","message":"hello",` + jsonKeyVals + "}",
		},
		{
			name:    "json print",
			logfn:   Print,
			format:  JSONFormat,
			msg:     "hello",
			keyVals: keyVals,
			want:    `{"severity":"INFO","message":"hello",` + jsonKeyVals + "}",
		},
		{
			name:    "json error",
			logfn:   Error,
			format:  JSONFormat,
			msg:     "hello",
			keyVals: keyVals,
			want:    `{"severity":"ERROR","message":"hello",` + jsonKeyVals + "}",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			ctx := Context(context.Background(), WithOutput(&buf), WithFormat(tc.format), WithDebug())
			tc.logfn(ctx, tc.msg, tc.keyVals...)
			got := buf.String()
			if got != tc.want {
				t.Errorf("got %s, want %s", got, tc.want)
			}
		})
	}
}
