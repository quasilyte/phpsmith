package irgen

import (
	"fmt"

	"github.com/quasilyte/phpsmith/ir"
)

func typesIdentical(t1, t2 ir.Type) bool {
	switch t1 := t1.(type) {
	case *ir.ScalarType:
		t2, ok := t2.(*ir.ScalarType)
		return ok && t1.Kind == t2.Kind

	default:
		panic(fmt.Sprintf("unexpected type %T", t1))
	}
}
