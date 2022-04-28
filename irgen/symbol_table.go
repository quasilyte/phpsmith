package irgen

import "github.com/quasilyte/phpsmith/ir"

type symbolTable struct {
	funcs map[string]*ir.FuncType

	voidFuncs   []*ir.FuncType
	boolFuncs   []*ir.FuncType
	intFuncs    []*ir.FuncType
	floatFuncs  []*ir.FuncType
	stringFuncs []*ir.FuncType
	arrayFuncs  []*ir.FuncType
}

func newSymbolTable() *symbolTable {
	return &symbolTable{
		funcs: make(map[string]*ir.FuncType),
	}
}

func (symtab *symbolTable) AddFunc(fn *ir.FuncType) {
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
