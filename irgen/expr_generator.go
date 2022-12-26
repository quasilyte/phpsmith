package irgen

import (
	"fmt"
	"math/rand"
	"strings"
	"unicode"

	"github.com/quasilyte/phpsmith/ir"
	"github.com/quasilyte/phpsmith/randutil"
)

type exprGenerator struct {
	config *Config

	rand *rand.Rand

	scope *scope

	symtab *symbolTable

	exprDepth int

	condChoices   exprChoiceList
	boolChoices   exprChoiceList
	intChoices    exprChoiceList
	floatChoices  exprChoiceList
	stringChoices exprChoiceList
}

type exprChoiceList struct {
	indexMap []uint16
	options  []exprChoice
	fallback func() *ir.Node
}

type exprChoice struct {
	freq     int
	generate func() *ir.Node
	fallback func() *ir.Node
}

func newExprGenerator(config *Config, s *scope, symtab *symbolTable) *exprGenerator {
	g := &exprGenerator{
		config: config,
		scope:  s,
		symtab: symtab,
		rand:   config.Rand,
	}

	makeChoicesList := func(fallback func() *ir.Node, options []exprChoice) exprChoiceList {
		indexes := make([]uint16, 0, len(options)*4)
		for i, o := range options {
			for j := 0; j < o.freq; j++ {
				indexes = append(indexes, uint16(i))
			}
		}
		return exprChoiceList{
			indexMap: indexes,
			options:  options,
			fallback: fallback,
		}
	}

	cmpOpGenerator := func(op ir.Op) func() *ir.Node {
		return func() *ir.Node {
			typ := g.PickScalarType()
			x := g.GenerateValueOfType(typ)
			y := g.GenerateValueOfType(typ)
			resultOp := op
			if scalarType, ok := typ.(*ir.ScalarType); ok && scalarType.Kind == ir.ScalarFloat {
				switch resultOp {
				case ir.OpEqual2:
					resultOp = ir.OpFloatEqual2
				case ir.OpEqual3:
					resultOp = ir.OpFloatEqual3
				case ir.OpNotEqual2:
					resultOp = ir.OpNotFloatEqual2
				case ir.OpNotEqual3:
					resultOp = ir.OpNotFloatEqual3
				}
			}
			return &ir.Node{Op: resultOp, Args: []*ir.Node{g.maybeAddParens(x), g.maybeAddParens(y)}}
		}
	}

	binaryOpGenerator := func(op ir.Op, operandGenerator func() *ir.Node) func() *ir.Node {
		return func() *ir.Node {
			x := operandGenerator()
			y := operandGenerator()
			return &ir.Node{Op: op, Args: []*ir.Node{g.maybeAddParens(x), g.maybeAddParens(y)}}
		}
	}

	unaryOpGenerator := func(op ir.Op, operandGenerator func() *ir.Node) func() *ir.Node {
		return func() *ir.Node {
			x := operandGenerator()
			return &ir.Node{Op: op, Args: []*ir.Node{g.maybeAddParens(x)}}
		}
	}

	withCast := func(generator func() *ir.Node, typ ir.Type) func() *ir.Node {
		return func() *ir.Node {
			arg := g.maybeAddParens(generator())
			return &ir.Node{Op: ir.OpCast, Args: []*ir.Node{arg}, Type: typ}
		}
	}

	g.condChoices = makeChoicesList(g.boolLit, []exprChoice{
		{freq: 3, generate: cmpOpGenerator(ir.OpEqual2)},
		{freq: 3, generate: cmpOpGenerator(ir.OpEqual3)},
		{freq: 4, generate: binaryOpGenerator(ir.OpAnd, g.boolValue)},
		{freq: 4, generate: binaryOpGenerator(ir.OpOr, g.boolValue)},
		{freq: 4, generate: unaryOpGenerator(ir.OpNot, g.condValue)},
		{freq: 5, generate: g.boolVar, fallback: g.boolLit},
		{freq: 6, generate: g.boolCall},
		{freq: 1, generate: g.boolLit},
	})

	g.boolChoices = makeChoicesList(g.boolLit, []exprChoice{
		{freq: 1, generate: cmpOpGenerator(ir.OpEqual2)},
		{freq: 1, generate: cmpOpGenerator(ir.OpEqual3)},
		{freq: 3, generate: binaryOpGenerator(ir.OpAnd, g.boolValue)},
		{freq: 3, generate: binaryOpGenerator(ir.OpOr, g.boolValue)},
		{freq: 4, generate: unaryOpGenerator(ir.OpNot, g.condValue)},
		{freq: 6, generate: g.boolVar, fallback: g.boolLit},
		{freq: 3, generate: g.boolLit},
		{freq: 4, generate: g.boolCall},
	})

	g.intChoices = makeChoicesList(g.intLit, []exprChoice{
		{freq: 1, generate: g.intTernary},
		{freq: 2, generate: withCast(binaryOpGenerator(ir.OpAdd, g.intValue), ir.IntType)},
		{freq: 2, generate: binaryOpGenerator(ir.OpSub, g.intValue)},
		{freq: 1, generate: withCast(binaryOpGenerator(ir.OpMul, g.intValue), ir.IntType)},
		{freq: 1, generate: binaryOpGenerator(ir.OpBitAnd, g.intValue)},
		{freq: 1, generate: binaryOpGenerator(ir.OpBitOr, g.intValue)},
		{freq: 1, generate: binaryOpGenerator(ir.OpBitXor, g.intValue)},
		{freq: 1, generate: withCast(binaryOpGenerator(ir.OpExp, g.intValue), ir.IntType)},
		{freq: 1, generate: withCast(binaryOpGenerator(ir.OpDiv, g.intValue), ir.IntType)},
		{freq: 1, generate: withCast(binaryOpGenerator(ir.OpMod, g.intValue), ir.IntType)},
		{freq: 2, generate: g.intNegation},
		{freq: 2, generate: g.intCast},
		{freq: 7, generate: g.intCall},
		{freq: 4, generate: g.intLit},
		{freq: 6, generate: g.intVar, fallback: g.intLit},
	})

	g.floatChoices = makeChoicesList(g.floatLit, []exprChoice{
		{freq: 1, generate: g.floatTernary},
		{freq: 2, generate: binaryOpGenerator(ir.OpAdd, g.floatValue)},
		{freq: 2, generate: binaryOpGenerator(ir.OpSub, g.floatValue)},
		{freq: 1, generate: binaryOpGenerator(ir.OpDiv, g.floatValue)},
		{freq: 1, generate: binaryOpGenerator(ir.OpMul, g.floatValue)},
		{freq: 5, generate: g.floatCall},
		{freq: 6, generate: g.floatVar, fallback: g.floatLit},
		{freq: 5, generate: g.floatLit},
	})

	g.stringChoices = makeChoicesList(g.stringLit, []exprChoice{
		{freq: 2, generate: g.stringCast},
		{freq: 5, generate: g.stringCall},
		{freq: 4, generate: binaryOpGenerator(ir.OpConcat, g.stringValue)},
		{freq: 5, generate: g.stringLit},
		{freq: 5, generate: g.interpolatedString},
		{freq: 6, generate: g.stringVar, fallback: g.stringLit},
		{freq: 2, generate: g.stringIndex, fallback: g.interpolatedString},
	})

	return g
}

func (g *exprGenerator) PickType() ir.Type {
	return g.pickType(0)
}

func (g *exprGenerator) pickType(depth int) ir.Type {
	if depth >= 5 {
		return g.PickScalarType()
	}

	switch g.rand.Intn(6 + depth*3) {
	case 0:
		elemType := g.pickType(depth + 1)
		return &ir.ArrayType{Elem: elemType}

	case 1:
		valueType := g.PickScalarType().(*ir.ScalarType)
		enumType := &ir.EnumType{ValueType: valueType}
		switch valueType.Kind {
		case ir.ScalarInt:
			enumType.Values = append(enumType.Values, int64(1), int64(2), int64(3))
		case ir.ScalarFloat:
			enumType.Values = append(enumType.Values, 0.424, -24.3, 32.5)
		case ir.ScalarString:
			enumType.Values = append(enumType.Values, "a", "b", "c")
		default:
			// Retry.
			return g.pickType(depth + 1)
		}
		return enumType

	default:
		return g.PickScalarType()
	}
}

func (g *exprGenerator) PickScalarType() ir.Type {
	return scalarTypes[g.rand.Intn(len(scalarTypes))]
}

func (g *exprGenerator) GenerateValueOfType(typ ir.Type) *ir.Node {
	switch typ := typ.(type) {
	case *ir.ScalarType:
		switch typ.Kind {
		case ir.ScalarBool:
			return g.boolValue()
		case ir.ScalarInt:
			return g.intValue()
		case ir.ScalarFloat:
			return g.floatValue()
		case ir.ScalarString:
			return g.stringValue()
		case ir.ScalarMixed:
			return g.mixedValue(true)
		default:
			panic(fmt.Sprintf("unexpected %s scalar type", typ.Kind))
		}

	case *ir.EnumType:
		roll := g.rand.Float64()
		if roll < 0.6 {
			if v := g.varOfType(typ); v != nil {
				return v
			}
		}
		switch typ.ValueType.Kind {
		case ir.ScalarInt:
			return ir.NewIntLit(randutil.Elem(g.rand, typ.Values).(int64))
		case ir.ScalarFloat:
			return ir.NewFloatLit(randutil.Elem(g.rand, typ.Values).(float64))
		case ir.ScalarString:
			return ir.NewStringLit(randutil.Elem(g.rand, typ.Values).(string))
		default:
			panic(fmt.Sprintf("unexpected %s enum type", typ.ValueType))
		}

	case *ir.ArrayType:
		return g.arrayValue(typ.Elem)

	default:
		panic(fmt.Sprintf("unexpected %T type", typ))
	}
}

func (g *exprGenerator) chooseExpr(list *exprChoiceList) *ir.Node {
	if g.exprDepth > 10 {
		return list.fallback()
	}
	g.exprDepth++
	defer func() { g.exprDepth-- }()

	for {
		probe := g.rand.Intn(len(list.indexMap))
		option := list.options[list.indexMap[probe]]
		n := option.generate()
		if n == nil && option.fallback != nil {
			n = option.fallback()
		}
		if n != nil {
			addParens := g.rand.Intn(10) <= 3
			if addParens {
				n = ir.NewParens(n)
			}
			return n
		}
	}
}

func (g *exprGenerator) condValue() *ir.Node {
	return g.chooseExpr(&g.condChoices)
}

func (g *exprGenerator) boolValue() *ir.Node {
	return g.chooseExpr(&g.boolChoices)
}

func (g *exprGenerator) intValue() *ir.Node {
	return g.chooseExpr(&g.intChoices)
}

func (g *exprGenerator) floatValue() *ir.Node {
	return g.chooseExpr(&g.floatChoices)
}

func (g *exprGenerator) stringValue() *ir.Node {
	return g.chooseExpr(&g.stringChoices)
}

func (g *exprGenerator) mixedValue(permitArray bool) *ir.Node {
	maxRoll := 4
	if g.exprDepth >= 10 || !permitArray {
		maxRoll = 3
	}
	switch randutil.IntRange(g.rand, 0, maxRoll) {
	case 0:
		return g.boolValue()
	case 1:
		return g.intValue()
	case 2:
		return g.floatValue()
	case 3:
		return g.stringValue()
	case 4:
		return g.arrayValue(g.PickScalarType())
	}
	panic("unreachable")
}

func (g *exprGenerator) newTernary(cond, x, y *ir.Node) *ir.Node {
	if randutil.Bool(g.rand) {
		x = ir.NewParens(x)
	}
	if randutil.Bool(g.rand) {
		y = ir.NewParens(y)
	}
	ternary := ir.NewTernary(cond, x, y)
	return ir.NewParens(ternary)
}

func (g *exprGenerator) intTernary() *ir.Node {
	return g.newTernary(g.condValue(), g.intValue(), g.intValue())
}

func (g *exprGenerator) floatTernary() *ir.Node {
	return g.newTernary(g.condValue(), g.floatValue(), g.floatValue())
}

func (g *exprGenerator) boolLit() *ir.Node {
	return boolLitValues[g.rand.Intn(len(boolLitValues))]
}

func (g *exprGenerator) intLit() *ir.Node {
	switch g.rand.Intn(8) {
	case 0, 1:
		return ir.NewIntLit(int64(g.rand.Intn(0xffff)))
	case 2, 3:
		return ir.NewIntLit(-int64(g.rand.Intn(0xffff)))
	case 4:
		return ir.NewIntLit(int64(randutil.IntRange(g.rand, 100000, 19438420511)))
	default:
		return intLitValues[g.rand.Intn(len(intLitValues))]
	}
}

func (g *exprGenerator) floatLit() *ir.Node {
	switch g.rand.Intn(8) {
	case 0:
		return ir.NewFloatLit(g.rand.Float64())
	case 2, 3:
		return ir.NewFloatLit(g.rand.Float64() * float64(g.rand.Intn(1000)))
	case 4:
		return ir.NewFloatLit(g.rand.Float64() * float64(g.rand.Intn(10000000)))
	default:
		return floatLitValues[g.rand.Intn(len(floatLitValues))]
	}
}

func (g *exprGenerator) interpolatedString() *ir.Node {
	numParts := randutil.IntRange(g.rand, 3, 8)
	n := &ir.Node{
		Op:   ir.OpInterpolatedString,
		Args: make([]*ir.Node, 0, numParts),
	}
	for i := 0; i < numParts; i++ {
		var part *ir.Node
		if randutil.Bool(g.rand) {
			v := g.varOfType(g.PickScalarType())
			if v != nil {
				part = v
			} else {
				part = g.stringLit()
			}
		} else {
			part = g.stringLit()
		}
		n.Args = append(n.Args, part)
	}
	return n
}

func (g *exprGenerator) stringLit() *ir.Node {
	var s strings.Builder
	count := randutil.IntRange(g.rand, 1, 6)
	for i := 0; i < count; i++ {
		ch := g.rand.Intn(unicode.MaxASCII)
		if !unicode.IsPrint(rune(ch)) || ch == '$' {
			s.WriteString(stringLitValues[g.rand.Intn(len(stringLitValues))].Value.(string))
		} else {
			s.WriteByte(byte(ch))
		}
	}

	return ir.NewStringLit(s.String())
}

func (g *exprGenerator) varOfType(typ ir.Type) *ir.Node {
	v := g.scope.FindVarOfType(typ)
	if v == nil {
		return nil
	}
	return ir.NewVar(v.name, v.typ)
}

func (g *exprGenerator) boolVar() *ir.Node   { return g.varOfType(ir.BoolType) }
func (g *exprGenerator) intVar() *ir.Node    { return g.varOfType(ir.IntType) }
func (g *exprGenerator) floatVar() *ir.Node  { return g.varOfType(ir.FloatType) }
func (g *exprGenerator) stringVar() *ir.Node { return g.varOfType(ir.StringType) }

func (g *exprGenerator) callOfType(fn *ir.FuncType) *ir.Node {
	g.exprDepth++
	defer func() { g.exprDepth-- }()

	numArgs := randutil.IntRange(g.rand, fn.MinArgsNum, len(fn.Params))
	callArgs := make([]*ir.Node, numArgs)
	for i := range callArgs {
		arg := g.GenerateValueOfType(fn.Params[i].Type)
		if fn.Params[i].Strict {
			arg = &ir.Node{Op: ir.OpCast, Args: []*ir.Node{g.maybeAddParens(arg)}, Type: fn.Params[i].Type}
		}
		callArgs[i] = arg
	}
	funcExpr := ir.NewName(fn.Name)
	result := ir.NewCall(funcExpr, callArgs...)
	if fn.NeedCast {
		result = &ir.Node{Op: ir.OpCast, Args: []*ir.Node{result}, Type: fn.Result}
	}
	return result
}

func (g *exprGenerator) boolCall() *ir.Node {
	return g.callOfType(g.symtab.boolFuncs[g.rand.Intn(len(g.symtab.boolFuncs))])
}

func (g *exprGenerator) intCall() *ir.Node {
	return g.callOfType(g.symtab.intFuncs[g.rand.Intn(len(g.symtab.intFuncs))])
}

func (g *exprGenerator) floatCall() *ir.Node {
	return g.callOfType(g.symtab.floatFuncs[g.rand.Intn(len(g.symtab.floatFuncs))])
}

func (g *exprGenerator) stringCall() *ir.Node {
	return g.callOfType(g.symtab.stringFuncs[g.rand.Intn(len(g.symtab.stringFuncs))])
}

func (g *exprGenerator) maybeAddParens(n *ir.Node) *ir.Node {
	if isSimpleNode(n) {
		return n
	}
	return ir.NewParens(n)
}

func (g *exprGenerator) intNegation() *ir.Node {
	return ir.NewNegation(g.maybeAddParens(g.intValue()))
}

func (g *exprGenerator) castToType(typ ir.Type) *ir.Node {
	arg := g.maybeAddParens(g.mixedValue(false))
	return &ir.Node{Op: ir.OpCast, Args: []*ir.Node{arg}, Type: typ}
}

func (g *exprGenerator) intCast() *ir.Node    { return g.castToType(ir.IntType) }
func (g *exprGenerator) stringCast() *ir.Node { return g.castToType(ir.StringType) }

func (g *exprGenerator) arrayValue(elemType ir.Type) *ir.Node {
	g.exprDepth++
	defer func() { g.exprDepth-- }()

	maxNumElems := 4
	if g.exprDepth >= 10 {
		maxNumElems = 2
	}
	numElems := randutil.IntRange(g.rand, 1, maxNumElems)
	elems := make([]*ir.Node, numElems)
	for i := 0; i < numElems; i++ {
		elems[i] = g.GenerateValueOfType(elemType)
	}
	return &ir.Node{Op: ir.OpArrayLit, Args: elems}
}

func (g *exprGenerator) lvalueOfType(typ ir.Type) *ir.Node {
	if v := g.varOfType(typ); v != nil {
		return v
	}
	return nil
}

func (g *exprGenerator) stringIndex() *ir.Node {
	lvalue := g.lvalueOfType(ir.StringType)
	if lvalue == nil {
		return nil
	}
	s := g.maybeAddParens(lvalue)
	var key *ir.Node
	if randutil.IntRange(g.rand, 0, 10) > 2 {
		key = g.intValue()
	} else {
		key = ir.NewIntLit(-1)
	}
	return ir.NewIndex(s, key)
}

var boolLitValues = []*ir.Node{
	ir.NewBoolLit(false),
	ir.NewBoolLit(true),
}

var intLitValues = []*ir.Node{
	ir.NewIntLit(0),
	ir.NewIntLit(-1),
	ir.NewIntLit(0xff),
	ir.NewIntLit(9284128),
	ir.NewIntLit(128412288),
	ir.NewIntLit(-9284120),
	ir.NewIntLit(-0xff),
}

var floatLitValues = []*ir.Node{
	ir.NewFloatLit(0),
	ir.NewFloatLit(-1),
	ir.NewFloatLit(2.51),
	ir.NewFloatLit(329.5),
	ir.NewFloatLit(0.00043),
	ir.NewFloatLit(21948.293242),
	ir.NewFloatLit(-2222.9999),
	ir.NewFloatLit(2842.6378),
}

var stringLitValues = []*ir.Node{
	ir.NewStringLit(""),
	ir.NewStringLit(","),
	ir.NewStringLit(" "),
	ir.NewStringLit("0x1f"),
	ir.NewStringLit("000"),
	ir.NewStringLit("24"),
	ir.NewStringLit("-123"),
	ir.NewStringLit("\x00"),
	ir.NewStringLit("simple string"),
	ir.NewStringLit("1\n2"),
	ir.NewStringLit("<div/>"),
	ir.NewStringLit("<h1>ok</h1>"),
	ir.NewStringLit("<p>"),
	ir.NewStringLit("</p>"),
	ir.NewStringLit(`{"key":1}`),
	ir.NewStringLit(`["val"]`),
}

var scalarTypes = []ir.Type{
	ir.BoolType,
	ir.IntType,
	ir.FloatType,
	ir.StringType,
}
