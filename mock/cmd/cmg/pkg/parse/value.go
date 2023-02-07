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

	value struct {
		p         *packages.Package
		ident     *ast.Ident
		valueType ast.Expr
	}
)

func newValue(p *packages.Package, ident *ast.Ident, valueType ast.Expr) Value {
	return &value{p: p, ident: ident, valueType: valueType}
}

func (v *value) Name() string {
	if v.ident != nil {
		return v.ident.Name
	}
	return ""
}

func (v *value) Type() types.Type {
	return v.p.TypesInfo.Types[v.valueType].Type
}
