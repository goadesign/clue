package parse

import (
	"go/ast"
	"go/types"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/packages"
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
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			typ := newType(nil, tc.Ident, nil)
			name := typ.Name()
			assert.Equal(t, tc.Expected, name)
		})
	}
}

func TestType_Constraint(t *testing.T) {
	var (
		comparableIdent            = ast.NewIdent("comparable")
		comparableType  types.Type = &fakeType{"comparable"}
	)

	cases := []struct {
		Name     string
		Package  *packages.Package
		TypeType ast.Expr
		Expected types.Type
	}{
		{
			Name: "success",
			Package: &packages.Package{TypesInfo: &types.Info{Types: map[ast.Expr]types.TypeAndValue{
				comparableIdent: {Type: comparableType},
			}}},
			TypeType: comparableIdent,
			Expected: comparableType,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			value := newType(tc.Package, nil, tc.TypeType)
			constraint := value.Constraint()
			assert.Equal(t, tc.Expected, constraint)
		})
	}
}
