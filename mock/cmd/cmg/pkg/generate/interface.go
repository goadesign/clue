package generate

import (
	"fmt"
	"math"
	"strconv"
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
		HasMore() string
	}

	interfaceImpl struct {
		parse.Interface

		methods              []Method
		maxFuncLen           int
		typeNames, typeZeros typeMap
		hasMoreName          string
	}

	interfaceScope map[string]struct{}
)

func newInterface(i parse.Interface, typeNames, typeZeros typeMap, stdImports, extImports, intImports importMap, modPath string) Interface {
	for _, p := range i.TypeParameters() {
		addType(p.Constraint(), typeNames, typeZeros, stdImports, extImports, intImports, modPath)
	}
	var (
		iface = &interfaceImpl{Interface: i, typeNames: typeNames, typeZeros: typeZeros}
		ms    = i.Methods()
		is    = newInterfaceScope(ms)
	)
	iface.hasMoreName = is.uniqueName("HasMore")
	for _, m := range ms {
		method := newMethod(m, i, typeNames, typeZeros, stdImports, extImports, intImports, modPath, is)
		if l := len(method.Func() + iface.TypeParameters()); l > iface.maxFuncLen {
			iface.maxFuncLen = l
		}
		iface.methods = append(iface.methods, method)
	}
	return iface
}

func (i *interfaceImpl) Constructor() string {
	return "New" + i.Name()
}

func (i *interfaceImpl) ConstructorFmt(pkgName string) string {
	return fmt.Sprintf("%%-%vv", 3+len(pkgName)+len(i.Name())+len(i.TypeParameterVars()))
}

func (i *interfaceImpl) TypeParameters() string {
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

func (i *interfaceImpl) TypeParameterVars() string {
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

func (i *interfaceImpl) Methods() []Method {
	return i.methods
}

func (i *interfaceImpl) Var() string {
	return "m"
}

func (i *interfaceImpl) MaxFuncLenFmt() string {
	return fmt.Sprintf("%%-%vv", i.maxFuncLen)
}

func (i *interfaceImpl) HasMore() string {
	return i.hasMoreName
}

func newInterfaceScope(ms []parse.Method) interfaceScope {
	is := make(interfaceScope, len(ms))
	for _, m := range ms {
		is[m.Name()] = struct{}{}
	}
	return is
}

func (is interfaceScope) uniqueName(name string) string {
	_, ok := is[name]
	if !ok {
		is[name] = struct{}{}
		return name
	}

	name += "Mock"
	if _, ok := is[name]; !ok {
		is[name] = struct{}{}
		return name
	}

	var newName string
	for i := 1; i <= math.MaxInt; i++ {
		newName = name + strconv.Itoa(i)
		if _, ok := is[newName]; !ok {
			is[newName] = struct{}{}
			break
		}
	}
	return newName
}
