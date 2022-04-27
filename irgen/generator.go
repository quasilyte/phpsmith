package irgen

import (
	"fmt"

	"github.com/quasilyte/phpsmith/ir"
)

type generator struct {
	config *Config

	files []*File

	varNameSeq int

	currentBlock *ir.Node

	scope *scope

	symtab *symbolTable
	expr   *exprGenerator
}

func newGenerator(config *Config) *generator {
	symtab := newSymbolTable()
	s := newScope()
	return &generator{
		config: config,
		symtab: symtab,
		scope:  s,
		expr:   newExprGenerator(config, s, symtab),
	}
}

func (g *generator) CreateProgram() *Program {
	g.files = append(g.files, g.createFile("main.php"))
	return &Program{Files: g.files}
}

func (g *generator) createFile(name string) *File {
	f := &File{Name: name}
	f.Nodes = append(f.Nodes, g.createFunc("main"))
	return f
}

func (g *generator) createFunc(name string) *ir.RootFuncDecl {
	fn := &ir.RootFuncDecl{
		Name: name,
		Body: ir.NewBlock(),
	}

	g.scope.Enter()
	defer g.scope.Leave()

	g.varNameSeq = 0
	g.currentBlock = fn.Body
	for i := 0; i < 8; i++ {
		g.pushVarDecl(g.genVarname())
	}

	return fn
}

func (g *generator) genVarname() string {
	varname := fmt.Sprintf("v%d", g.varNameSeq)
	g.varNameSeq++
	return varname
}

func (g *generator) pushVarDecl(name string) {
	typ := g.expr.PickScalarType()
	lhs := ir.NewVar(name, typ)
	rhs := g.expr.GenerateValueOfType(typ)
	assign := ir.NewAssign(lhs, rhs)
	g.currentBlock.Args = append(g.currentBlock.Args, assign)
	g.scope.PushVar(name, typ)
}
