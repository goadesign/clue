package parse

import (
	"go/ast"
	"go/types"
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

func TestInterface_TypeParameters(t *testing.T) {
	p := &packages.Package{}
	cases := []struct {
		Name     string
		TypeSpec *ast.TypeSpec
		Expected []Type
	}{
		{
			Name: "success",
			TypeSpec: &ast.TypeSpec{TypeParams: &ast.FieldList{List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent("K")},
					Type:  ast.NewIdent("comparable"),
				},
				{
					Names: []*ast.Ident{ast.NewIdent("V"), ast.NewIdent("X")},
					Type:  ast.NewIdent("any"),
				},
			}}},
			Expected: []Type{
				newType(p, ast.NewIdent("K"), ast.NewIdent("comparable")),
				newType(p, ast.NewIdent("V"), ast.NewIdent("any")),
				newType(p, ast.NewIdent("X"), ast.NewIdent("any")),
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			i := newInterface(p, "", tc.TypeSpec, nil)
			parameters := i.TypeParameters()
			assert.Equal(t, tc.Expected, parameters)
		})
	}
}

func TestInterface_Methods(t *testing.T) {
	ps, err := packages.Load(&packages.Config{
		Dir:  "_tests",
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedImports | packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo | packages.NeedModule,
	}, "./doer")
	require.NoError(t, err)
	require.Len(t, ps, 1)
	p := ps[0]

	var (
		externalDoerSelectorExpr *ast.SelectorExpr
		externalDoerInterface    *types.Interface
	)
	for at, tv := range p.TypesInfo.Types {
		if se, ok := at.(*ast.SelectorExpr); ok {
			if i, ok := se.X.(*ast.Ident); ok {
				if i.Name != "external" {
					continue
				}
				if se.Sel.Name == "Doer" {
					externalDoerSelectorExpr = se
					externalDoerInterface, _ = tv.Type.Underlying().(*types.Interface)
					break
				}
			}
		}
	}
	require.NotNil(t, externalDoerSelectorExpr)
	require.NotNil(t, externalDoerInterface)
	require.Greater(t, externalDoerInterface.NumMethods(), 0)

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
				newASTMethod(p, ast.NewIdent("Do"), &ast.FuncType{
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
		{
			Name:     "embedded interface",
			TypeSpec: &ast.TypeSpec{Name: ast.NewIdent("EmbeddedDoer")},
			InterfaceType: &ast.InterfaceType{Methods: &ast.FieldList{List: []*ast.Field{{
				Type: &ast.Ident{Name: "Doer", Obj: &ast.Object{Kind: ast.Typ, Decl: &ast.TypeSpec{
					Type: &ast.InterfaceType{Methods: &ast.FieldList{List: []*ast.Field{
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
				}}}}},
			}},
			Expected: []Method{
				newASTMethod(p, ast.NewIdent("Do"), &ast.FuncType{
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
		{
			Name:          "external embedded interface",
			TypeSpec:      &ast.TypeSpec{Name: ast.NewIdent("ExternalEmbeddedDoer")},
			InterfaceType: &ast.InterfaceType{Methods: &ast.FieldList{List: []*ast.Field{{Type: externalDoerSelectorExpr}}}},
			Expected: []Method{
				newTypesMethod(externalDoerInterface.Method(0)),
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
