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
