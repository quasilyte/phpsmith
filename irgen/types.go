package irgen

import (
	"fmt"

	"github.com/quasilyte/phpsmith/ir"
)

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

func typesIdentical(t1, t2 ir.Type) bool {
	switch t1 := t1.(type) {
	case *ir.ScalarType:
		t2, ok := t2.(*ir.ScalarType)
		return ok && t1.Kind == t2.Kind

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
