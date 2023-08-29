package parse

import (
	"go/ast"
	"go/types"

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

	astMethod struct {
		p        *packages.Package
		ident    *ast.Ident
		funcType *ast.FuncType
		variadic bool
	}

	typesMethod struct {
		f *types.Func
		s *types.Signature
	}
)

func newASTMethod(p *packages.Package, ident *ast.Ident, funcType *ast.FuncType, variadic bool) Method {
	return &astMethod{p: p, ident: ident, funcType: funcType, variadic: variadic}
}

func (am *astMethod) Name() string {
	return am.ident.Name
}

func (am *astMethod) IsExported() bool {
	return am.ident.IsExported()
}

func (am *astMethod) Parameters() (parameters []Value) {
	for _, p := range am.funcType.Params.List {
		idents := []*ast.Ident{nil}
		if p.Names != nil {
			idents = p.Names
		}
		for _, ident := range idents {
			parameters = append(parameters, newASTValue(am.p, ident, p.Type))
		}
	}
	return
}

func (am *astMethod) Results() (results []Value) {
	if am.funcType.Results != nil {
		for _, r := range am.funcType.Results.List {
			idents := []*ast.Ident{nil}
			if r.Names != nil {
				idents = r.Names
			}
			for _, ident := range idents {
				results = append(results, newASTValue(am.p, ident, r.Type))
			}
		}
	}
	return
}

func (am *astMethod) Variadic() bool {
	return am.variadic
}

func newTypesMethod(f *types.Func) Method {
	return &typesMethod{f: f, s: f.Type().(*types.Signature)}
}

func (tm *typesMethod) Name() string {
	return tm.f.Name()
}

func (tm *typesMethod) IsExported() bool {
	return tm.f.Exported()
}

func (tm *typesMethod) Parameters() (parameters []Value) {
	for i := 0; i < tm.s.Params().Len(); i++ {
		parameters = append(parameters, newTypesValue(tm.s.Params().At(i)))
	}
	return
}

func (tm *typesMethod) Results() (results []Value) {
	for i := 0; i < tm.s.Results().Len(); i++ {
		results = append(results, newTypesValue(tm.s.Results().At(i)))
	}
	return
}

func (tm *typesMethod) Variadic() bool {
	return tm.s.Variadic()
}
