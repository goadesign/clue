package parse

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/packages"
)

type (
	Value interface {
		Name() string
		Type() types.Type
	}

	astValue struct {
		p         *packages.Package
		ident     *ast.Ident
		valueType ast.Expr
	}

	typesValue struct {
		v *types.Var
	}
)

func newASTValue(p *packages.Package, ident *ast.Ident, valueType ast.Expr) Value {
	return &astValue{p: p, ident: ident, valueType: valueType}
}

func (av *astValue) Name() string {
	if av.ident != nil {
		return av.ident.Name
	}
	return ""
}

func (av *astValue) Type() types.Type {
	return av.p.TypesInfo.Types[av.valueType].Type
}

func newTypesValue(v *types.Var) Value {
	return &typesValue{v: v}
}

func (tv *typesValue) Name() string {
	return tv.v.Name()
}

func (tv *typesValue) Type() types.Type {
	return tv.v.Type()
}
