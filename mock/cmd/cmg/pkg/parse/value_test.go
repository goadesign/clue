package parse

import (
	"go/ast"
	"go/types"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/packages"
)

type (
	fakeType struct {
		Name string
	}
)

func (t *fakeType) Underlying() types.Type {
	return t
}

func (t *fakeType) String() string {
	return t.Name
}

func TestValue_Name(t *testing.T) {
	cases := []struct {
		Name, Expected string
		Ident          *ast.Ident
	}{
		{
			Name:     "empty",
			Ident:    nil,
			Expected: "",
		},
		{
			Name:     "success",
			Ident:    ast.NewIdent("a"),
			Expected: "a",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			value := newValue(nil, tc.Ident, nil)
			name := value.Name()
			assert.Equal(t, tc.Expected, name)
		})
	}
}

func TestValue_Type(t *testing.T) {
	var (
		stringIdent            = ast.NewIdent("string")
		stringType  types.Type = &fakeType{"string"}
	)

	cases := []struct {
		Name      string
		Package   *packages.Package
		ValueType ast.Expr
		Expected  types.Type
	}{
		{
			Name: "success",
			Package: &packages.Package{TypesInfo: &types.Info{Types: map[ast.Expr]types.TypeAndValue{
				stringIdent: {Type: stringType},
			}}},
			ValueType: stringIdent,
			Expected:  stringType,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			value := newValue(tc.Package, nil, tc.ValueType)
			typ := value.Type()
			assert.Equal(t, tc.Expected, typ)
		})
	}
}
