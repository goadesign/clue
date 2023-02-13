package parse

import (
	"go/ast"

	"golang.org/x/tools/go/packages"
)

type (
	Method interface {
		Name() string
		IsExported() bool
		Parameters() []Value
		Results() []Value
		Variadic() bool
	}

	method struct {
		p        *packages.Package
		ident    *ast.Ident
		funcType *ast.FuncType
		variadic bool
	}
)

func newMethod(p *packages.Package, ident *ast.Ident, funcType *ast.FuncType, variadic bool) Method {
	return &method{p: p, ident: ident, funcType: funcType, variadic: variadic}
}

func (m *method) Name() string {
	return m.ident.Name
}

func (m *method) IsExported() bool {
	return m.ident.IsExported()
}

func (m *method) Parameters() (parameters []Value) {
	for _, p := range m.funcType.Params.List {
		idents := []*ast.Ident{nil}
		if p.Names != nil {
			idents = p.Names
		}
		for _, ident := range idents {
			parameters = append(parameters, newValue(m.p, ident, p.Type))
		}
	}
	return
}

func (m *method) Results() (results []Value) {
	if m.funcType.Results != nil {
		for _, r := range m.funcType.Results.List {
			idents := []*ast.Ident{nil}
			if r.Names != nil {
				idents = r.Names
			}
			for _, ident := range idents {
				results = append(results, newValue(m.p, ident, r.Type))
			}
		}
	}
	return
}

func (m *method) Variadic() bool {
	return m.variadic
}
