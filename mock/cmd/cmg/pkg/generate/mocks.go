package generate

import (
	"context"
	_ "embed"
	"io"
	"sort"
	"text/template"

	"goa.design/clue/mock/cmd/cmg/pkg/parse"
)

type (
	Mocks interface {
		PkgName() string
		PkgImport() Import
		StdImports() []Import
		ExtImports() []Import
		IntImports() []Import
		Interfaces() []Interface
		ToolVersion() string
		ToolCommandLine() string
		Render(ctx context.Context, w io.Writer) error
	}

	ToolVersionFunc func() string

	mocks struct {
		pkgName, pkgPath                   string
		pkgImport                          Import
		stdImports, extImports, intImports []Import
		interfaces                         []Interface
		toolVersionFunc                    ToolVersionFunc
	}
)

//go:embed mocks.go.tmpl
var mocksStr string

var (
	mocksTmpl = template.Must(template.New("mocks").Parse(mocksStr))
)

func NewMocks(prefix string, p parse.Package, interfaces []parse.Interface, toolVersionFunc ToolVersionFunc) Mocks {
	var (
		m = &mocks{
			pkgName:         prefix + p.Name(),
			pkgPath:         p.PkgPath(),
			pkgImport:       newImport(p.PkgPath(), p.Name()),
			toolVersionFunc: toolVersionFunc,
		}
		stdImports = importMap{"testing": newImport("testing")}
		extImports = importMap{"mock": newImport("goa.design/clue/mock")}
		intImports = make(importMap)
		modPath    = p.ModPath()
		typeNames  = make(typeMap)
		typeZeros  = make(typeMap)
	)

	addImport(m.pkgImport, stdImports, extImports, intImports, modPath)

	for _, i := range interfaces {
		m.interfaces = append(m.interfaces, newInterface(i, typeNames, typeZeros, stdImports, extImports, intImports, modPath))
	}

	for _, i := range stdImports {
		m.stdImports = append(m.stdImports, i)
	}
	sort.Slice(m.stdImports, func(i, j int) bool { return m.stdImports[i].PkgPath() < m.stdImports[j].PkgPath() })

	for _, i := range extImports {
		m.extImports = append(m.extImports, i)
	}
	sort.Slice(m.extImports, func(i, j int) bool { return m.extImports[i].PkgPath() < m.extImports[j].PkgPath() })

	for _, i := range intImports {
		m.intImports = append(m.intImports, i)
	}
	sort.Slice(m.intImports, func(i, j int) bool { return m.intImports[i].PkgPath() < m.intImports[j].PkgPath() })

	return m
}

func (m *mocks) PkgName() string {
	return m.pkgName
}

func (m *mocks) PkgImport() Import {
	return m.pkgImport
}

func (m *mocks) StdImports() []Import {
	return m.stdImports
}

func (m *mocks) ExtImports() []Import {
	return m.extImports
}

func (m *mocks) IntImports() []Import {
	return m.intImports
}

func (m *mocks) Interfaces() []Interface {
	return m.interfaces
}

func (m *mocks) ToolVersion() string {
	return m.toolVersionFunc()
}

func (m *mocks) ToolCommandLine() string {
	return "$ cmg gen " + m.pkgPath
}

func (m *mocks) Render(ctx context.Context, w io.Writer) error {
	return mocksTmpl.Execute(w, m)
}
