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

	scope *scope

	symtab *symbolTable

	condChoices   exprChoiceList
	boolChoices   exprChoiceList
	intChoices    exprChoiceList
	floatChoices  exprChoiceList
	stringChoices exprChoiceList
}

type exprChoiceList struct {
	indexMap []uint16
	options  []exprChoice
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

	makeChoicesList := func(options []exprChoice) exprChoiceList {
		indexes := make([]uint16, 0, len(options)*4)
		for i, o := range options {
			for j := 0; j < o.freq; j++ {
				indexes = append(indexes, uint16(i))
			}
		}
		return exprChoiceList{
			indexMap: indexes,
			options:  options,
		}
	}

	binaryOpGenerator := func(op ir.Op, operandGenerator func() *ir.Node) func() *ir.Node {
		return func() *ir.Node {
			x := operandGenerator()
			y := operandGenerator()
			return &ir.Node{Op: op, Args: []*ir.Node{x, y}}
		}
	}

	unaryOpGenerator := func(op ir.Op, operandGenerator func() *ir.Node) func() *ir.Node {
		return func() *ir.Node {
			x := operandGenerator()
			return &ir.Node{Op: op, Args: []*ir.Node{x}}
		}
	}

	g.condChoices = makeChoicesList([]exprChoice{
		{freq: 3, generate: binaryOpGenerator(ir.OpAnd, g.boolValue)},
		{freq: 3, generate: binaryOpGenerator(ir.OpOr, g.boolValue)},
		{freq: 4, generate: unaryOpGenerator(ir.OpNot, g.condValue)},
		{freq: 5, generate: g.boolVar, fallback: g.boolLit},
		{freq: 1, generate: g.boolLit},
	})

	g.boolChoices = makeChoicesList([]exprChoice{
		{freq: 6, generate: g.boolVar, fallback: g.boolLit},
		{freq: 4, generate: g.boolLit},
	})

	g.intChoices = makeChoicesList([]exprChoice{
		{freq: 1, generate: g.intTernary},
		{freq: 2, generate: binaryOpGenerator(ir.OpAdd, g.intValue)},
		{freq: 2, generate: binaryOpGenerator(ir.OpSub, g.intValue)},
		{freq: 6, generate: g.intVar, fallback: g.intLit},
		{freq: 4, generate: g.intCall},
		{freq: 5, generate: g.intLit},
	})

	g.floatChoices = makeChoicesList([]exprChoice{
		{freq: 1, generate: g.floatTernary},
		{freq: 2, generate: binaryOpGenerator(ir.OpAdd, g.floatValue)},
		{freq: 2, generate: binaryOpGenerator(ir.OpSub, g.floatValue)},
		{freq: 6, generate: g.floatVar, fallback: g.floatLit},
		{freq: 5, generate: g.floatLit},
	})

	g.stringChoices = makeChoicesList([]exprChoice{
		{freq: 4, generate: binaryOpGenerator(ir.OpConcat, g.stringValue)},
		{freq: 6, generate: g.stringVar, fallback: g.stringLit},
		{freq: 5, generate: g.stringLit},
	})

	return g
}

func (g *exprGenerator) PickType() ir.Type {
	switch g.rand.Intn(3) {
	case 0:
		return g.PickArrayType()
	default:
		return g.PickScalarType()
	}
}

func (g *exprGenerator) PickArrayType() ir.Type {
	elemType := g.PickType()
	return &ir.ArrayType{Elem: elemType}
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
		default:
			panic(fmt.Sprintf("unexpected %s scalar type", typ.Kind))
		}

	case *ir.ArrayType:
		return g.arrayValue(typ.Elem)

	default:
		panic(fmt.Sprintf("unexpected %T type", typ))
	}
}

func (g *exprGenerator) chooseExpr(list *exprChoiceList) *ir.Node {
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

func (g *exprGenerator) intTernary() *ir.Node {
	return ir.NewTernary(g.condValue(), g.intValue(), g.intValue())
}

func (g *exprGenerator) floatTernary() *ir.Node {
	return ir.NewTernary(g.condValue(), g.floatValue(), g.floatValue())
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
		return ir.NewIntLit(g.rand.Int63())
	default:
		return intLitValues[g.rand.Intn(len(intLitValues))]
	}
}

func (g *exprGenerator) floatLit() *ir.Node {
	return floatLitValues[g.rand.Intn(len(floatLitValues))]
}

func (g *exprGenerator) stringLit() *ir.Node {
	return stringLitValues[g.rand.Intn(len(stringLitValues))]
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

func (g *exprGenerator) intCall() *ir.Node {
	fn := g.symtab.intFuncs[g.rand.Intn(len(g.symtab.intFuncs))]
	callArgs := make([]*ir.Node, len(fn.Params))
	for i := range callArgs {
		callArgs[i] = g.GenerateValueOfType(fn.Params[i].Type)
	}
	funcExpr := ir.NewName(fn.Name)
	return ir.NewCall(funcExpr, callArgs...)
}

func (g *exprGenerator) arrayValue(elemType ir.Type) *ir.Node {
	numElems := randutil.IntRange(g.rand, 1, 4)
	elems := make([]*ir.Node, numElems)
	for i := 0; i < numElems; i++ {
		elems[i] = g.GenerateValueOfType(elemType)
	}
	return &ir.Node{Op: ir.OpArrayLit, Args: elems}
}

var boolLitValues = []*ir.Node{
	ir.NewBoolLit(false),
	ir.NewBoolLit(true),
}

var intLitValues = []*ir.Node{
	ir.NewIntLit(0),
	ir.NewIntLit(-1),
	ir.NewIntLit(0xff),
	ir.NewIntLit(-0xff),
}

var floatLitValues = []*ir.Node{
	ir.NewFloatLit(0),
	ir.NewFloatLit(-1),
	ir.NewFloatLit(2.51),
	ir.NewFloatLit(2842.6378),
}

var stringLitValues = []*ir.Node{
	ir.NewStringLit(""),
	ir.NewStringLit("\x00"),
	ir.NewStringLit("simple string"),
	ir.NewStringLit("1\n2"),
}

var scalarTypes = []ir.Type{
	ir.BoolType,
	ir.IntType,
	ir.FloatType,
	ir.StringType,
}
