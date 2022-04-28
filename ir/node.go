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

	// $Args holds array elements
	OpArrayLit

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

	// '-' $Args[0]
	OpNegation

	// '+' $Args[0]
	OpUnaryPlus

	// $Args[0] '.' $Args[1]
	OpConcat

	// $Args[0] '+' $Args[1]
	OpAdd

	// $Args[0] '-' $Args[1]
	OpSub

	// $Args[0] '/' $Args[1]
	OpDiv

	// $Args[0] '*' $Args[1]
	OpMul

	// $Args[0] '%' $Args[1]
	OpMod

	// $Args[0] '**' $Args[1]
	OpExp

	// $Args[0] '&&' $Args[1]
	OpAnd

	// $Args[0] 'and' $Args[1]
	OpAndWord

	// $Args[0] '||' $Args[1]
	OpOr

	// $Args[0] 'or' $Args[1]
	OpOrWord

	// $Args[0] 'xor' $Args[1]
	OpXorWord

	// $Args[0] '?' $Args[1] ':' $Args[2]
	OpTernary

	// $Args[0] '(' $Args[1:]... ')'
	OpCall

	// $Args[0] '<' $Args[1]
	OpLess

	// $Args[0] '<=' $Args[1]
	OpLessOrEqual

	// $Args[0] '>' $Args[1]
	OpGreater

	// $Args[0] '>=' $Args[1]
	OpGreaterOrEqual

	// $Args[0] '==' $Args[1]
	OpEqual2

	// $Args[0] '===' $Args[1]
	OpEqual3

	// $Args[0] '!=' $Args[1]
	OpNotEqual2

	// $Args[0] '!==' $Args[1]
	OpNotEqual3

	// $Args[0] '<=>' $Args[1]
	OpSpaceship

	// $Args[0] '++'
	OpPostInc

	// '++' $Args[0]
	OpPreInc

	// $Args[0] '--'
	OpPostDec

	// '--' $Args[0]
	OpPreDec

	// ($Type)$Args[0]
	OpCast

	// $Args[0] & $Args[1]
	OpBitAnd

	// $Args[0] | $Args[1]
	OpBitOr

	// $Args[0] ^ $Args[1]
	OpBitXor

	// ~ $Args[0]
	OpBitNot

	// $Args[0] << $Args[1]
	OpBitShiftLeft

	// $Args[0] >> $Args[1]
	OpBitShiftRight

	// $Args[0] ?? $Args[1]
	OpNullCoales
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

func NewAssignModify(op Op, lhs, rhs *Node) *Node {
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

func NewLess(x, y *Node) *Node {
	return &Node{Op: OpLess, Args: []*Node{x, y}}
}

func NewLessOrEqual(x, y *Node) *Node {
	return &Node{Op: OpLessOrEqual, Args: []*Node{x, y}}
}

func NewGreater(x, y *Node) *Node {
	return &Node{Op: OpGreater, Args: []*Node{x, y}}
}

func NewGreaterOrEqual(x, y *Node) *Node {
	return &Node{Op: OpGreaterOrEqual, Args: []*Node{x, y}}
}

func NewEqual2(x, y *Node) *Node {
	return &Node{Op: OpEqual2, Args: []*Node{x, y}}
}

func NewEqual3(x, y *Node) *Node {
	return &Node{Op: OpEqual3, Args: []*Node{x, y}}
}

func NewPostInc(x *Node) *Node {
	return &Node{Op: OpPostInc, Args: []*Node{x}}
}

func NewPostDec(x *Node) *Node {
	return &Node{Op: OpPostDec, Args: []*Node{x}}
}

func NewPreInc(x *Node) *Node {
	return &Node{Op: OpPreInc, Args: []*Node{x}}
}

func NewPreDec(x *Node) *Node {
	return &Node{Op: OpPreDec, Args: []*Node{x}}
}

func NewAndWord(x, y *Node) *Node {
	return &Node{Op: OpAndWord, Args: []*Node{x, y}}
}

func NewOrWord(x, y *Node) *Node {
	return &Node{Op: OpOrWord, Args: []*Node{x, y}}
}

func NewXorWord(x, y *Node) *Node {
	return &Node{Op: OpXorWord, Args: []*Node{x, y}}
}

func NewNotEqual2(x, y *Node) *Node {
	return &Node{Op: OpNotEqual2, Args: []*Node{x, y}}
}

func NewNotEqual3(x, y *Node) *Node {
	return &Node{Op: OpNotEqual3, Args: []*Node{x, y}}
}

func NewSpaceship(x, y *Node) *Node {
	return &Node{Op: OpSpaceship, Args: []*Node{x, y}}
}

func NewMod(x, y *Node) *Node {
	return &Node{Op: OpMod, Args: []*Node{x, y}}
}

func NewExp(x, y *Node) *Node {
	return &Node{Op: OpExp, Args: []*Node{x, y}}
}

func NewMul(x, y *Node) *Node {
	return &Node{Op: OpMul, Args: []*Node{x, y}}
}

func NewDiv(x, y *Node) *Node {
	return &Node{Op: OpDiv, Args: []*Node{x, y}}
}

func NewNegation(x *Node) *Node {
	return &Node{Op: OpNegation, Args: []*Node{x}}
}

func NewUnaryPlus(x *Node) *Node {
	return &Node{Op: OpUnaryPlus, Args: []*Node{x}}
}

func NewBitAnd(x, y *Node) *Node {
	return &Node{Op: OpBitAnd, Args: []*Node{x, y}}
}

func NewBitOr(x, y *Node) *Node {
	return &Node{Op: OpBitOr, Args: []*Node{x, y}}
}

func NewBitXor(x, y *Node) *Node {
	return &Node{Op: OpBitXor, Args: []*Node{x, y}}
}

func NewBitNot(x *Node) *Node {
	return &Node{Op: OpBitNot, Args: []*Node{x}}
}

func NewBitShiftLeft(x, y *Node) *Node {
	return &Node{Op: OpBitShiftLeft, Args: []*Node{x, y}}
}

func NewBitShiftRight(x, y *Node) *Node {
	return &Node{Op: OpBitShiftRight, Args: []*Node{x, y}}
}

func NewNullCoales(x, y *Node) *Node {
	return &Node{Op: OpNullCoales, Args: []*Node{x, y}}
}
