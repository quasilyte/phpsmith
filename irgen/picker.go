package irgen

import (
	"fmt"
	"math/rand"

	"github.com/quasilyte/phpsmith/ir"
)

type picker struct {
	config *Config

	rand *rand.Rand

	scope *scope

	symtab *symbolTable
}

func newPicker(config *Config, s *scope, symtab *symbolTable) *picker {
	return &picker{
		config: config,
		scope:  s,
		symtab: symtab,
		rand:   config.Rand,
	}
}

func (p *picker) PickScalarType() ir.Type {
	return scalarTypes[p.rand.Intn(len(scalarTypes))]
}

func (p *picker) PickValueOfType(typ ir.Type) *ir.Node {
	switch typ := typ.(type) {
	case *ir.ScalarType:
		switch typ.Kind {
		case ir.ScalarBool:
			return p.boolValue()
		case ir.ScalarInt:
			return p.intValue()
		case ir.ScalarFloat:
			return p.floatValue()
		case ir.ScalarString:
			return p.stringValue()
		default:
			panic(fmt.Sprintf("unexpected %s scalar type", typ.Kind))
		}

	default:
		panic(fmt.Sprintf("unexpected %T type", typ))
	}
}

func (p *picker) boolValue() *ir.Node {
	return ir.NewBoolLit(true)
}

func (p *picker) intValue() *ir.Node {
	return ir.NewIntLit(4)
}

func (p *picker) floatValue() *ir.Node {
	return ir.NewFloatLit(2.5)
}

func (p *picker) stringValue() *ir.Node {
	pick := p.rand.Intn(3)
	switch pick {
	case 0:
		return p.stringLit()
	case 1:
		v := p.scope.FindVarOfType(stringType)
		if v == nil {
			return p.stringValue()
		}
		return ir.NewVar(v.name, v.typ)
	default:
		return ir.NewConcat(p.stringValue(), p.stringValue())
	}
}

func (p *picker) stringLit() *ir.Node {
	return ir.NewStringLit(stringLitValues[p.rand.Intn(len(stringLitValues))])
}

var stringLitValues = []string{
	"",
	"simple string",
}

var scalarTypes = []ir.Type{
	boolType,
	intType,
	floatType,
	stringType,
}
