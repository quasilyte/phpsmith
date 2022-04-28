package irgen

import "github.com/quasilyte/phpsmith/ir"

type symbolTable struct {
	funcs map[string]*ir.FuncType
}

func newSymbolTable() *symbolTable {
	return &symbolTable{
		funcs: make(map[string]*ir.FuncType),
	}
}
