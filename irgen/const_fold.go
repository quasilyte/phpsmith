package irgen

import "github.com/quasilyte/phpsmith/ir"

func extractValue(n *ir.Node) any {
	switch n.Op {
	case ir.OpIntLit, ir.OpStringLit, ir.OpFloatLit, ir.OpBoolLit:
		return n.Value
	default:
		return nil
	}
}
