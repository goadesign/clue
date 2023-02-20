package generate

import (
	"go/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	fakeType struct {
		Name string
	}
)

func (t *fakeType) Underlying() types.Type {
	return t
}

func (t *fakeType) String() string {
	return t.Name
}

func TestTypeAdder_name(t *testing.T) {
	cases := []struct {
		Name, ExpectedPanic string
		Type                types.Type
	}{
		{
			Name:          "unhandled type",
			Type:          &fakeType{"unhandled"},
			ExpectedPanic: `unknown name for type: &generate.fakeType{Name:"unhandled"} (*generate.fakeType)`,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			assert := assert.New(t)

			ta := typeAdder{}

			if tc.ExpectedPanic != "" {
				assert.PanicsWithError(tc.ExpectedPanic, func() {
					ta.name(tc.Type)
				})
			}
		})
	}
}

func TestTypeAdder_zero(t *testing.T) {
	cases := []struct {
		Name, ExpectedPanic string
		Type                types.Type
	}{
		{
			Name:          "unhandled type",
			Type:          &fakeType{"unhandled"},
			ExpectedPanic: `unknown zero for type: &generate.fakeType{Name:"unhandled"} (*generate.fakeType)`,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			assert := assert.New(t)

			ta := typeAdder{}

			if tc.ExpectedPanic != "" {
				assert.PanicsWithError(tc.ExpectedPanic, func() {
					ta.zero(tc.Type)
				})
			}
		})
	}
}
