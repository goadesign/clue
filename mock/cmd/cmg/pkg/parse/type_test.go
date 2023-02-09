package parse

import (
	"go/ast"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestType_Name(t *testing.T) {
	cases := []struct {
		Name, Expected string
		Ident          *ast.Ident
	}{
		{
			Name:     "success",
			Ident:    ast.NewIdent("T"),
			Expected: "T",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			typ := newType(nil, tc.Ident, nil)
			name := typ.Name()
			assert.Equal(t, tc.Expected, name)
		})
	}
}
