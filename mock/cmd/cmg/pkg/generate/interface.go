package generate

import (
	"fmt"

	"goa.design/clue/mock/cmd/cmg/pkg/parse"
)

type (
	Interface interface {
		Name() string
		Constructor() string
		Methods() []Method
		MaxFuncLenFmt() string
		Var() string
	}

	interface_ struct {
		parse.Interface

		methods    []Method
		maxFuncLen int
	}
)

func newInterface(i parse.Interface, typeNames, typeZeros typeMap, stdImports, extImports, intImports importMap, modPath string) Interface {
	var (
		methods    []Method
		maxFuncLen int
	)
	for _, m := range i.Methods() {
		method := newMethod(m, i, typeNames, typeZeros, stdImports, extImports, intImports, modPath)
		if l := len(method.Func()); l > maxFuncLen {
			maxFuncLen = l
		}
		methods = append(methods, method)
	}
	return &interface_{i, methods, maxFuncLen}
}

func (i *interface_) Constructor() string {
	return "New" + i.Name()
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
