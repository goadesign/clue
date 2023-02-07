package parse

import (
	"go/ast"

	"golang.org/x/tools/go/packages"
)

type (
	Package interface {
		Name() string
		PkgPath() string
		ModPath() string
		Interfaces() ([]Interface, error)
	}

	package_ struct {
		p *packages.Package
	}

	interfaceVisitor struct {
		p          *packages.Package
		file       string
		interfaces []Interface
	}
)

func LoadPackages(patterns []string) ([]Package, error) {
	c := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedImports | packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo | packages.NeedModule,
	}

	ps, err := packages.Load(c, patterns...)
	if err != nil {
		return nil, err
	}

	packages := make([]Package, len(ps))
	for i, p := range ps {
		packages[i] = &package_{p: p}
	}

	return packages, nil
}

func (p *package_) Name() string {
	return p.p.Name
}

func (p *package_) PkgPath() string {
	return p.p.PkgPath
}

func (p *package_) ModPath() string {
	return p.p.Module.Path
}

func (p *package_) Interfaces() ([]Interface, error) {
	if len(p.p.Errors) > 0 {
		return nil, p.p.Errors[0]
	}

	var interfaces []Interface
	for i, gf := range p.p.GoFiles {
		iv := &interfaceVisitor{p: p.p, file: gf}
		ast.Walk(iv, p.p.Syntax[i])
		interfaces = append(interfaces, iv.interfaces...)
	}

	return interfaces, nil
}

func (iv *interfaceVisitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.TypeSpec:
		switch t := n.Type.(type) {
		case *ast.InterfaceType:
			iv.interfaces = append(iv.interfaces, newInterface(iv.p, iv.file, n, t))
		}
	}
	return iv
}
