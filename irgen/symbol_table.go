package irgen

import (
	"fmt"
	"sort"

	"github.com/quasilyte/phpsmith/ir"
)

type symbolTable struct {
	funcs   map[string]*ir.FuncType
	classes map[string]*ir.ClassType

	sorted bool

	boolFields   []fieldRef
	intFields    []fieldRef
	floatFields  []fieldRef
	stringFields []fieldRef
	arrayFields  []fieldRef

	voidFuncs   []*ir.FuncType
	boolFuncs   []*ir.FuncType
	intFuncs    []*ir.FuncType
	floatFuncs  []*ir.FuncType
	stringFuncs []*ir.FuncType
	arrayFuncs  []*ir.FuncType
}

type fieldRef struct {
	index int
	class *ir.ClassType
}

func (ref fieldRef) Get() *ir.TypeField {
	return &ref.class.Fields[ref.index]
}

func newSymbolTable() *symbolTable {
	return &symbolTable{
		funcs:   make(map[string]*ir.FuncType),
		classes: make(map[string]*ir.ClassType),
	}
}

func (symtab *symbolTable) Sort() {
	if symtab.sorted {
		panic("double symtab sorting")
	}
	symtab.sorted = true

	sort.SliceStable(symtab.arrayFuncs, func(i, j int) bool {
		return typeLess(symtab.arrayFuncs[i].Result, symtab.arrayFuncs[j].Result)
	})
	sort.SliceStable(symtab.arrayFields, func(i, j int) bool {
		return typeLess(symtab.arrayFields[i].Get().Type, symtab.arrayFields[j].Get().Type)
	})
}

func (symtab *symbolTable) FindFuncsOfType(typ ir.Type) []*ir.FuncType {
	switch typ := typ.(type) {
	case *ir.ScalarType:
		switch typ.Kind {
		case ir.ScalarBool:
			return symtab.boolFuncs
		case ir.ScalarInt:
			return symtab.intFuncs
		case ir.ScalarFloat:
			return symtab.floatFuncs
		case ir.ScalarString:
			return symtab.stringFuncs
		case ir.ScalarVoid:
			return symtab.voidFuncs
		default:
			return nil
		}

	case *ir.ArrayType:
		i := sort.Search(len(symtab.arrayFuncs), func(i int) bool {
			return !typeLess(symtab.arrayFuncs[i].Result, typ)
		})
		if i < len(symtab.arrayFuncs) && typesIdentical(symtab.arrayFuncs[i].Result, typ) {
			j := i
			for j < len(symtab.arrayFuncs)-1 && typesIdentical(symtab.arrayFuncs[j+1].Result, typ) {
				j++
			}
			return symtab.arrayFuncs[i : j+1]
		}
		return nil

	default:
		return nil
	}

}

func (symtab *symbolTable) PickRandomClass() *ir.ClassType {
	for _, c := range symtab.classes {
		return c
	}
	return nil
}

func (symtab *symbolTable) DeclareClass(name string) {
	if symtab.classes[name] != nil {
		panic(fmt.Sprintf("class %s is already declared", name))
	}
	symtab.classes[name] = &ir.ClassType{Name: name}
}

func (symtab *symbolTable) DefineClass(c *ir.ClassType) {
	if symtab.sorted {
		panic("adding to a sorted symtab")
	}
	declared := symtab.classes[c.Name]
	if declared == nil {
		panic(fmt.Sprintf("class %s was not declared and can't be defined", c.Name))
	}
	*declared = *c

	for i, field := range c.Fields {
		switch fieldType := field.Type.(type) {
		case *ir.ScalarType:
			switch fieldType.Kind {
			case ir.ScalarBool:
				symtab.boolFields = append(symtab.boolFields, fieldRef{index: i, class: c})
			case ir.ScalarInt:
				symtab.intFields = append(symtab.intFields, fieldRef{index: i, class: c})
			case ir.ScalarFloat:
				symtab.floatFields = append(symtab.floatFields, fieldRef{index: i, class: c})
			case ir.ScalarString:
				symtab.stringFields = append(symtab.stringFields, fieldRef{index: i, class: c})
			}

		case *ir.ArrayType:
			symtab.arrayFields = append(symtab.arrayFields, fieldRef{index: i, class: c})
		}
	}
}

func (symtab *symbolTable) AddFunc(fn *ir.FuncType) {
	if symtab.sorted {
		panic("adding to a sorted symtab")
	}
	symtab.funcs[fn.Name] = fn

	switch resultType := fn.Result.(type) {
	case *ir.ScalarType:
		switch resultType.Kind {
		case ir.ScalarVoid:
			symtab.voidFuncs = append(symtab.voidFuncs, fn)
		case ir.ScalarBool:
			symtab.boolFuncs = append(symtab.boolFuncs, fn)
		case ir.ScalarInt:
			symtab.intFuncs = append(symtab.intFuncs, fn)
		case ir.ScalarFloat:
			symtab.floatFuncs = append(symtab.floatFuncs, fn)
		case ir.ScalarString:
			symtab.stringFuncs = append(symtab.stringFuncs, fn)
		}

	case *ir.ArrayType:
		symtab.arrayFuncs = append(symtab.arrayFuncs, fn)
	}
}
