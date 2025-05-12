package generate

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"goa.design/clue/mock/cmd/cmg/pkg/parse"
)

var updateGolden = false

func init() {
	flag.BoolVar(&updateGolden, "update-golden", false, "update golden files")
}

func TestMocks_Render(t *testing.T) {
	cases := []struct {
		Name, Pattern string
		ExpectedFiles []string
		Testify       bool
	}{
		{
			Name:          "extensive",
			Pattern:       "./extensive",
			ExpectedFiles: []string{"extensive.go"},
		},
		{
			Name:          "conflicts",
			Pattern:       "./conflicts",
			ExpectedFiles: []string{"conflicts.go"},
		},
		{
			Name:          "testify",
			Pattern:       "./testify",
			ExpectedFiles: []string{"testify.go"},
			Testify:       true,
		},
	}

	for _, tc := range cases {
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

			mocksDir := filepath.Join("_tests", tc.Pattern, "mocks")

			if updateGolden {
				require.NoError(os.MkdirAll(mocksDir, 0750))

				for f, is := range interfacesByFile {
					of, err := os.OpenFile(filepath.Join(mocksDir, f), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0640)
					require.NoError(err)
					t.Cleanup(func() { assert.NoError(of.Close()) })

					m := NewMocks("mock", p, is, toolVersion, tc.Testify)
					require.NoError(m.Render(of))
				}

				return
			}

			var files []string
			for f := range interfacesByFile {
				files = append(files, f)
			}
			assert.ElementsMatch(tc.ExpectedFiles, files)

			for f, is := range interfacesByFile {
				f := filepath.Join(mocksDir, f)
				m := NewMocks("mock", p, is, toolVersion, tc.Testify)
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

func toolVersion() string {
	return "TEST VERSION"
}
