package mock

import (
	"testing"
)

type namedFunc struct {
	Name string
	Func func() string
}

func TestAdd(t *testing.T) {
	cases := []struct {
		Name     string
		Funcs    []namedFunc
		NextArgs []string
		Expected []func() string
	}{
		{
			Name:     "single",
			Funcs:    []namedFunc{{"f1", f1}},
			NextArgs: []string{"f1"},
			Expected: []func() string{f1},
		},
		{
			Name:     "multiple",
			Funcs:    []namedFunc{{"f1", f1}, {"f2", f2}},
			NextArgs: []string{"f1", "f2"},
			Expected: []func() string{f1, f2},
		},
		{
			Name:     "nofunc",
			Funcs:    []namedFunc{{"f1", f1}},
			NextArgs: []string{"f2"},
			Expected: []func() string{nil},
		},
		{
			Name:     "consumed",
			Funcs:    []namedFunc{{"f1", f1}},
			NextArgs: []string{"f1", "f1"},
			Expected: []func() string{f1, nil},
		},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			m := New()
			for _, nf := range c.Funcs {
				m.Add(nf.Name, nf.Func)
			}
			for i, a := range c.NextArgs {
				got := m.Next(a)
				compareFuncs(t, got, c.Expected[i])
			}
		})
	}
}

func TestSet(t *testing.T) {
	cases := []struct {
		Name     string
		Funcs    []namedFunc
		NextArgs []string
		Expected []func() string
	}{
		{
			Name:     "single",
			Funcs:    []namedFunc{{"f1", f1}},
			NextArgs: []string{"f1"},
			Expected: []func() string{f1},
		},
		{
			Name:     "multiple",
			Funcs:    []namedFunc{{"f1", f1}, {"f2", f2}},
			NextArgs: []string{"f1", "f2"},
			Expected: []func() string{f1, f2},
		},
		{
			Name:     "nofunc",
			Funcs:    []namedFunc{{"f1", f1}},
			NextArgs: []string{"f2"},
			Expected: []func() string{nil},
		},
		{
			Name:     "permanent",
			Funcs:    []namedFunc{{"f1", f1}},
			NextArgs: []string{"f1", "f1"},
			Expected: []func() string{f1, f1},
		},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			m := New()
			for _, nf := range c.Funcs {
				m.Set(nf.Name, nf.Func)
			}
			for i, a := range c.NextArgs {
				got := m.Next(a)
				compareFuncs(t, got, c.Expected[i])
			}
		})
	}
}

func TestHasMore(t *testing.T) {
	cases := []struct {
		Name     string
		Funcs    []namedFunc
		NextArgs []string
		Expected []bool
	}{
		{
			Name:     "no sequence",
			Expected: []bool{false},
		},
		{
			Name:     "single",
			Funcs:    []namedFunc{{"f1", f1}},
			NextArgs: []string{"f1"},
			Expected: []bool{true, false},
		},
		{
			Name:     "multiple",
			Funcs:    []namedFunc{{"f1", f1}, {"f2", f2}},
			NextArgs: []string{"f1", "f2"},
			Expected: []bool{true, true, false},
		},
		{
			Name:     "nofunc",
			Funcs:    []namedFunc{{"f1", f1}},
			NextArgs: []string{"f2"},
			Expected: []bool{true, true},
		},
		{
			Name:     "consumed",
			Funcs:    []namedFunc{{"f1", f1}},
			NextArgs: []string{"f1", "f1"},
			Expected: []bool{true, false, false},
		},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			m := New()
			for _, nf := range c.Funcs {
				m.Add(nf.Name, nf.Func)
			}
			if c.Expected[0] != m.HasMore() {
				t.Errorf("HasMore() should be %v", c.Expected[0])
			}
			for i, a := range c.NextArgs {
				m.Next(a)
				got := m.HasMore()
				if got != c.Expected[i+1] {
					t.Errorf("got %v, want %v", got, c.Expected[i])
				}
			}
		})
	}
}

func TestNext(t *testing.T) {
	cases := []struct {
		Name     string
		Seq      []namedFunc
		Set      namedFunc
		NextArgs []string
		Expected []func() string
	}{
		{
			Name:     "no sequence",
			Expected: []func() string{nil},
		},
		{
			Name:     "single",
			Seq:      []namedFunc{{"f1", f1}},
			NextArgs: []string{"f1"},
			Expected: []func() string{f1},
		},
		{
			Name:     "multiple",
			Seq:      []namedFunc{{"f1", f1}, {"f2", f2}},
			NextArgs: []string{"f1", "f2"},
			Expected: []func() string{f1, f2},
		},
		{
			Name:     "nofunc",
			Seq:      []namedFunc{{"f1", f1}},
			NextArgs: []string{"f2"},
			Expected: []func() string{nil},
		},
		{
			Name:     "consumed",
			Seq:      []namedFunc{{"f1", f1}},
			NextArgs: []string{"f1", "f1"},
			Expected: []func() string{f1, nil},
		},
		{
			Name:     "set",
			Set:      namedFunc{"f1", f1},
			NextArgs: []string{"f1", "f1", "f1", "f1"},
			Expected: []func() string{f1, f1, f1, f1},
		},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			m := New()
			for _, nf := range c.Seq {
				m.Add(nf.Name, nf.Func)
			}
			if nf := c.Set; nf.Name != "" {
				m.Set(nf.Name, nf.Func)
			}
			for i, a := range c.NextArgs {
				got := m.Next(a)
				compareFuncs(t, got, c.Expected[i])
			}
		})
	}
}

func compareFuncs(t *testing.T, got interface{}, expected func() string) {
	if got == nil && expected == nil {
		return
	}
	f, ok := got.(func() string)
	if !ok {
		t.Errorf("got %v, want func() string", got)
		return
	}
	if f() != expected() {
		t.Errorf("got %s %p, expected %s %p", f(), got, expected(), expected)
	}
}

func f1() string { return "f1" }
func f2() string { return "f2" }
