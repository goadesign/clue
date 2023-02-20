package generate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddImport(t *testing.T) {
	cases := []struct {
		Name, ExpectedLocation             string
		Import, ExpectedImport             Import
		StdImports, ExtImports, IntImports importMap
	}{
		{
			Name:             "standard library duplicate",
			Import:           newImport("text/template"),
			StdImports:       importMap{"template": newImport("html/template")},
			ExpectedImport:   newImport("text/template", "template", "template1"),
			ExpectedLocation: "std",
		},
		{
			Name:             "external duplicate",
			Import:           newImport("example.com/a/b/template"),
			StdImports:       importMap{"template": newImport("html/template"), "template1": newImport("text/template", "template", "template1")},
			IntImports:       importMap{"template2": newImport("example.com/b/c/template", "template", "template2")},
			ExpectedImport:   newImport("example.com/a/b/template", "template", "template3"),
			ExpectedLocation: "ext",
		},
		{
			Name:             "internal duplicate",
			Import:           newImport("example.com/b/c/template"),
			StdImports:       importMap{"template": newImport("html"), "template1": newImport("text/template", "template", "template1")},
			ExtImports:       importMap{"template2": newImport("example.com/a/b/template", "template", "template2")},
			ExpectedImport:   newImport("example.com/b/c/template", "template", "template3"),
			ExpectedLocation: "int",
		},
		{
			Name:             "already aliased",
			Import:           newImport("text/template"),
			StdImports:       importMap{"template": newImport("html/template"), "template1": newImport("text/template", "template", "template1")},
			ExpectedImport:   newImport("text/template", "template", "template1"),
			ExpectedLocation: "std",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			assert := assert.New(t)

			if tc.StdImports == nil {
				tc.StdImports = make(importMap)
			}
			if tc.ExtImports == nil {
				tc.ExtImports = make(importMap)
			}
			if tc.IntImports == nil {
				tc.IntImports = make(importMap)
			}

			i := addImport(tc.Import, tc.StdImports, tc.ExtImports, tc.IntImports, "example.com/b/c")
			if assert.Equal(tc.ExpectedImport, i) && assert.Contains([]string{"std", "ext", "int"}, tc.ExpectedLocation) {
				a := i.AliasOrPkgName()
				switch tc.ExpectedLocation {
				case "std":
					assert.Contains(tc.StdImports, a)
					assert.NotContains(tc.ExtImports, a)
					assert.NotContains(tc.IntImports, a)
				case "ext":
					assert.NotContains(tc.StdImports, a)
					assert.Contains(tc.ExtImports, a)
					assert.NotContains(tc.IntImports, a)
				case "int":
					assert.NotContains(tc.StdImports, a)
					assert.NotContains(tc.ExtImports, a)
					assert.Contains(tc.IntImports, a)
				}
			}
		})
	}
}
