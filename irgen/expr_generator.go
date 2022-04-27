package irgen

import (
	"fmt"
	"math/rand"

	"github.com/quasilyte/phpsmith/ir"
)

type exprGenerator struct {
	config *Config

	rand *rand.Rand

	scope *scope

	symtab *symbolTable
}

func newExprGenerator(config *Config, s *scope, symtab *symbolTable) *exprGenerator {
	return &exprGenerator{
		config: config,
		scope:  s,
		symtab: symtab,
		rand:   config.Rand,
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
		default:
			panic(fmt.Sprintf("unexpected %s scalar type", typ.Kind))
		}

	default:
		panic(fmt.Sprintf("unexpected %T type", typ))
	}
}

func (g *exprGenerator) boolValue() *ir.Node {
	return ir.NewBoolLit(true)
}

func (g *exprGenerator) intValue() *ir.Node {
	pick := g.rand.Intn(3)
	switch pick {
	case 0:
		op := intBinOp[g.rand.Intn(len(intBinOp))]
		return &ir.Node{Op: op, Args: []*ir.Node{g.intValue(), g.intValue()}}
	case 1:
		v := g.scope.FindVarOfType(intType)
		if v == nil {
			return g.intValue()
		}
		return ir.NewVar(v.name, v.typ)
	default:
		return g.intLit()
	}
}

func (g *exprGenerator) floatValue() *ir.Node {
	return ir.NewFloatLit(2.5)
}

func (g *exprGenerator) stringValue() *ir.Node {
	pick := g.rand.Intn(3)
	switch pick {
	case 0:
		return ir.NewConcat(g.stringValue(), g.stringValue())
	case 1:
		v := g.scope.FindVarOfType(stringType)
		if v == nil {
			return g.stringValue()
		}
		return ir.NewVar(v.name, v.typ)
	default:
		return g.stringLit()
	}
}

func (g *exprGenerator) intLit() *ir.Node {
	return intLitValues[g.rand.Intn(len(intLitValues))]
}

func (g *exprGenerator) stringLit() *ir.Node {
	return stringLitValues[g.rand.Intn(len(stringLitValues))]
}

var intBinOp = []ir.Op{
	ir.OpAdd,
	ir.OpSub,
}

var intLitValues = []*ir.Node{
	ir.NewIntLit(0),
	ir.NewIntLit(-1),
	ir.NewIntLit(0xff),
	ir.NewIntLit(-0xff),
}

var stringLitValues = []*ir.Node{
	ir.NewStringLit(""),
	ir.NewStringLit("\x00"),
	ir.NewStringLit("simple string"),
	ir.NewStringLit("1\n2"),
}

var scalarTypes = []ir.Type{
	boolType,
	intType,
	floatType,
	stringType,
}
