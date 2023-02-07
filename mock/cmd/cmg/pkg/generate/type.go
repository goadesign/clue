package generate

import (
	"fmt"
	"go/types"
	"strings"
)

type (
	typeMap map[types.Type]string

	typeAdder struct {
		names, zeros                       typeMap
		stdImports, extImports, intImports importMap
		modPath                            string
	}
)

func addType(tt types.Type, typeNames, typeZeros typeMap, stdImports, extImports, intImports importMap, modPath string) (name, zero string) {
	ta := &typeAdder{typeNames, typeZeros, stdImports, extImports, intImports, modPath}
	name = ta.name(tt)
	zero = ta.zero(tt)
	return
}

func (ta *typeAdder) name(tt types.Type) (name string) {
	name, ok := ta.names[tt]
	if !ok {
		switch t := tt.(type) {
		case *types.Array:
			name = fmt.Sprintf("[%v]%v", t.Len(), ta.name(t.Elem()))
		case *types.Basic:
			switch t.Kind() {
			case types.UnsafePointer:
				i := addImport(newImport("unsafe"), ta.stdImports, ta.extImports, ta.intImports, ta.modPath)
				name = fmt.Sprintf("%v.Pointer", i.AliasOrPkgName())
			default:
				name = t.Name()
			}
		case *types.Interface:
			// TODO
		case *types.Map:
			name = fmt.Sprintf("map[%v]%v", ta.name(t.Key()), ta.name(t.Elem()))
		case *types.Named:
			o := t.Obj()
			if p := o.Pkg(); p != nil {
				i := addImport(newImport(p.Path(), p.Name()), ta.stdImports, ta.extImports, ta.intImports, ta.modPath)
				name = fmt.Sprintf("%v.%v", i.AliasOrPkgName(), o.Name())
			} else {
				name = o.Name()
			}
		case *types.Pointer:
			name = "*" + ta.name(t.Elem())
		case *types.Signature:
			// TODO
			b := &strings.Builder{}
			name = b.String()
		case *types.Slice:
			name = "[]" + ta.name(t.Elem())
		case *types.Struct:
			// TODO
		default:
			panic(fmt.Sprintf("unknown name for type: %#v (%T)", t, t))
		}
		ta.names[tt] = name
	}
	return
}

func (ta *typeAdder) zero(tt types.Type) (zero string) {
	zero, ok := ta.zeros[tt]
	if !ok {
		switch t := tt.(type) {
		case *types.Array, *types.Struct:
			zero = ta.name(t) + "{}"
		case *types.Basic:
			switch t.Kind() {
			case types.Bool:
				zero = "false"
			case types.Complex64, types.Complex128:
				zero = "0+0i"
			case types.String:
				zero = `""`
			case types.UnsafePointer:
				zero = "nil"
			default:
				zero = "0"
			}
		case *types.Chan, *types.Interface, *types.Map, *types.Pointer, *types.Signature, *types.Slice:
			zero = "nil"
		case *types.Named:
			zero = ta.zero(t.Underlying())
		default:
			panic(fmt.Sprintf("unknown zero for type: %#v (%T)", t, t))
		}
		ta.zeros[tt] = zero
	}
	return
}