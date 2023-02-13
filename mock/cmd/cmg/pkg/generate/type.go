package generate

import (
	"bytes"
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
		case *types.Chan:
			switch t.Dir() {
			case types.SendRecv:
				name = "chan " + ta.name(t.Elem())
			case types.SendOnly:
				name = "chan<- " + ta.name(t.Elem())
			case types.RecvOnly:
				name = "<-chan " + ta.name(t.Elem())
			}
		case *types.Interface:
			if t.Empty() {
				name = t.String()
			} else {
				es := make([]string, 0, t.NumEmbeddeds()+t.NumExplicitMethods())
				for i := 0; i < t.NumEmbeddeds(); i++ {
					es = append(es, ta.name(t.EmbeddedType(i)))
				}
				for i := 0; i < t.NumExplicitMethods(); i++ {
					m := t.ExplicitMethod(i)
					es = append(es, fmt.Sprintf("%v%v", m.Name(), ta.name(m.Type())[4:]))
				}
				name = fmt.Sprintf("interface{%v}", strings.Join(es, "; "))
			}
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
			if tas := t.TypeArgs(); tas != nil {
				as := make([]string, 0, tas.Len())
				for i := 0; i < tas.Len(); i++ {
					as = append(as, ta.name(tas.At(i)))
				}
				name += "[" + strings.Join(as, ", ") + "]"
			}
		case *types.Pointer:
			name = "*" + ta.name(t.Elem())
		case *types.Signature:
			b := &bytes.Buffer{}
			types.WriteSignature(b, t, nil)
			name = "func" + b.String()
		case *types.Slice:
			name = "[]" + ta.name(t.Elem())
		case *types.Struct:
			fs := make([]string, 0, t.NumFields())
			for i := 0; i < t.NumFields(); i++ {
				f := t.Field(i)
				if f.Embedded() {
					fs = append(fs, ta.name(f.Type()))
				} else {
					fs = append(fs, fmt.Sprintf("%v %v", f.Name(), ta.name(f.Type())))
				}
			}
			name = fmt.Sprintf("struct{%v}", strings.Join(fs, "; "))
		case *types.TypeParam:
			name = t.Obj().Name()
		case *types.Union:
			ts := make([]string, 0, t.Len())
			for i := 0; i < t.Len(); i++ {
				term := t.Term(i)
				if term.Tilde() {
					ts = append(ts, "~"+ta.name(term.Type()))
				} else {
					ts = append(ts, ta.name(term.Type()))
				}
			}
			name = strings.Join(ts, " | ")
		default:
			panic(fmt.Errorf("unknown name for type: %#v (%T)", t, t))
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
				zero = "0i"
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
			switch t.Underlying().(type) {
			case *types.Array, *types.Struct:
				zero = ta.name(t) + "{}"
			default:
				zero = ta.zero(t.Underlying())
			}
		case *types.TypeParam:
			zero = fmt.Sprintf("*new(%v)", t.Obj().Name())
		case *types.Union:
			// the zero value for a type constraint union will not be used
		default:
			panic(fmt.Errorf("unknown zero for type: %#v (%T)", t, t))
		}
		ta.zeros[tt] = zero
	}
	return
}
