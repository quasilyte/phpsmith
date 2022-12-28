package irgen

import (
	"fmt"

	"github.com/quasilyte/phpsmith/ir"
)

func canConstexprInitialize(t ir.Type) bool {
	switch t := t.(type) {
	case *ir.ScalarType:
		return t.Kind != ir.ScalarFloat
	case *ir.EnumType:
		return t.ValueType.Kind != ir.ScalarFloat
	case *ir.ClassType:
		return true
	case *ir.ArrayType:
		return canConstexprInitialize(t.Elem)
	default:
		return false
	}
}

func canDump(t ir.Type) bool {
	switch t := t.(type) {
	case *ir.ScalarType, *ir.EnumType:
		return true
	case *ir.ArrayType:
		return canDump(t.Elem)
	default:
		return false
	}
}

func typeLess(t1, t2 ir.Type) bool {
	tag1 := t1.Tag()
	tag2 := t2.Tag()
	if tag1 < tag2 {
		return true
	}
	if tag1 > tag2 {
		return false
	}

	// If tags are identical, the underlying concrete types are identical too.
	switch t1 := t1.(type) {
	case *ir.ScalarType:
		return t1.Kind < t2.(*ir.ScalarType).Kind
	case *ir.ClassType:
		return t1.Name < t2.(*ir.ClassType).Name
	case *ir.FuncType:
		return t1.Name < t2.(*ir.FuncType).Name
	case *ir.ArrayType:
		return typeLess(t1.Elem, t2.(*ir.ArrayType).Elem)
	case *ir.TupleType:
		t2 := t2.(*ir.TupleType)
		if len(t1.Elems) < len(t2.Elems) {
			return true
		}
		if len(t1.Elems) > len(t2.Elems) {
			return false
		}
		for i, e1 := range t1.Elems {
			e2 := t2.Elems[i]
			if typeLess(e1, e2) {
				return true
			}
		}
		return false
	case *ir.EnumType:
		t2 := t2.(*ir.EnumType)
		if len(t1.Values) < len(t2.Values) || t1.ValueType.Kind < t2.ValueType.Kind {
			return true
		}
		if len(t1.Values) > len(t2.Values) {
			return false
		}
		for i, v1 := range t1.Values {
			v2 := t2.Values[i]
			switch v1 := v1.(type) {
			case string:
				if v1 < v2.(string) {
					return true
				}
			case int64:
				if v1 < v2.(int64) {
					return true
				}
			case float64:
				if v1 < v2.(float64) {
					return true
				}
			case bool:
				if !v1 && v2.(bool) {
					return true
				}
			default:
				panic(fmt.Sprintf("unexpected enum value: %T", v1))
			}
		}
		return false

	default:
		panic(fmt.Sprintf("unexpected type %T", t1))
	}
}

func typesIdentical(t1, t2 ir.Type) bool {
	switch t1 := t1.(type) {
	case *ir.ScalarType:
		t2, ok := t2.(*ir.ScalarType)
		return ok && t1.Kind == t2.Kind

	case *ir.ClassType:
		t2, ok := t2.(*ir.ClassType)
		return ok && t1.Name == t2.Name

	case *ir.FuncType:
		t2, ok := t2.(*ir.FuncType)
		return ok && t1.Name == t2.Name

	case *ir.ArrayType:
		t2, ok := t2.(*ir.ArrayType)
		return ok && typesIdentical(t1.Elem, t2.Elem)

	case *ir.TupleType:
		t2, ok := t2.(*ir.TupleType)
		if !ok || len(t1.Elems) != len(t2.Elems) {
			return false
		}
		for i, e1 := range t1.Elems {
			e2 := t2.Elems[i]
			if !typesIdentical(e1, e2) {
				return false
			}
		}
		return true

	case *ir.EnumType:
		t2, ok := t2.(*ir.EnumType)
		if !ok || len(t1.Values) != len(t2.Values) || !typesIdentical(t1.ValueType, t2.ValueType) {
			return false
		}
		for i, v1 := range t1.Values {
			v2 := t2.Values[i]
			if v1 != v2 {
				return false
			}
		}
		return true

	default:
		panic(fmt.Sprintf("unexpected type %T", t1))
	}
}
