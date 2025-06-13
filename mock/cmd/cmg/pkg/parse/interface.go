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

	interfaceAlias struct {
		Interface
		aliasedInterface *types.Interface
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
		case *ast.Ident, *ast.SelectorExpr, *ast.IndexExpr, *ast.IndexListExpr:
			if tv, ok := i.p.TypesInfo.Types[t]; ok {
				if ti, ok := tv.Type.Underlying().(*types.Interface); ok {
					for m := range ti.Methods() {
						methods = append(methods, newTypesMethod(m))
					}
				}
			}
		}
	}
	return
}

func newInterfaceAlias(p *packages.Package, file string, typeSpec *ast.TypeSpec, aliasedInterface *types.Interface) Interface {
	return &interfaceAlias{
		Interface:        newInterface(p, file, typeSpec, nil),
		aliasedInterface: aliasedInterface,
	}
}

func (i *interfaceAlias) Methods() (methods []Method) {
	for m := range i.aliasedInterface.Methods() {
		methods = append(methods, newTypesMethod(m))
	}
	return
}
