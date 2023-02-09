package parse

import (
	"go/ast"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/packages"
)

func TestMethod_Name(t *testing.T) {
	cases := []struct {
		Name, Expected string
		Ident          *ast.Ident
	}{
		{
			Name:     "success",
			Ident:    ast.NewIdent("Do"),
			Expected: "Do",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			method := newMethod(nil, tc.Ident, nil, false)
			name := method.Name()
			assert.Equal(t, tc.Expected, name)
		})
	}
}

func TestMethod_(t *testing.T) {
	p := &packages.Package{}
	cases := []struct {
		Name     string
		FuncType *ast.FuncType
		Expected []Value
	}{
		{
			Name: "success",
			FuncType: &ast.FuncType{Results: &ast.FieldList{List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent("a"), ast.NewIdent("b")},
					Type:  ast.NewIdent("int"),
				},
				{
					Names: []*ast.Ident{ast.NewIdent("err")},
					Type:  ast.NewIdent("error"),
				},
			}}},
			Expected: []Value{
				newValue(p, ast.NewIdent("a"), ast.NewIdent("int")),
				newValue(p, ast.NewIdent("b"), ast.NewIdent("int")),
				newValue(p, ast.NewIdent("err"), ast.NewIdent("error")),
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			method := newMethod(p, nil, tc.FuncType, false)
			results := method.Results()
			assert.Equal(t, tc.Expected, results)
		})
	}
}
func TestMethod_IsExported(t *testing.T) {
	cases := []struct {
		Name     string
		Ident    *ast.Ident
		Expected bool
	}{
		{
			Name:     "exported",
			Ident:    ast.NewIdent("Do"),
			Expected: true,
		},
		{
			Name:     "not exported",
			Ident:    ast.NewIdent("do"),
			Expected: false,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			method := newMethod(nil, tc.Ident, nil, false)
			exported := method.IsExported()
			assert.Equal(t, tc.Expected, exported)
		})
	}
}

func TestMethod_TypeParameters(t *testing.T) {
	p := &packages.Package{}
	cases := []struct {
		Name     string
		FuncType *ast.FuncType
		Expected []Type
	}{
		{
			Name: "success",
			FuncType: &ast.FuncType{TypeParams: &ast.FieldList{List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent("K")},
					Type:  ast.NewIdent("comparable"),
				},
				{
					Names: []*ast.Ident{ast.NewIdent("V")},
					Type:  ast.NewIdent("any"),
				},
			}}},
			Expected: []Type{
				newType(p, ast.NewIdent("K"), ast.NewIdent("comparable")),
				newType(p, ast.NewIdent("V"), ast.NewIdent("any")),
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			method := newMethod(p, nil, tc.FuncType, false)
			parameters := method.TypeParameters()
			assert.Equal(t, tc.Expected, parameters)
		})
	}
}

func TestMethod_Parameters(t *testing.T) {
	p := &packages.Package{}
	cases := []struct {
		Name     string
		FuncType *ast.FuncType
		Expected []Value
	}{
		{
			Name: "success",
			FuncType: &ast.FuncType{Params: &ast.FieldList{List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent("a"), ast.NewIdent("b")},
					Type:  ast.NewIdent("int"),
				},
				{
					Names: []*ast.Ident{ast.NewIdent("c")},
					Type:  ast.NewIdent("float64"),
				},
			}}},
			Expected: []Value{
				newValue(p, ast.NewIdent("a"), ast.NewIdent("int")),
				newValue(p, ast.NewIdent("b"), ast.NewIdent("int")),
				newValue(p, ast.NewIdent("c"), ast.NewIdent("float64")),
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			method := newMethod(p, nil, tc.FuncType, false)
			parameters := method.Parameters()
			assert.Equal(t, tc.Expected, parameters)
		})
	}
}

func TestMethod_Results(t *testing.T) {
	p := &packages.Package{}
	cases := []struct {
		Name     string
		FuncType *ast.FuncType
		Expected []Value
	}{
		{
			Name: "success",
			FuncType: &ast.FuncType{Results: &ast.FieldList{List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent("a"), ast.NewIdent("b")},
					Type:  ast.NewIdent("int"),
				},
				{
					Names: []*ast.Ident{ast.NewIdent("err")},
					Type:  ast.NewIdent("error"),
				},
			}}},
			Expected: []Value{
				newValue(p, ast.NewIdent("a"), ast.NewIdent("int")),
				newValue(p, ast.NewIdent("b"), ast.NewIdent("int")),
				newValue(p, ast.NewIdent("err"), ast.NewIdent("error")),
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			method := newMethod(p, nil, tc.FuncType, false)
			results := method.Results()
			assert.Equal(t, tc.Expected, results)
		})
	}
}

func TestMethod_Variadic(t *testing.T) {
	cases := []struct {
		Name               string
		Variadic, Expected bool
	}{
		{
			Name:     "variadic",
			Variadic: true,
			Expected: true,
		},
		{
			Name:     "not variadic",
			Variadic: false,
			Expected: false,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			method := newMethod(nil, nil, nil, tc.Variadic)
			variadic := method.Variadic()
			assert.Equal(t, tc.Expected, variadic)
		})
	}
}
