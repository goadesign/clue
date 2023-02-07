package generate

import (
	"fmt"
	"strings"

	"goa.design/clue/mock/cmd/cmg/pkg/parse"
)

type (
	Method interface {
		Name() string
		Func() string
		Add() string
		Set() string
		Parameters() string
		ParameterVars() string
		Results() string
		ZeroResults() string
	}

	method struct {
		parse.Method

		func_, add, set      string
		typeNames, typeZeros typeMap
	}
)

func newMethod(m parse.Method, i parse.Interface, typeNames, typeZeros typeMap, stdImports, extImports, intImports importMap, modPath string) Method {
	for _, t := range append(m.Parameters(), m.Results()...) {
		addType(t.Type(), typeNames, typeZeros, stdImports, extImports, intImports, modPath)
	}
	return &method{
		m,
		i.Name() + m.Name() + "Func",
		"Add" + m.Name(),
		"Set" + m.Name(),
		typeNames,
		typeZeros,
	}
}

func (m *method) Func() string {
	return m.func_
}

func (m *method) Add() string {
	return m.add
}

func (m *method) Set() string {
	return m.set
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
			fmt.Fprintf(b, "a%v", i)
		}
		b.WriteString(" ")
		if m.Method.Variadic() && i == last {
			b.WriteString("..." + m.typeNames[p.Type()][2:])
		} else {
			b.WriteString(m.typeNames[p.Type()])
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
			v = fmt.Sprintf("a%v", i)
		}
		if m.Method.Variadic() && i == last {
			v = v + "..."
		}
		vars = append(vars, v)
	}
	return strings.Join(vars, ", ")
}

func (m *method) Results() string {
	var (
		results []string
		named   bool
	)
	for _, r := range m.Method.Results() {
		var b strings.Builder
		if n, _ := b.WriteString(r.Name()); n > 0 {
			named = true
			b.WriteString(" ")
		}
		b.WriteString(m.typeNames[r.Type()])
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
