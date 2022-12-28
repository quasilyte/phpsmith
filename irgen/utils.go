package irgen

import "github.com/quasilyte/phpsmith/ir"

func canBeNull(n *ir.Node) bool {
	switch n.Op {
	case ir.OpParens:
		return canBeNull(n.Args[0])
	case ir.OpNew:
		return false
	case ir.OpVar:
		return n.Value.(string) != "this"
	default:
		return true
	}
}

func newSimpleCall(fn string, args ...*ir.Node) *ir.Node {
	return ir.NewCall(ir.NewName(fn), args...)
}

func isSimpleNode(n *ir.Node) bool {
	switch n.Op {
	case ir.OpVar, ir.OpName, ir.OpParens, ir.OpStringLit, ir.OpBoolLit, ir.OpCall, ir.OpArrayLit, ir.OpMemberAccess:
		return true
	case ir.OpIntLit:
		return n.Value.(int64) >= 0
	case ir.OpFloatLit:
		return n.Value.(float64) >= 0
	default:
		return false
	}
}
