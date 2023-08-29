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
		Value          Value
	}{
		{
			Name:     "AST empty",
			Value:    newASTValue(nil, nil, nil),
			Expected: "",
		},
		{
			Name:     "AST success",
			Value:    newASTValue(nil, ast.NewIdent("a"), nil),
			Expected: "a",
		},
		{
			Name:     "types success",
			Value:    newTypesValue(types.NewVar(0, nil, "a", nil)),
			Expected: "a",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			name := tc.Value.Name()
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
		Name     string
		Value    Value
		Expected types.Type
	}{
		{
			Name: "AST success",
			Value: newASTValue(&packages.Package{TypesInfo: &types.Info{Types: map[ast.Expr]types.TypeAndValue{
				stringIdent: {Type: stringType},
			}}}, nil, stringIdent),
			Expected: stringType,
		},
		{
			Name:     "types success",
			Value:    newTypesValue(types.NewVar(0, nil, "", stringType)),
			Expected: stringType,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			typ := tc.Value.Type()
			assert.Equal(t, tc.Expected, typ)
		})
	}
}
