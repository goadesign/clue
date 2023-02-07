package parse

import (
	"go/ast"

	"golang.org/x/tools/go/packages"
)

type (
	Method interface {
		Name() string
		IsExported() bool
		TypeParameters() []Type
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

func (m *method) TypeParameters() (typeParameters []Type) {
	if m.funcType.TypeParams != nil {
		for _, tp := range m.funcType.TypeParams.List {
			typeParameters = append(typeParameters, newType(m.p, tp.Names[0], tp.Type))
		}
	}
	return
}

func (m *method) Parameters() (parameters []Value) {
	for _, p := range m.funcType.Params.List {
		var ident *ast.Ident
		if p.Names != nil {
			ident = p.Names[0]
		}
		parameters = append(parameters, newValue(m.p, ident, p.Type))
	}
	return
}

func (m *method) Results() (results []Value) {
	if m.funcType.Results != nil {
		for _, r := range m.funcType.Results.List {
			var ident *ast.Ident
			if r.Names != nil {
				ident = r.Names[0]
			}
			results = append(results, newValue(m.p, ident, r.Type))
		}
	}
	return
}

func (m *method) Variadic() bool {
	return m.variadic
}
