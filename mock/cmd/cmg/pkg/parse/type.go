package parse

import (
	"go/ast"

	"golang.org/x/tools/go/packages"
)

type (
	Type interface {
		Name() string
	}

	type_ struct {
		p        *packages.Package
		ident    *ast.Ident
		typeType ast.Expr
	}
)

func newType(p *packages.Package, ident *ast.Ident, typeType ast.Expr) Type {
	return &type_{p: p, ident: ident, typeType: typeType}
}

func (t *type_) Name() string {
	return t.ident.Name
}
