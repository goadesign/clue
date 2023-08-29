package parse

import (
	"go/ast"
	"go/types"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/packages"
)

func TestMethod_Name(t *testing.T) {
	cases := []struct {
		Name, Expected string
		Method         Method
	}{
		{
			Name:     "AST success",
			Method:   newASTMethod(nil, ast.NewIdent("Do"), nil, false),
			Expected: "Do",
		},
		{
			Name:     "types success",
			Method:   newTypesMethod(types.NewFunc(0, nil, "Do", &types.Signature{})),
			Expected: "Do",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			name := tc.Method.Name()
			assert.Equal(t, tc.Expected, name)
		})
	}
}

func TestMethod_IsExported(t *testing.T) {
	cases := []struct {
		Name     string
		Method   Method
		Expected bool
	}{
		{
			Name:     "AST exported",
			Method:   newASTMethod(nil, ast.NewIdent("Do"), nil, false),
			Expected: true,
		},
		{
			Name:     "AST not exported",
			Method:   newASTMethod(nil, ast.NewIdent("do"), nil, false),
			Expected: false,
		},
		{
			Name:     "types exported",
			Method:   newTypesMethod(types.NewFunc(0, nil, "Do", &types.Signature{})),
			Expected: true,
		},
		{
			Name:     "types not exported",
			Method:   newTypesMethod(types.NewFunc(0, nil, "do", &types.Signature{})),
			Expected: false,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			exported := tc.Method.IsExported()
			assert.Equal(t, tc.Expected, exported)
		})
	}
}

func TestMethod_Parameters(t *testing.T) {
	var (
		intType     types.Type = &fakeType{"int"}
		float64Type types.Type = &fakeType{"float64"}
	)

	p := &packages.Package{}
	cases := []struct {
		Name     string
		Method   Method
		Expected []Value
	}{
		{
			Name: "AST success",
			Method: newASTMethod(p, nil, &ast.FuncType{Params: &ast.FieldList{List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent("a"), ast.NewIdent("b")},
					Type:  ast.NewIdent("int"),
				},
				{
					Names: []*ast.Ident{ast.NewIdent("c")},
					Type:  ast.NewIdent("float64"),
				},
			}}}, false),
			Expected: []Value{
				newASTValue(p, ast.NewIdent("a"), ast.NewIdent("int")),
				newASTValue(p, ast.NewIdent("b"), ast.NewIdent("int")),
				newASTValue(p, ast.NewIdent("c"), ast.NewIdent("float64")),
			},
		},
		{
			Name: "types success",
			Method: newTypesMethod(types.NewFunc(0, nil, "", types.NewSignatureType(nil, nil, nil, types.NewTuple(
				types.NewParam(0, nil, "a", intType),
				types.NewParam(0, nil, "b", intType),
				types.NewParam(0, nil, "c", float64Type),
			), nil, false))),
			Expected: []Value{
				newTypesValue(types.NewParam(0, nil, "a", intType)),
				newTypesValue(types.NewParam(0, nil, "b", intType)),
				newTypesValue(types.NewParam(0, nil, "c", float64Type)),
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			parameters := tc.Method.Parameters()
			assert.Equal(t, tc.Expected, parameters)
		})
	}
}

func TestMethod_Results(t *testing.T) {
	var (
		intType   types.Type = &fakeType{"int"}
		errorType types.Type = &fakeType{"error"}
	)

	p := &packages.Package{}
	cases := []struct {
		Name     string
		Method   Method
		Expected []Value
	}{
		{
			Name: "AST success",
			Method: newASTMethod(p, nil, &ast.FuncType{Results: &ast.FieldList{List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent("a"), ast.NewIdent("b")},
					Type:  ast.NewIdent("int"),
				},
				{
					Names: []*ast.Ident{ast.NewIdent("err")},
					Type:  ast.NewIdent("error"),
				},
			}}}, false),
			Expected: []Value{
				newASTValue(p, ast.NewIdent("a"), ast.NewIdent("int")),
				newASTValue(p, ast.NewIdent("b"), ast.NewIdent("int")),
				newASTValue(p, ast.NewIdent("err"), ast.NewIdent("error")),
			},
		},
		{
			Name: "types success",
			Method: newTypesMethod(types.NewFunc(0, nil, "", types.NewSignatureType(nil, nil, nil, nil, types.NewTuple(
				types.NewVar(0, nil, "a", intType),
				types.NewVar(0, nil, "b", intType),
				types.NewVar(0, nil, "err", errorType),
			), false))),
			Expected: []Value{
				newTypesValue(types.NewVar(0, nil, "a", intType)),
				newTypesValue(types.NewVar(0, nil, "b", intType)),
				newTypesValue(types.NewVar(0, nil, "err", errorType)),
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			results := tc.Method.Results()
			assert.Equal(t, tc.Expected, results)
		})
	}
}

func TestMethod_Variadic(t *testing.T) {
	var intType types.Type = &fakeType{"int"}

	cases := []struct {
		Name     string
		Method   Method
		Expected bool
	}{
		{
			Name:     "AST variadic",
			Method:   newASTMethod(nil, nil, nil, true),
			Expected: true,
		},
		{
			Name:     "AST not variadic",
			Method:   newASTMethod(nil, nil, nil, false),
			Expected: false,
		},
		{
			Name: "types variadic",
			Method: newTypesMethod(types.NewFunc(0, nil, "", types.NewSignatureType(nil, nil, nil, types.NewTuple(
				types.NewParam(0, nil, "a", types.NewSlice(intType)),
			), nil, true))),
			Expected: true,
		},
		{
			Name: "types not variadic",
			Method: newTypesMethod(types.NewFunc(0, nil, "", types.NewSignatureType(nil, nil, nil, types.NewTuple(
				types.NewParam(0, nil, "a", types.NewSlice(intType)),
			), nil, false))),
			Expected: false,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			variadic := tc.Method.Variadic()
			assert.Equal(t, tc.Expected, variadic)
		})
	}
}
