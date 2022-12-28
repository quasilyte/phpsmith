package irgen

import (
	"fmt"
	"math/rand"

	"github.com/quasilyte/phpsmith/ir"
	"github.com/quasilyte/phpsmith/randutil"
)

type exprGenerator struct {
	config *Config

	rand *rand.Rand

	valueGenerator *valueGenerator

	scope *scope

	symtab  *symbolTable
	scratch []*scopeVar

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
		config:         config,
		scope:          s,
		symtab:         symtab,
		rand:           config.Rand,
		valueGenerator: newValueGenerator(config.Rand),
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

	binaryOpGenerator := func(op ir.Op, typeHint ir.Type, operandGenerator func() *ir.Node) func() *ir.Node {
		return func() *ir.Node {
			x := operandGenerator()
			y := operandGenerator()
			return &ir.Node{Op: op, Args: []*ir.Node{g.maybeAddParens(x), g.maybeAddParens(y)}, Type: typeHint}
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
		{freq: 2, generate: g.boolFieldAccess, fallback: g.boolLit},
		{freq: 3, generate: cmpOpGenerator(ir.OpEqual2)},
		{freq: 3, generate: cmpOpGenerator(ir.OpEqual3)},
		{freq: 4, generate: binaryOpGenerator(ir.OpAnd, nil, g.boolValue)},
		{freq: 4, generate: binaryOpGenerator(ir.OpOr, nil, g.boolValue)},
		{freq: 4, generate: unaryOpGenerator(ir.OpNot, g.condValue)},
		{freq: 6, generate: g.boolCall},
		{freq: 1, generate: g.boolLit},
	})

	g.boolChoices = makeChoicesList(g.boolLit, []exprChoice{
		{freq: 1, generate: cmpOpGenerator(ir.OpEqual2)},
		{freq: 1, generate: cmpOpGenerator(ir.OpEqual3)},
		{freq: 2, generate: g.boolFieldAccess, fallback: g.boolLit},
		{freq: 3, generate: binaryOpGenerator(ir.OpAnd, nil, g.boolValue)},
		{freq: 3, generate: binaryOpGenerator(ir.OpOr, nil, g.boolValue)},
		{freq: 3, generate: g.boolLit},
		{freq: 4, generate: unaryOpGenerator(ir.OpNot, g.condValue)},
		{freq: 4, generate: g.boolCall},
	})

	g.intChoices = makeChoicesList(g.intLit, []exprChoice{
		{freq: 1, generate: g.intTernary},
		{freq: 2, generate: withCast(binaryOpGenerator(ir.OpAdd, ir.IntType, g.intValue), ir.IntType)},
		{freq: 2, generate: binaryOpGenerator(ir.OpSub, ir.IntType, g.intValue)},
		{freq: 1, generate: withCast(binaryOpGenerator(ir.OpMul, ir.IntType, g.intValue), ir.IntType)},
		{freq: 1, generate: binaryOpGenerator(ir.OpBitAnd, ir.IntType, g.intValue)},
		{freq: 1, generate: binaryOpGenerator(ir.OpBitOr, ir.IntType, g.intValue)},
		{freq: 1, generate: binaryOpGenerator(ir.OpBitXor, ir.IntType, g.intValue)},
		{freq: 1, generate: withCast(binaryOpGenerator(ir.OpExp, ir.IntType, g.intValue), ir.IntType)},
		{freq: 1, generate: withCast(binaryOpGenerator(ir.OpDiv, ir.IntType, g.intValue), ir.IntType)},
		{freq: 1, generate: withCast(binaryOpGenerator(ir.OpMod, ir.IntType, g.intValue), ir.IntType)},
		{freq: 2, generate: g.intFieldAccess, fallback: g.intLit},
		{freq: 2, generate: g.intNegation},
		{freq: 2, generate: g.intCast},
		{freq: 7, generate: g.intCall},
		{freq: 4, generate: g.intLit},
	})

	g.floatChoices = makeChoicesList(g.floatLit, []exprChoice{
		{freq: 1, generate: g.floatTernary},
		{freq: 2, generate: binaryOpGenerator(ir.OpAdd, ir.FloatType, g.floatValue)},
		{freq: 2, generate: binaryOpGenerator(ir.OpSub, ir.FloatType, g.floatValue)},
		{freq: 2, generate: g.floatFieldAccess, fallback: g.floatLit},
		{freq: 1, generate: binaryOpGenerator(ir.OpDiv, ir.FloatType, g.floatValue)},
		{freq: 1, generate: binaryOpGenerator(ir.OpMul, ir.FloatType, g.floatValue)},
		{freq: 5, generate: g.floatCall},
		{freq: 5, generate: g.floatLit},
	})

	g.stringChoices = makeChoicesList(g.stringLit, []exprChoice{
		{freq: 2, generate: g.stringCast},
		{freq: 2, generate: g.stringFieldAccess, fallback: g.stringLit},
		{freq: 5, generate: g.stringCall},
		{freq: 4, generate: binaryOpGenerator(ir.OpConcat, ir.StringType, g.stringValue)},
		{freq: 5, generate: g.stringLit},
		{freq: 5, generate: g.interpolatedString},
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

	switch g.rand.Intn(7 + depth*3) {
	case 0:
		elemType := g.pickType(depth + 1)
		return &ir.ArrayType{Elem: elemType}

	case 1:
		return g.pickTupleType(depth + 2)

	case 2:
		return g.PickEnumType()

	case 3, 4:
		return g.maybePickClassType(depth)

	default:
		return g.PickScalarType()
	}
}

func (g *exprGenerator) maybePickClassType(depth int) ir.Type {
	c := g.symtab.PickRandomClass()
	if c != nil {
		return c
	}
	return g.pickType(depth + 1)
}

func (g *exprGenerator) pickTupleType(depth int) ir.Type {
	numElems := randutil.IntRange(g.rand, 1, 12)
	tuple := &ir.TupleType{
		Elems: make([]ir.Type, 0, numElems),
	}
	for i := 0; i < numElems; i++ {
		tuple.Elems = append(tuple.Elems, g.pickType(depth))
	}
	return tuple
}

func (g *exprGenerator) PickEnumType() ir.Type {
	valueType := g.PickScalarTypeNoBool().(*ir.ScalarType)
	enumType := &ir.EnumType{ValueType: valueType}
	switch valueType.Kind {
	case ir.ScalarInt:
		enumType.Values = toEfaceSlice(generateUniqueValues(randutil.IntRange(g.rand, 2, 30), func() int64 {
			return g.valueGenerator.IntValue()
		}))
	case ir.ScalarFloat:
		enumType.Values = toEfaceSlice(generateUniqueValues(randutil.IntRange(g.rand, 2, 16), func() float64 {
			return g.valueGenerator.FloatValue()
		}))
	case ir.ScalarString:
		enumType.Values = toEfaceSlice(generateUniqueValues(randutil.IntRange(g.rand, 2, 20), func() string {
			return g.valueGenerator.StringValue()
		}))
	default:
		panic("unreachable")
	}
	return enumType
}

func (g *exprGenerator) PickScalarTypeNoBool() ir.Type {
	return scalarTypesNoBool[g.rand.Intn(len(scalarTypesNoBool))]
}

func (g *exprGenerator) PickScalarType() ir.Type {
	return scalarTypes[g.rand.Intn(len(scalarTypes))]
}

func (g *exprGenerator) GenerateConstValueOfType(typ ir.Type) *ir.Node {
	switch typ := typ.(type) {
	case *ir.ScalarType:
		switch typ.Kind {
		case ir.ScalarBool:
			return g.boolLit()
		case ir.ScalarInt:
			return g.intLit()
		case ir.ScalarFloat:
			return g.floatLit()
		case ir.ScalarString:
			return g.stringLit()
		case ir.ScalarMixed:
			return g.stringLit()
		default:
			panic(fmt.Sprintf("unexpected %s scalar type", typ.Kind))
		}

	case *ir.EnumType:
		return g.enumValue(typ)

	case *ir.ArrayType:
		return g.constArrayValue(typ.Elem)

	case *ir.TupleType:
		return g.constTupleValue(typ)

	case *ir.ClassType:
		return ir.NewName("null")

	default:
		panic(fmt.Sprintf("unexpected %T type", typ))
	}
}

func (g *exprGenerator) enumValue(typ *ir.EnumType) *ir.Node {
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
}

func (g *exprGenerator) GenerateValueOfType(typ ir.Type) *ir.Node {
	if randutil.Chance(g.rand, 0.4) {
		v := g.varOfType(typ)
		if v != nil {
			return v
		}
	}

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
		return g.enumValue(typ)

	case *ir.ArrayType:
		return g.arrayValue(typ.Elem)

	case *ir.TupleType:
		return g.tupleValue(typ)

	case *ir.ClassType:
		return g.classValue(typ)

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
	return ir.NewBoolLit(g.valueGenerator.BoolValue())
}

func (g *exprGenerator) intLit() *ir.Node {
	return ir.NewIntLit(g.valueGenerator.IntValue())
}

func (g *exprGenerator) floatLit() *ir.Node {
	return ir.NewFloatLit(g.valueGenerator.FloatValue())
}

func (g *exprGenerator) boolFieldAccess() *ir.Node   { return g.fieldAccessOfType(ir.BoolType) }
func (g *exprGenerator) intFieldAccess() *ir.Node    { return g.fieldAccessOfType(ir.IntType) }
func (g *exprGenerator) floatFieldAccess() *ir.Node  { return g.fieldAccessOfType(ir.FloatType) }
func (g *exprGenerator) stringFieldAccess() *ir.Node { return g.fieldAccessOfType(ir.StringType) }

func (g *exprGenerator) fieldAccessOfType(typ ir.Type) *ir.Node {
	var list []fieldRef
	switch typ := typ.(type) {
	case *ir.ScalarType:
		switch typ.Kind {
		case ir.ScalarBool:
			list = g.symtab.boolFields
		case ir.ScalarInt:
			list = g.symtab.intFields
		case ir.ScalarFloat:
			list = g.symtab.floatFields
		case ir.ScalarString:
			list = g.symtab.stringFields
		}
	}
	if len(list) == 0 {
		return nil
	}
	ref := list[g.rand.Intn(len(list))]
	instance := g.maybeAddParens(g.GenerateValueOfType(ref.class))
	field := ref.class.Fields[ref.index]
	return ir.NewMemberAccess(instance, field.Name)
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
	return ir.NewStringLit(g.valueGenerator.StringValue())
}

func (g *exprGenerator) varOfType(typ ir.Type) *ir.Node {
	v := g.scope.FindVarOfType(typ)
	if v == nil {
		return nil
	}
	return ir.NewVar(v.name, v.typ)
}

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
	var funcExpr *ir.Node
	if fn.Class == nil {
		funcExpr = ir.NewName(fn.Name)
	} else {
		instance := g.maybeAddParens(g.GenerateValueOfType(fn.Class))
		funcExpr = ir.NewMemberAccess(instance, fn.Name)
	}
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

func (g *exprGenerator) classValue(typ *ir.ClassType) *ir.Node {
	g.exprDepth++
	defer func() { g.exprDepth-- }()

	return &ir.Node{
		Op:    ir.OpNew,
		Value: typ.Name,
		Type:  typ,
	}
}

func (g *exprGenerator) constTupleValue(typ *ir.TupleType) *ir.Node {
	return g.makeTupleValue(typ, g.GenerateConstValueOfType)
}

func (g *exprGenerator) tupleValue(typ *ir.TupleType) *ir.Node {
	return g.makeTupleValue(typ, g.GenerateValueOfType)
}

func (g *exprGenerator) makeTupleValue(typ *ir.TupleType, elemGen func(ir.Type) *ir.Node) *ir.Node {
	g.exprDepth++
	defer func() { g.exprDepth-- }()

	numElems := len(typ.Elems)
	elems := make([]*ir.Node, numElems)
	for i := 0; i < numElems; i++ {
		elems[i] = elemGen(typ.Elems[i])
	}
	return ir.NewCall(ir.NewName("tuple"), elems...)
}

func (g *exprGenerator) constArrayValue(elemType ir.Type) *ir.Node {
	return g.makeArrayValue(elemType, g.GenerateConstValueOfType)
}

func (g *exprGenerator) arrayValue(elemType ir.Type) *ir.Node {
	return g.makeArrayValue(elemType, g.GenerateValueOfType)
}

func (g *exprGenerator) makeArrayValue(elemType ir.Type, elemGen func(ir.Type) *ir.Node) *ir.Node {
	g.exprDepth++
	defer func() { g.exprDepth-- }()

	maxNumElems := 4
	if g.exprDepth >= 10 {
		maxNumElems = 2
	}
	numElems := randutil.IntRange(g.rand, 1, maxNumElems)
	elems := make([]*ir.Node, numElems)
	for i := 0; i < numElems; i++ {
		elems[i] = elemGen(elemType)
	}
	return &ir.Node{Op: ir.OpArrayLit, Args: elems}
}

func (g *exprGenerator) findRandomVar(predicate func(v *scopeVar) bool) *scopeVar {
	g.scratch = g.scratch[:0]
	g.scope.FindVar(func(v *scopeVar) bool {
		if !predicate(v) {
			return false
		}
		g.scratch = append(g.scratch, v)
		return false
	})
	if len(g.scratch) == 0 {
		return nil
	}
	return randutil.Elem(g.rand, g.scratch)
}

func (g *exprGenerator) PickLvalue() (*ir.Node, ir.Type) {
	roll := g.rand.Float64()

	if roll < 0.4 {
		classVar := g.findRandomVar(func(v *scopeVar) bool {
			_, ok := v.typ.(*ir.ClassType)
			return ok
		})
		if classVar != nil {
			classType := classVar.typ.(*ir.ClassType)
			if len(classType.Fields) != 0 {
				field := randutil.Elem(g.rand, classType.Fields)
				return ir.NewMemberAccess(ir.NewVar(classVar.name, classVar.typ), field.Name), field.Type
			}
		}
	}

	blockVars := g.scope.CurrentBlockVars()
	if len(blockVars) == 0 {
		return nil, nil
	}
	v := blockVars[g.rand.Intn(len(blockVars))]
	return ir.NewVar(v.name, v.typ), v.typ
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
	return &ir.Node{
		Op:   ir.OpCast,
		Args: []*ir.Node{ir.NewIndex(s, key)},
		Type: ir.StringType,
	}
}

var scalarTypes = []ir.Type{
	ir.BoolType,
	ir.IntType,
	ir.FloatType,
	ir.StringType,
}

var scalarTypesNoBool = []ir.Type{
	ir.IntType,
	ir.FloatType,
	ir.StringType,
}
