package generate

import (
	"fmt"
	"strings"

	"goa.design/clue/mock/cmd/cmg/pkg/parse"
)

type (
	Interface interface {
		Name() string
		Constructor() string
		ConstructorFmt(pkgName string) string
		TypeParameters() string
		TypeParameterVars() string
		Methods() []Method
		MaxFuncLenFmt() string
		Var() string
	}

	interface_ struct {
		parse.Interface

		methods              []Method
		maxFuncLen           int
		typeNames, typeZeros typeMap
	}
)

func newInterface(i parse.Interface, typeNames, typeZeros typeMap, stdImports, extImports, intImports importMap, modPath string) Interface {
	for _, p := range i.TypeParameters() {
		addType(p.Constraint(), typeNames, typeZeros, stdImports, extImports, intImports, modPath)
	}
	iface := &interface_{i, nil, 0, typeNames, typeZeros}
	for _, m := range i.Methods() {
		method := newMethod(m, i, typeNames, typeZeros, stdImports, extImports, intImports, modPath)
		if l := len(method.Func()); l > iface.maxFuncLen {
			iface.maxFuncLen = l
		}
		iface.methods = append(iface.methods, method)
	}
	return iface
}

func (i *interface_) Constructor() string {
	return "New" + i.Name()
}

func (i *interface_) ConstructorFmt(pkgName string) string {
	return fmt.Sprintf("%%-%vv", 3+len(pkgName)+len(i.Name())+len(i.TypeParameterVars()))
}

func (i *interface_) TypeParameters() string {
	ps := i.Interface.TypeParameters()
	if len(ps) == 0 {
		return ""
	}
	var parameters []string
	for j, p := range ps {
		parameter := p.Name()
		if !(j+1 < len(ps) && p.Constraint() == ps[j+1].Constraint()) {
			parameter += " " + i.typeNames[p.Constraint()]
		}
		parameters = append(parameters, parameter)
	}
	return "[" + strings.Join(parameters, ", ") + "]"
}

func (i *interface_) TypeParameterVars() string {
	ps := i.Interface.TypeParameters()
	if len(ps) == 0 {
		return ""
	}
	var vars []string
	for _, p := range ps {
		vars = append(vars, p.Name())
	}
	return "[" + strings.Join(vars, ", ") + "]"
}

func (i *interface_) Methods() []Method {
	return i.methods
}

func (i *interface_) Var() string {
	return "m"
}

func (i *interface_) MaxFuncLenFmt() string {
	return fmt.Sprintf("%%-%vv", i.maxFuncLen)
}
