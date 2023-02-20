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

	interface_ struct {
		p             *packages.Package
		file          string
		typeSpec      *ast.TypeSpec
		interfaceType *ast.InterfaceType
	}
)

func newInterface(p *packages.Package, file string, typeSpec *ast.TypeSpec, interfaceType *ast.InterfaceType) Interface {
	return &interface_{p: p, file: file, typeSpec: typeSpec, interfaceType: interfaceType}
}

func (i *interface_) Name() string {
	return i.typeSpec.Name.Name
}

func (i *interface_) IsExported() bool {
	return i.typeSpec.Name.IsExported()
}

func (i *interface_) File() string {
	return i.file
}

func (i *interface_) TypeParameters() (typeParameters []Type) {
	if i.typeSpec.TypeParams != nil {
		for _, tp := range i.typeSpec.TypeParams.List {
			for _, ident := range tp.Names {
				typeParameters = append(typeParameters, newType(i.p, ident, tp.Type))
			}
		}
	}
	return
}

func (i *interface_) Methods() (methods []Method) {
	for _, m := range i.interfaceType.Methods.List {
		switch t := m.Type.(type) {
		case *ast.FuncType:
			for _, n := range m.Names {
				o, _, _ := types.LookupFieldOrMethod(i.p.Types.Scope().Lookup(i.Name()).Type(), true, i.p.Types, n.Name)
				methods = append(methods, newMethod(i.p, n, t, o.Type().Underlying().(*types.Signature).Variadic()))
			}
		}
	}
	return
}
