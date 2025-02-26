package parse

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/packages"
)

type (
	Interface interface {
		Name() string
		IsExported() bool
		File() string
		TypeParameters() []Type
		Methods() []Method
	}

	interfaceImpl struct {
		p             *packages.Package
		file          string
		typeSpec      *ast.TypeSpec
		interfaceType *ast.InterfaceType
	}
)

func newInterface(p *packages.Package, file string, typeSpec *ast.TypeSpec, interfaceType *ast.InterfaceType) Interface {
	return &interfaceImpl{p: p, file: file, typeSpec: typeSpec, interfaceType: interfaceType}
}

func (i *interfaceImpl) Name() string {
	return i.typeSpec.Name.Name
}

func (i *interfaceImpl) IsExported() bool {
	return i.typeSpec.Name.IsExported()
}

func (i *interfaceImpl) File() string {
	return i.file
}

func (i *interfaceImpl) TypeParameters() (typeParameters []Type) {
	if i.typeSpec.TypeParams != nil {
		for _, tp := range i.typeSpec.TypeParams.List {
			for _, ident := range tp.Names {
				typeParameters = append(typeParameters, newType(i.p, ident, tp.Type))
			}
		}
	}
	return
}

func (i *interfaceImpl) Methods() []Method {
	return i.methods(i.interfaceType)
}

func (i *interfaceImpl) methods(it *ast.InterfaceType) (methods []Method) {
	for _, m := range it.Methods.List {
		switch t := m.Type.(type) {
		case *ast.FuncType:
			for _, n := range m.Names {
				o, _, _ := types.LookupFieldOrMethod(i.p.Types.Scope().Lookup(i.Name()).Type(), true, i.p.Types, n.Name)
				methods = append(methods, newASTMethod(i.p, n, t, o.Type().Underlying().(*types.Signature).Variadic()))
			}
		case *ast.Ident:
			switch dt := t.Obj.Decl.(type) {
			case *ast.TypeSpec:
				switch t := dt.Type.(type) {
				case *ast.InterfaceType:
					methods = append(methods, i.methods(t)...)
				}
			}
		case *ast.SelectorExpr:
			if tv, ok := i.p.TypesInfo.Types[t]; ok {
				if ti, ok := tv.Type.Underlying().(*types.Interface); ok {
					for j := 0; j < ti.NumMethods(); j++ {
						methods = append(methods, newTypesMethod(ti.Method(j)))
					}
				}
			}
		}
	}
	return
}
