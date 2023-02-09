package parse

import (
	"go/ast"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/packages"
)

func TestInterface_Name(t *testing.T) {
	cases := []struct {
		Name, Expected string
		TypeSpec       *ast.TypeSpec
	}{
		{
			Name:     "success",
			TypeSpec: &ast.TypeSpec{Name: ast.NewIdent("Doer")},
			Expected: "Doer",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			i := newInterface(nil, "", tc.TypeSpec, nil)
			name := i.Name()
			assert.Equal(t, tc.Expected, name)
		})
	}
}

func TestInterface_IsExported(t *testing.T) {
	cases := []struct {
		Name     string
		TypeSpec *ast.TypeSpec
		Expected bool
	}{
		{
			Name:     "exported",
			TypeSpec: &ast.TypeSpec{Name: ast.NewIdent("Doer")},
			Expected: true,
		},
		{
			Name:     "not exported",
			TypeSpec: &ast.TypeSpec{Name: ast.NewIdent("doer")},
			Expected: false,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			i := newInterface(nil, "", tc.TypeSpec, nil)
			exported := i.IsExported()
			assert.Equal(t, tc.Expected, exported)
		})
	}
}

func TestInterface_File(t *testing.T) {
	cases := []struct {
		Name, File, Expected string
	}{
		{
			Name:     "exported",
			File:     "doer.go",
			Expected: "doer.go",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			i := newInterface(nil, tc.File, nil, nil)
			file := i.File()
			assert.Equal(t, tc.Expected, file)
		})
	}
}

func TestInterface_Methods(t *testing.T) {
	ps, err := packages.Load(&packages.Config{
		Dir:  filepath.Join(".", "_tests"),
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedImports | packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo | packages.NeedModule,
	}, "./doer")
	require.NoError(t, err)
	require.Len(t, ps, 1)
	p := ps[0]

	cases := []struct {
		Name          string
		TypeSpec      *ast.TypeSpec
		InterfaceType *ast.InterfaceType
		Expected      []Method
	}{
		{
			Name:     "success",
			TypeSpec: &ast.TypeSpec{Name: ast.NewIdent("Doer")},
			InterfaceType: &ast.InterfaceType{Methods: &ast.FieldList{List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent("Do")},
					Type: &ast.FuncType{
						Params: &ast.FieldList{List: []*ast.Field{
							{Names: []*ast.Ident{ast.NewIdent("a"), ast.NewIdent("b")}, Type: ast.NewIdent("int")},
							{Names: []*ast.Ident{ast.NewIdent("c")}, Type: ast.NewIdent("float64")},
						}},
						Results: &ast.FieldList{List: []*ast.Field{
							{Names: []*ast.Ident{ast.NewIdent("d"), ast.NewIdent("e")}, Type: ast.NewIdent("int")},
							{Names: []*ast.Ident{ast.NewIdent("err")}, Type: ast.NewIdent("error")},
						}},
					},
				},
			}}},
			Expected: []Method{
				newMethod(p, ast.NewIdent("Do"), &ast.FuncType{
					Params: &ast.FieldList{List: []*ast.Field{
						{Names: []*ast.Ident{ast.NewIdent("a"), ast.NewIdent("b")}, Type: ast.NewIdent("int")},
						{Names: []*ast.Ident{ast.NewIdent("c")}, Type: ast.NewIdent("float64")},
					}},
					Results: &ast.FieldList{List: []*ast.Field{
						{Names: []*ast.Ident{ast.NewIdent("d"), ast.NewIdent("e")}, Type: ast.NewIdent("int")},
						{Names: []*ast.Ident{ast.NewIdent("err")}, Type: ast.NewIdent("error")},
					}},
				}, false),
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			i := newInterface(p, "", tc.TypeSpec, tc.InterfaceType)
			methods := i.Methods()
			assert.Equal(t, tc.Expected, methods)
		})
	}
}
