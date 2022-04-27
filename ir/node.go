package ir

type Node struct {
	Op Op

	Args []*Node

	Value interface{}

	Type Type
}

func (n *Node) IsStatement() bool {
	return !miscOpsMap[n.Op] && statementOpsMap[n.Op]
}

func (n *Node) IsExpression() bool {
	return !miscOpsMap[n.Op] && !statementOpsMap[n.Op]
}

//go:generate stringer -type Op -trimprefix Op
type Op int

const (
	// OpInvalid is a marker for unset op values; should never be used.
	OpInvalid Op = iota

	// OpBad is a special node that is not valid for PHP.
	// $Value.(string) contains a text to be inserted "as is".
	OpBad

	// break $Value.(int)
	// A value of 0 means "no explicit label".
	OpBreak

	// continue $Value.(int)
	// A value of 0 means "no explicit label".
	OpContinue

	// 'if' '(' $Args[0] ')' $Args[1]
	OpIf

	// 'if' '(' $Args[0] ')' $Args[1] 'else' $Args[2]
	OpIfElse

	// 'while' '(' $Args[0] ')' $Args[1]
	OpWhile

	// 'do' $Args[0] 'while' $Args[1]
	OpDoWhile

	// '{' $Args[:]... '}'
	OpBlock

	// 'return' $Args[0]
	OpReturn

	// 'return'
	OpReturnVoid

	// 'echo' $Args[:]...
	OpEcho

	// '(' $Args[0] ')'
	OpParens

	// $Args[0] '=' $Args[1]
	OpAssign

	// $Args[0] <op>'=' $Args[1]
	// $Value.(Op) contains the operation (like OpAdd)
	OpAssignModify

	// $Value.(bool)
	OpBoolLit
	// $Value.(int64)
	OpIntLit
	// $Value.(float64)
	OpFloatLit
	// $Value.(string)
	OpStringLit

	// $Value.(string) contains a variable name
	// $Type contains a variable type
	OpVar

	// $Value.(string) contains a symbol name
	OpName

	// '!' $Args[0]
	OpNot

	// $Args[0] '->' $Value.(string)
	OpProp

	// $Args[0] '[' $Args[1] ']'
	OpIndex

	// $Args[0] '.' $Args[1]
	OpConcat
	// $Args[0] '+' $Args[1]
	OpAdd
	// $Args[0] '-' $Args[1]
	OpSub

	// $Args[0] '&&' $Args[1]
	OpAnd
	// $Args[0] '||' $Args[1]
	OpOr

	// $Args[0] '?' $Args[1] ':' $Args[2]
	OpTernary

	// $Args[0] '(' $Args[1:]... ')'
	OpCall
)

var statementOpsMap = [...]bool{
	OpBreak:      true,
	OpContinue:   true,
	OpIf:         true,
	OpIfElse:     true,
	OpWhile:      true,
	OpBlock:      true,
	OpReturn:     true,
	OpReturnVoid: true,
	OpEcho:       true,
}

var miscOpsMap = [...]bool{
	OpInvalid: true,
}

func NewBreak(value int) *Node {
	return &Node{Op: OpBreak, Value: value}
}

func NewContinue(value int) *Node {
	return &Node{Op: OpContinue, Value: value}
}

func NewIf(cond, body *Node) *Node {
	return &Node{Op: OpIf, Args: []*Node{cond, body}}
}

func NewIfElse(cond, body, elseNode *Node) *Node {
	return &Node{Op: OpIfElse, Args: []*Node{cond, body, elseNode}}
}

func NewWhile(cond, body *Node) *Node {
	return &Node{Op: OpWhile, Args: []*Node{cond, body}}
}

func NewDoWhile(body, cond *Node) *Node {
	return &Node{Op: OpDoWhile, Args: []*Node{body, cond}}
}

func NewBlock(statements ...*Node) *Node {
	return &Node{Op: OpBlock, Args: statements}
}

func NewReturn(x *Node) *Node {
	return &Node{Op: OpReturn, Args: []*Node{x}}
}

func NewReturnVoid() *Node {
	return &Node{Op: OpReturnVoid}
}

func NewEcho(args ...*Node) *Node {
	return &Node{Op: OpEcho, Args: args}
}

func NewParens(x *Node) *Node {
	return &Node{Op: OpParens, Args: []*Node{x}}
}

func NewAssign(lhs, rhs *Node) *Node {
	return &Node{Op: OpAssign, Args: []*Node{lhs, rhs}}
}

func NewAssignModify(op Op, rhs, lhs *Node) *Node {
	return &Node{Op: OpAssignModify, Value: op, Args: []*Node{lhs, rhs}}
}

func NewBoolLit(value bool) *Node {
	return &Node{Op: OpBoolLit, Value: value}
}

func NewIntLit(value int64) *Node {
	return &Node{Op: OpIntLit, Value: value}
}

func NewFloatLit(value float64) *Node {
	return &Node{Op: OpFloatLit, Value: value}
}

func NewStringLit(value string) *Node {
	return &Node{Op: OpStringLit, Value: value}
}

func NewVar(name string, typ Type) *Node {
	return &Node{Op: OpVar, Value: name, Type: typ}
}

func NewName(name string) *Node {
	return &Node{Op: OpName, Value: name}
}

func NewNot(x *Node) *Node {
	return &Node{Op: OpNot, Args: []*Node{x}}
}

func NewProp(obj *Node, propName string) *Node {
	return &Node{Op: OpProp, Value: propName, Args: []*Node{obj}}
}

func NewIndex(array, key *Node) *Node {
	return &Node{Op: OpIndex, Args: []*Node{array, key}}
}

func NewConcat(x, y *Node) *Node {
	return &Node{Op: OpConcat, Args: []*Node{x, y}}
}

func NewAdd(x, y *Node) *Node {
	return &Node{Op: OpAdd, Args: []*Node{x, y}}
}

func NewSub(x, y *Node) *Node {
	return &Node{Op: OpSub, Args: []*Node{x, y}}
}

func NewAnd(x, y *Node) *Node {
	return &Node{Op: OpAnd, Args: []*Node{x, y}}
}

func NewOr(x, y *Node) *Node {
	return &Node{Op: OpOr, Args: []*Node{x, y}}
}

func NewTernary(cond, trueExpr, falseExpr *Node) *Node {
	return &Node{Op: OpTernary, Args: []*Node{cond, trueExpr, falseExpr}}
}

func NewCall(fn *Node, args ...*Node) *Node {
	allArgs := make([]*Node, len(args)+1)
	allArgs[0] = fn
	copy(allArgs[1:], args)
	return &Node{Op: OpCall, Args: allArgs}
}
