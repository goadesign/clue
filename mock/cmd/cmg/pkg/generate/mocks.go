package generate

import (
	_ "embed"
	"io"
	"sort"
	"text/template"

	"goa.design/clue/mock/cmd/cmg/pkg/parse"
)

type (
	// Mocks is the interface for the mocks.
	Mocks interface {
		// PkgName returns the package name.
		PkgName() string
		// PkgImport returns the package import.
		PkgImport() Import
		// StdImports returns the standard imports.
		StdImports() []Import
		// ExtImports returns the external imports.
		ExtImports() []Import
		// IntImports returns the internal imports.
		IntImports() []Import
		// Interfaces returns the interfaces.
		Interfaces() []Interface
		// ToolVersion returns the tool version.
		ToolVersion() string
		// ToolCommandLine returns the tool command line.
		ToolCommandLine() string
		// Testify returns true if testify should be used.
		Testify() bool
		// Render renders the mocks to the given writer.
		Render(w io.Writer) error
	}

	// ToolVersionFunc is the function type for the tool version.
	ToolVersionFunc func() string

	// mocks is the implementation of the Mocks interface.
	mocks struct {
		pkgName, pkgPath                   string
		pkgImport                          Import
		stdImports, extImports, intImports []Import
		interfaces                         []Interface
		toolVersionFunc                    ToolVersionFunc
		testify                            bool
	}
)

//go:embed mocks.go.tmpl
var mocksStr string

var (
	// mocksTmpl is the template for the mocks.
	mocksTmpl = template.Must(template.New("mocks").Parse(mocksStr))
)

// NewMocks creates a new Mocks instance.
// If testify is true, it uses github.com/stretchr/testify for assertions.
func NewMocks(prefix string, p parse.Package, interfaces []parse.Interface, toolVersionFunc ToolVersionFunc, testify bool) Mocks {
	var (
		stdImports = importMap{"testing": newImport("testing")}
		extImports = importMap{"mock": newImport("goa.design/clue/mock")}
		intImports = make(importMap)
	)
	if testify {
		extImports["assert"] = newImport("github.com/stretchr/testify/assert")
	}

	var (
		modPath = p.ModPath()
		m       = &mocks{
			pkgName:         prefix + p.Name(),
			pkgPath:         p.PkgPath(),
			pkgImport:       addImport(newImport(p.PkgPath(), p.Name()), stdImports, extImports, intImports, modPath),
			toolVersionFunc: toolVersionFunc,
			testify:         testify,
		}
		typeNames = make(typeMap)
		typeZeros = make(typeMap)
	)

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

// PkgName returns the package name.
func (m *mocks) PkgName() string {
	return m.pkgName
}

// PkgImport returns the package import.
func (m *mocks) PkgImport() Import {
	return m.pkgImport
}

// StdImports returns the standard imports.
func (m *mocks) StdImports() []Import {
	return m.stdImports
}

// ExtImports returns the external imports.
func (m *mocks) ExtImports() []Import {
	return m.extImports
}

// IntImports returns the internal imports.
func (m *mocks) IntImports() []Import {
	return m.intImports
}

// Interfaces returns the interfaces.
func (m *mocks) Interfaces() []Interface {
	return m.interfaces
}

// ToolVersion returns the tool version.
func (m *mocks) ToolVersion() string {
	return m.toolVersionFunc()
}

// ToolCommandLine returns the tool command line.
func (m *mocks) ToolCommandLine() string {
	return "$ cmg gen " + m.pkgPath
}

// Testify returns true if testify should be used.
func (m *mocks) Testify() bool {
	return m.testify
}

// Render renders the mocks to the given writer.
func (m *mocks) Render(w io.Writer) error {
	return mocksTmpl.Execute(w, m)
}
