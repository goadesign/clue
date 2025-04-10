// Code generated by Clue Mock Generator TEST VERSION, DO NOT EDIT.
//
// Command:
// $ cmg gen example.com/c/d/extensive

package mockextensive

import (
	"io"
	"testing"
	"unsafe"

	"goa.design/clue/mock"
	goa "goa.design/goa/v3/pkg"

	"example.com/c/d/extensive"
	imported "example.com/c/d/extensive/aliased"
)

type (
	Extensive struct {
		m *mock.Mock
		t *testing.T
	}

	ExtensiveSimpleFunc            func(p0 int, p1 string) float64
	ExtensiveNoResultFunc          func()
	ExtensiveMultipleResultsFunc   func() (bool, complex64, complex128, string, unsafe.Pointer, error)
	ExtensiveNamedResultFunc       func() (err error)
	ExtensiveRepeatedTypesFunc     func(a, b int, c, d float64) (e, f int, g, h float64, err error)
	ExtensiveVariadicFunc          func(args ...string)
	ExtensiveComplexTypesFunc      func(p0 [5]string, p1 []string, p2 map[string]string, p3 *string, p4 chan int, p5 chan<- int, p6 <-chan int) ([5]string, []string, map[string]string, *string, chan int, chan<- int, <-chan int)
	ExtensiveMoreComplexTypesFunc  func(p0 interface{}, p1 interface{io.ReadWriter; A(int) error; B()}, p2 struct{extensive.Struct; A, B int; C float64}, p3 func(int) (bool, error)) (interface{}, interface{io.ReadWriter; A(int) error; B()}, struct{extensive.Struct; A, B int; C float64}, func(int) (bool, error))
	ExtensiveNamedTypesFunc        func(p0 extensive.Struct, p1 extensive.Array, p2 io.Reader, p3 imported.Type, p4 goa.Endpoint, p5 extensive.Generic[uint, string, extensive.Struct, extensive.Array]) (extensive.Struct, extensive.Array, io.Reader, imported.Type, goa.Endpoint, extensive.Generic[uint, string, extensive.Struct, extensive.Array])
	ExtensiveFuncNamedTypesFunc    func(p0 func(extensive.Struct, extensive.Array, io.Reader, imported.Type, goa.Endpoint, extensive.Generic[uint, string, extensive.Struct, extensive.Array])) func(extensive.Struct, extensive.Array, io.Reader, imported.Type, goa.Endpoint, extensive.Generic[uint, string, extensive.Struct, extensive.Array])
	ExtensiveVariableConflictsFunc func(f, m uint)
	ExtensiveAliasedTypesFunc      func(p0 extensive.IntAlias, p1 extensive.ArrayAlias, p2 extensive.StructAlias, p3 extensive.IntSetAlias, p4 extensive.SetAlias[string]) (extensive.IntAlias, extensive.ArrayAlias, extensive.StructAlias, extensive.IntSetAlias, extensive.SetAlias[string])
	ExtensiveAliasedFuncTypesFunc  func(p0 func(extensive.IntAlias, extensive.ArrayAlias, extensive.StructAlias, extensive.IntSetAlias, extensive.SetAlias[string])) func(extensive.IntAlias, extensive.ArrayAlias, extensive.StructAlias, extensive.IntSetAlias, extensive.SetAlias[string])
	ExtensiveEmbeddedFunc          func(p0 int8) int8
	ExtensiveImportedFunc          func(p0 imported.Type) imported.Type

	Embedded struct {
		m *mock.Mock
		t *testing.T
	}

	EmbeddedEmbeddedFunc func(p0 int8) int8

	Generic[K comparable, V ~int | bool | string, X, Y any] struct {
		m *mock.Mock
		t *testing.T
	}

	GenericSimpleFunc[K comparable, V ~int | bool | string, X, Y any]  func(k K, v V, x X, y Y) (K, V, X, Y)
	GenericComplexFunc[K comparable, V ~int | bool | string, X, Y any] func(p0 map[K]V, p1 []X, p2 *Y, p3 extensive.Set[K]) (map[K]V, []X, *Y, extensive.Set[K])
)

func NewExtensive(t *testing.T) *Extensive {
	var (
		m                     = &Extensive{mock.New(), t}
		_ extensive.Extensive = m
	)
	return m
}

func (m *Extensive) AddSimple(f ExtensiveSimpleFunc) {
	m.m.Add("Simple", f)
}

func (m *Extensive) SetSimple(f ExtensiveSimpleFunc) {
	m.m.Set("Simple", f)
}

func (m *Extensive) Simple(p0 int, p1 string) float64 {
	if f := m.m.Next("Simple"); f != nil {
		return f.(ExtensiveSimpleFunc)(p0, p1)
	}
	m.t.Helper()
	m.t.Error("unexpected Simple call")
	return 0
}

func (m *Extensive) AddNoResult(f ExtensiveNoResultFunc) {
	m.m.Add("NoResult", f)
}

func (m *Extensive) SetNoResult(f ExtensiveNoResultFunc) {
	m.m.Set("NoResult", f)
}

func (m *Extensive) NoResult() {
	if f := m.m.Next("NoResult"); f != nil {
		f.(ExtensiveNoResultFunc)()
		return
	}
	m.t.Helper()
	m.t.Error("unexpected NoResult call")
}

func (m *Extensive) AddMultipleResults(f ExtensiveMultipleResultsFunc) {
	m.m.Add("MultipleResults", f)
}

func (m *Extensive) SetMultipleResults(f ExtensiveMultipleResultsFunc) {
	m.m.Set("MultipleResults", f)
}

func (m *Extensive) MultipleResults() (bool, complex64, complex128, string, unsafe.Pointer, error) {
	if f := m.m.Next("MultipleResults"); f != nil {
		return f.(ExtensiveMultipleResultsFunc)()
	}
	m.t.Helper()
	m.t.Error("unexpected MultipleResults call")
	return false, 0i, 0i, "", nil, nil
}

func (m *Extensive) AddNamedResult(f ExtensiveNamedResultFunc) {
	m.m.Add("NamedResult", f)
}

func (m *Extensive) SetNamedResult(f ExtensiveNamedResultFunc) {
	m.m.Set("NamedResult", f)
}

func (m *Extensive) NamedResult() (err error) {
	if f := m.m.Next("NamedResult"); f != nil {
		return f.(ExtensiveNamedResultFunc)()
	}
	m.t.Helper()
	m.t.Error("unexpected NamedResult call")
	return nil
}

func (m *Extensive) AddRepeatedTypes(f ExtensiveRepeatedTypesFunc) {
	m.m.Add("RepeatedTypes", f)
}

func (m *Extensive) SetRepeatedTypes(f ExtensiveRepeatedTypesFunc) {
	m.m.Set("RepeatedTypes", f)
}

func (m *Extensive) RepeatedTypes(a, b int, c, d float64) (e, f int, g, h float64, err error) {
	if f := m.m.Next("RepeatedTypes"); f != nil {
		return f.(ExtensiveRepeatedTypesFunc)(a, b, c, d)
	}
	m.t.Helper()
	m.t.Error("unexpected RepeatedTypes call")
	return 0, 0, 0, 0, nil
}

func (m *Extensive) AddVariadic(f ExtensiveVariadicFunc) {
	m.m.Add("Variadic", f)
}

func (m *Extensive) SetVariadic(f ExtensiveVariadicFunc) {
	m.m.Set("Variadic", f)
}

func (m *Extensive) Variadic(args ...string) {
	if f := m.m.Next("Variadic"); f != nil {
		f.(ExtensiveVariadicFunc)(args...)
		return
	}
	m.t.Helper()
	m.t.Error("unexpected Variadic call")
}

func (m *Extensive) AddComplexTypes(f ExtensiveComplexTypesFunc) {
	m.m.Add("ComplexTypes", f)
}

func (m *Extensive) SetComplexTypes(f ExtensiveComplexTypesFunc) {
	m.m.Set("ComplexTypes", f)
}

func (m *Extensive) ComplexTypes(p0 [5]string, p1 []string, p2 map[string]string, p3 *string, p4 chan int, p5 chan<- int, p6 <-chan int) ([5]string, []string, map[string]string, *string, chan int, chan<- int, <-chan int) {
	if f := m.m.Next("ComplexTypes"); f != nil {
		return f.(ExtensiveComplexTypesFunc)(p0, p1, p2, p3, p4, p5, p6)
	}
	m.t.Helper()
	m.t.Error("unexpected ComplexTypes call")
	return [5]string{}, nil, nil, nil, nil, nil, nil
}

func (m *Extensive) AddMoreComplexTypes(f ExtensiveMoreComplexTypesFunc) {
	m.m.Add("MoreComplexTypes", f)
}

func (m *Extensive) SetMoreComplexTypes(f ExtensiveMoreComplexTypesFunc) {
	m.m.Set("MoreComplexTypes", f)
}

func (m *Extensive) MoreComplexTypes(p0 interface{}, p1 interface{io.ReadWriter; A(int) error; B()}, p2 struct{extensive.Struct; A, B int; C float64}, p3 func(int) (bool, error)) (interface{}, interface{io.ReadWriter; A(int) error; B()}, struct{extensive.Struct; A, B int; C float64}, func(int) (bool, error)) {
	if f := m.m.Next("MoreComplexTypes"); f != nil {
		return f.(ExtensiveMoreComplexTypesFunc)(p0, p1, p2, p3)
	}
	m.t.Helper()
	m.t.Error("unexpected MoreComplexTypes call")
	return nil, nil, struct{extensive.Struct; A, B int; C float64}{}, nil
}

func (m *Extensive) AddNamedTypes(f ExtensiveNamedTypesFunc) {
	m.m.Add("NamedTypes", f)
}

func (m *Extensive) SetNamedTypes(f ExtensiveNamedTypesFunc) {
	m.m.Set("NamedTypes", f)
}

func (m *Extensive) NamedTypes(p0 extensive.Struct, p1 extensive.Array, p2 io.Reader, p3 imported.Type, p4 goa.Endpoint, p5 extensive.Generic[uint, string, extensive.Struct, extensive.Array]) (extensive.Struct, extensive.Array, io.Reader, imported.Type, goa.Endpoint, extensive.Generic[uint, string, extensive.Struct, extensive.Array]) {
	if f := m.m.Next("NamedTypes"); f != nil {
		return f.(ExtensiveNamedTypesFunc)(p0, p1, p2, p3, p4, p5)
	}
	m.t.Helper()
	m.t.Error("unexpected NamedTypes call")
	return extensive.Struct{}, extensive.Array{}, nil, 0, nil, nil
}

func (m *Extensive) AddFuncNamedTypes(f ExtensiveFuncNamedTypesFunc) {
	m.m.Add("FuncNamedTypes", f)
}

func (m *Extensive) SetFuncNamedTypes(f ExtensiveFuncNamedTypesFunc) {
	m.m.Set("FuncNamedTypes", f)
}

func (m *Extensive) FuncNamedTypes(p0 func(extensive.Struct, extensive.Array, io.Reader, imported.Type, goa.Endpoint, extensive.Generic[uint, string, extensive.Struct, extensive.Array])) func(extensive.Struct, extensive.Array, io.Reader, imported.Type, goa.Endpoint, extensive.Generic[uint, string, extensive.Struct, extensive.Array]) {
	if f := m.m.Next("FuncNamedTypes"); f != nil {
		return f.(ExtensiveFuncNamedTypesFunc)(p0)
	}
	m.t.Helper()
	m.t.Error("unexpected FuncNamedTypes call")
	return nil
}

func (m *Extensive) AddVariableConflicts(f ExtensiveVariableConflictsFunc) {
	m.m.Add("VariableConflicts", f)
}

func (m *Extensive) SetVariableConflicts(f ExtensiveVariableConflictsFunc) {
	m.m.Set("VariableConflicts", f)
}

func (m1 *Extensive) VariableConflicts(f, m uint) {
	if f1 := m1.m.Next("VariableConflicts"); f1 != nil {
		f1.(ExtensiveVariableConflictsFunc)(f, m)
		return
	}
	m1.t.Helper()
	m1.t.Error("unexpected VariableConflicts call")
}

func (m *Extensive) AddAliasedTypes(f ExtensiveAliasedTypesFunc) {
	m.m.Add("AliasedTypes", f)
}

func (m *Extensive) SetAliasedTypes(f ExtensiveAliasedTypesFunc) {
	m.m.Set("AliasedTypes", f)
}

func (m *Extensive) AliasedTypes(p0 extensive.IntAlias, p1 extensive.ArrayAlias, p2 extensive.StructAlias, p3 extensive.IntSetAlias, p4 extensive.SetAlias[string]) (extensive.IntAlias, extensive.ArrayAlias, extensive.StructAlias, extensive.IntSetAlias, extensive.SetAlias[string]) {
	if f := m.m.Next("AliasedTypes"); f != nil {
		return f.(ExtensiveAliasedTypesFunc)(p0, p1, p2, p3, p4)
	}
	m.t.Helper()
	m.t.Error("unexpected AliasedTypes call")
	return 0, extensive.ArrayAlias{}, extensive.StructAlias{}, nil, nil
}

func (m *Extensive) AddAliasedFuncTypes(f ExtensiveAliasedFuncTypesFunc) {
	m.m.Add("AliasedFuncTypes", f)
}

func (m *Extensive) SetAliasedFuncTypes(f ExtensiveAliasedFuncTypesFunc) {
	m.m.Set("AliasedFuncTypes", f)
}

func (m *Extensive) AliasedFuncTypes(p0 func(extensive.IntAlias, extensive.ArrayAlias, extensive.StructAlias, extensive.IntSetAlias, extensive.SetAlias[string])) func(extensive.IntAlias, extensive.ArrayAlias, extensive.StructAlias, extensive.IntSetAlias, extensive.SetAlias[string]) {
	if f := m.m.Next("AliasedFuncTypes"); f != nil {
		return f.(ExtensiveAliasedFuncTypesFunc)(p0)
	}
	m.t.Helper()
	m.t.Error("unexpected AliasedFuncTypes call")
	return nil
}

func (m *Extensive) AddEmbedded(f ExtensiveEmbeddedFunc) {
	m.m.Add("Embedded", f)
}

func (m *Extensive) SetEmbedded(f ExtensiveEmbeddedFunc) {
	m.m.Set("Embedded", f)
}

func (m *Extensive) Embedded(p0 int8) int8 {
	if f := m.m.Next("Embedded"); f != nil {
		return f.(ExtensiveEmbeddedFunc)(p0)
	}
	m.t.Helper()
	m.t.Error("unexpected Embedded call")
	return 0
}

func (m *Extensive) AddImported(f ExtensiveImportedFunc) {
	m.m.Add("Imported", f)
}

func (m *Extensive) SetImported(f ExtensiveImportedFunc) {
	m.m.Set("Imported", f)
}

func (m *Extensive) Imported(p0 imported.Type) imported.Type {
	if f := m.m.Next("Imported"); f != nil {
		return f.(ExtensiveImportedFunc)(p0)
	}
	m.t.Helper()
	m.t.Error("unexpected Imported call")
	return 0
}

func (m *Extensive) HasMore() bool {
	return m.m.HasMore()
}

func NewEmbedded(t *testing.T) *Embedded {
	var (
		m                    = &Embedded{mock.New(), t}
		_ extensive.Embedded = m
	)
	return m
}

func (m *Embedded) AddEmbedded(f EmbeddedEmbeddedFunc) {
	m.m.Add("Embedded", f)
}

func (m *Embedded) SetEmbedded(f EmbeddedEmbeddedFunc) {
	m.m.Set("Embedded", f)
}

func (m *Embedded) Embedded(p0 int8) int8 {
	if f := m.m.Next("Embedded"); f != nil {
		return f.(EmbeddedEmbeddedFunc)(p0)
	}
	m.t.Helper()
	m.t.Error("unexpected Embedded call")
	return 0
}

func (m *Embedded) HasMore() bool {
	return m.m.HasMore()
}

func NewGeneric[K comparable, V ~int | bool | string, X, Y any](t *testing.T) *Generic[K, V, X, Y] {
	var (
		m                               = &Generic[K, V, X, Y]{mock.New(), t}
		_ extensive.Generic[K, V, X, Y] = m
	)
	return m
}

func (m *Generic[K, V, X, Y]) AddSimple(f GenericSimpleFunc[K, V, X, Y]) {
	m.m.Add("Simple", f)
}

func (m *Generic[K, V, X, Y]) SetSimple(f GenericSimpleFunc[K, V, X, Y]) {
	m.m.Set("Simple", f)
}

func (m *Generic[K, V, X, Y]) Simple(k K, v V, x X, y Y) (K, V, X, Y) {
	if f := m.m.Next("Simple"); f != nil {
		return f.(GenericSimpleFunc[K, V, X, Y])(k, v, x, y)
	}
	m.t.Helper()
	m.t.Error("unexpected Simple call")
	return *new(K), *new(V), *new(X), *new(Y)
}

func (m *Generic[K, V, X, Y]) AddComplex(f GenericComplexFunc[K, V, X, Y]) {
	m.m.Add("Complex", f)
}

func (m *Generic[K, V, X, Y]) SetComplex(f GenericComplexFunc[K, V, X, Y]) {
	m.m.Set("Complex", f)
}

func (m *Generic[K, V, X, Y]) Complex(p0 map[K]V, p1 []X, p2 *Y, p3 extensive.Set[K]) (map[K]V, []X, *Y, extensive.Set[K]) {
	if f := m.m.Next("Complex"); f != nil {
		return f.(GenericComplexFunc[K, V, X, Y])(p0, p1, p2, p3)
	}
	m.t.Helper()
	m.t.Error("unexpected Complex call")
	return nil, nil, nil, nil
}

func (m *Generic[K, V, X, Y]) HasMore() bool {
	return m.m.HasMore()
}
