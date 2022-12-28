package irgen

import "github.com/quasilyte/phpsmith/ir"

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
}

func (symtab *symbolTable) PickRandomClass() *ir.ClassType {
	for _, c := range symtab.classes {
		return c
	}
	return nil
}

func (symtab *symbolTable) AddClass(c *ir.ClassType) {
	if symtab.sorted {
		panic("adding to a sorted symtab")
	}
	symtab.classes[c.Name] = c

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
