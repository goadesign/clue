package parse

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/packages"
)

func TestLoadPackages(t *testing.T) {
	cases := []struct {
		Name, ExpectedErr string
		Patterns          []string
		Expected          []struct{ Name, PkgPath, ModPath string }
	}{
		{
			Name:     "success",
			Patterns: []string{"./doer"},
			Expected: []struct{ Name, PkgPath, ModPath string }{
				{Name: "doer", PkgPath: "example.com/a/b/doer", ModPath: "example.com/a/b"},
			},
		},
		{
			Name:        "error",
			Patterns:    []string{"nonexistent=bogus"},
			ExpectedErr: `invalid query type "nonexistent" in query pattern "nonexistent=bogus"`,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			assert := assert.New(t)

			packages, err := LoadPackages(tc.Patterns, "_tests")
			if tc.ExpectedErr != "" {
				assert.ErrorContains(err, tc.ExpectedErr)
			} else if assert.NoError(err) && assert.Len(packages, len(tc.Expected)) {
				for i, p := range packages {
					expected := tc.Expected[i]
					assert.Equal(expected.Name, p.Name())
					assert.Equal(expected.PkgPath, p.PkgPath())
					assert.Equal(expected.ModPath, p.ModPath())
				}
			}
		})
	}
}

func TestPackage_Name(t *testing.T) {
	cases := []struct {
		Name, Expected string
		Package        *packages.Package
	}{
		{
			Name:     "success",
			Package:  &packages.Package{Name: "a"},
			Expected: "a",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			p := &packageImpl{tc.Package}
			name := p.Name()
			assert.Equal(t, tc.Expected, name)
		})
	}
}

func TestPackage_PkgPath(t *testing.T) {
	cases := []struct {
		Name, Expected string
		Package        *packages.Package
	}{
		{
			Name:     "success",
			Package:  &packages.Package{PkgPath: "example.com/a/b/c"},
			Expected: "example.com/a/b/c",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			p := &packageImpl{tc.Package}
			path := p.PkgPath()
			assert.Equal(t, tc.Expected, path)
		})
	}
}

func TestPackage_ModPath(t *testing.T) {
	cases := []struct {
		Name, Expected string
		Package        *packages.Package
	}{
		{
			Name:     "success",
			Package:  &packages.Package{Module: &packages.Module{Path: "example.com/a/b"}},
			Expected: "example.com/a/b",
		},
		{
			Name:     "empty",
			Package:  &packages.Package{},
			Expected: "",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			p := &packageImpl{tc.Package}
			path := p.ModPath()
			assert.Equal(t, tc.Expected, path)
		})
	}
}

func TestPackage_Interfaces(t *testing.T) {
	cases := []struct {
		Name, Pattern, ExpectedErr string
		Expected                   []struct {
			Name, File string
			IsExported bool
		}
	}{
		{
			Name:    "success",
			Pattern: "./doer",
			Expected: []struct {
				Name, File string
				IsExported bool
			}{
				{Name: "Doer", File: "doer.go", IsExported: true},
				{Name: "EmbeddedDoer", File: "doer.go", IsExported: true},
				{Name: "ExternalEmbeddedDoer", File: "doer.go", IsExported: true},
				{Name: "doer", File: "doer.go", IsExported: false},
			},
		},
		{
			Name:        "error",
			Pattern:     ".",
			ExpectedErr: "-: no Go files in ",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			assert := assert.New(t)

			packages, err := LoadPackages([]string{tc.Pattern}, "_tests")
			if assert.NoError(err) && assert.Len(packages, 1) {
				interfaces, err := packages[0].Interfaces()
				if tc.ExpectedErr != "" {
					assert.ErrorContains(err, tc.ExpectedErr)
				} else if assert.NoError(err) && assert.Len(interfaces, len(tc.Expected)) {
					for i, iface := range interfaces {
						expected := tc.Expected[i]
						assert.Equal(expected.Name, iface.Name())
						assert.Equal(expected.File, filepath.Base(iface.File()))
						assert.Equal(expected.IsExported, iface.IsExported())
					}
				}
			}
		})
	}
}
