package generate

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"goa.design/clue/mock/cmd/cmg/pkg/parse"
)

func TestMocks_Render(t *testing.T) {
	cases := []struct {
		Name, Pattern string
		ExpectedFiles []string
	}{
		{
			Name:          "extensive",
			Pattern:       "./extensive",
			ExpectedFiles: []string{"extensive.go"},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			assert := assert.New(t)
			require := require.New(t)

			ps, err := parse.LoadPackages([]string{tc.Pattern}, "_tests")
			require.NoError(err)
			require.Len(ps, 1)
			p := ps[0]

			is, err := p.Interfaces()
			require.NoError(err)

			interfacesByFile := make(map[string][]parse.Interface)
			for _, i := range is {
				f := filepath.Base(i.File())
				interfacesByFile[f] = append(interfacesByFile[f], i)
			}

			var files []string
			for f := range interfacesByFile {
				files = append(files, f)
			}
			assert.ElementsMatch(tc.ExpectedFiles, files)

			for f, is := range interfacesByFile {
				f := filepath.Join("_tests", tc.Pattern, "mocks", f)
				m := NewMocks("mock", p, is, func() string { return "TEST VERSION" })
				b := &bytes.Buffer{}

				err := m.Render(b)
				if assert.NoError(err) && assert.FileExists(f) {
					expected, err := os.ReadFile(f)
					if assert.NoError(err) {
						assert.Equal(string(expected), b.String())
					}
				}
			}
		})
	}
}
