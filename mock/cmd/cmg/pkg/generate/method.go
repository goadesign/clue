package generate

import (
	"fmt"
	"math"
	"strings"

	"goa.design/clue/mock/cmd/cmg/pkg/parse"
)

type (
	Method interface {
		Name() string
		Func() string
		Add() string
		Set() string
		InterfaceVar() string
		FuncVar() string
		Parameters() string
		ParameterVars() string
		Results() string
		ZeroResults() string
	}

	method struct {
		parse.Method

		funcName, addName, setName, interfaceVar, funcVar string
		typeNames, typeZeros                              typeMap
	}
)

func newMethod(m parse.Method, i parse.Interface, typeNames, typeZeros typeMap, stdImports, extImports, intImports importMap, modPath string, is interfaceScope) Method {
	parameterVars := make(map[string]struct{})
	for _, t := range m.Parameters() {
		parameterVars[t.Name()] = struct{}{}
		addType(t.Type(), typeNames, typeZeros, stdImports, extImports, intImports, modPath)
	}
	for _, t := range m.Results() {
		addType(t.Type(), typeNames, typeZeros, stdImports, extImports, intImports, modPath)
	}
	return &method{
		Method:       m,
		funcName:     i.Name() + m.Name() + "Func",
		addName:      is.uniqueName("Add" + m.Name()),
		setName:      is.uniqueName("Set" + m.Name()),
		interfaceVar: uniqueVar("m", parameterVars),
		funcVar:      uniqueVar("f", parameterVars),
		typeNames:    typeNames,
		typeZeros:    typeZeros,
	}
}

func (m *method) Func() string {
	return m.funcName
}

func (m *method) Add() string {
	return m.addName
}

func (m *method) Set() string {
	return m.setName
}

func (m *method) InterfaceVar() string {
	return m.interfaceVar
}

func (m *method) FuncVar() string {
	return m.funcVar
}

func (m *method) Parameters() string {
	var (
		parameters []string
		ps         = m.Method.Parameters()
		last       = len(ps) - 1
	)
	for i, p := range ps {
		b := &strings.Builder{}
		if n, _ := b.WriteString(p.Name()); n == 0 {
			fmt.Fprintf(b, "p%v", i)
		}
		if m.Method.Variadic() && i == last {
			b.WriteString(" ..." + m.typeNames[p.Type()][2:])
		} else if !(i+1 < len(ps) && p.Type() == ps[i+1].Type()) {
			b.WriteString(" " + m.typeNames[p.Type()])
		}
		parameters = append(parameters, b.String())
	}
	return strings.Join(parameters, ", ")
}

func (m *method) ParameterVars() string {
	var (
		vars []string
		ps   = m.Method.Parameters()
		last = len(ps) - 1
	)
	for i, p := range ps {
		v := p.Name()
		if v == "" {
			v = fmt.Sprintf("p%v", i)
		}
		if m.Method.Variadic() && i == last {
			v += "..."
		}
		vars = append(vars, v)
	}
	return strings.Join(vars, ", ")
}

func (m *method) Results() string {
	var (
		results []string
		rs      = m.Method.Results()
		named   bool
	)
	for i, r := range rs {
		var b strings.Builder
		if n, _ := b.WriteString(r.Name()); n > 0 {
			named = true
		}
		if named {
			if !(i+1 < len(rs) && r.Type() == rs[i+1].Type()) {
				b.WriteString(" " + m.typeNames[r.Type()])
			}
		} else {
			b.WriteString(m.typeNames[r.Type()])
		}
		results = append(results, b.String())
	}
	switch len(results) {
	case 0:
		return ""
	case 1:
		if !named {
			return results[0]
		}
		fallthrough
	default:
		return "(" + strings.Join(results, ", ") + ")"
	}
}

func (m *method) ZeroResults() string {
	var results []string
	for _, r := range m.Method.Results() {
		results = append(results, m.typeZeros[r.Type()])
	}
	return strings.Join(results, ", ")
}

func uniqueVar(prefix string, vars map[string]struct{}) string {
	v := prefix
	if _, ok := vars[v]; ok {
		for i := 1; i <= math.MaxInt; i++ {
			v = fmt.Sprintf("%v%v", prefix, i)
			if _, ok = vars[v]; !ok {
				break
			}
		}
	}
	return v
}
