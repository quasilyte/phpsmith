package irgen

import "github.com/quasilyte/phpsmith/ir"

func isSimpleNode(n *ir.Node) bool {
	switch n.Op {
	case ir.OpVar, ir.OpName, ir.OpParens, ir.OpStringLit, ir.OpBoolLit, ir.OpCall, ir.OpArrayLit, ir.OpCast:
		return true
	case ir.OpIntLit:
		return n.Value.(int64) >= 0
	case ir.OpFloatLit:
		return n.Value.(float64) >= 0
	default:
		return false
	}
}
