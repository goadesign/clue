package extensive

import (
	"io"
	"unsafe"

	goa "goa.design/goa/v3/pkg"

	imported "example.com/c/d/extensive/aliased"
)

type (
	Extensive interface {
		Simple(int, string) float64
		NoResult()
		MultipleResults() (bool, complex64, complex128, string, unsafe.Pointer, error)
		NamedResult() (err error)
		RepeatedTypes(a, b int, c, d float64) (e, f int, g, h float64, err error)
		Variadic(args ...string)
		ComplexTypes([5]string, []string, map[string]string, *string, chan int, chan<- int, <-chan int) ([5]string, []string, map[string]string, *string, chan int, chan<- int, <-chan int)
		MoreComplexTypes(interface{}, interface {
			io.ReadWriter
			A(int) error
			B()
		}, struct {
			Struct
			A, B int
			C    float64
		}, func(int) (bool, error)) (interface{}, interface {
			io.ReadWriter
			A(int) error
			B()
		}, struct {
			Struct
			A, B int
			C    float64
		}, func(int) (bool, error))
		NamedTypes(Struct, Array, io.Reader, imported.Type, goa.Endpoint, Generic[uint, string, Struct, Array]) (Struct, Array, io.Reader, imported.Type, goa.Endpoint, Generic[uint, string, Struct, Array])
		FuncNamedTypes(func(Struct, Array, io.Reader, imported.Type, goa.Endpoint, Generic[uint, string, Struct, Array])) func(Struct, Array, io.Reader, imported.Type, goa.Endpoint, Generic[uint, string, Struct, Array])
		VariableConflicts(f, m uint)
		AliasedTypes(IntAlias, ArrayAlias, StructAlias, IntSetAlias, SetAlias[string]) (IntAlias, ArrayAlias, StructAlias, IntSetAlias, SetAlias[string])
		AliasedFuncTypes(func(IntAlias, ArrayAlias, StructAlias, IntSetAlias, SetAlias[string])) func(IntAlias, ArrayAlias, StructAlias, IntSetAlias, SetAlias[string])

		Embedded
		imported.Interface
	}

	Embedded interface {
		Embedded(int8) int8
	}

	Generic[K comparable, V ~int | bool | string, X, Y any] interface {
		Simple(k K, v V, x X, y Y) (K, V, X, Y)
		Complex(map[K]V, []X, *Y, Set[K]) (map[K]V, []X, *Y, Set[K])
	}

	Struct                 struct{ A, B int }
	Array                  [5]Struct
	Set[K comparable]      map[K]Struct
	IntAlias               = int
	ArrayAlias             = [5]Struct
	StructAlias            = struct{ A, B int }
	IntSetAlias            = Set[int]
	SetAlias[K comparable] = Set[K]

	ExtensiveAlias                                               = Extensive
	ImportedAlias                                                = imported.Interface
	GenericAlias[K comparable, V ~int | bool | string, X, Y any] = Generic[K, V, X, Y]
	ConstrainedGenericAlias                                      = Generic[string, IntAlias, float32, float64]
)
