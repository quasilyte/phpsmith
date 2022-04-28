package irgen

import (
	"math/rand"
	"strconv"

	"github.com/quasilyte/phpsmith/ir"
	"github.com/quasilyte/phpsmith/phpfunc"
	"github.com/quasilyte/phpsmith/randutil"
)

type generator struct {
	config *Config

	rand *rand.Rand

	files []*File

	varNameSeq int

	currentBlock *ir.Node

	scope *scope

	symtab *symbolTable
	expr   *exprGenerator
}

func newGenerator(config *Config) *generator {
	symtab := newSymbolTable()
	phpfunc.Add(symtab.funcs)
	for _, fn := range symtab.funcs {
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

	s := newScope()
	return &generator{
		config: config,
		rand:   config.Rand,
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
	for i := 0; i < 2; i++ {
		f.Nodes = append(f.Nodes, g.createFunc("func"+strconv.Itoa(i)))
	}
	return f
}

func (g *generator) createFunc(name string) *ir.RootFuncDecl {
	fn := &ir.RootFuncDecl{
		Body: ir.NewBlock(),
		Type: &ir.FuncType{
			Name:   name,
			Result: ir.VoidType,
		},
	}

	g.scope.Enter()
	defer g.scope.Leave()

	g.varNameSeq = 0
	g.currentBlock = fn.Body

	blockVars := make([]string, randutil.IntRange(g.rand, 3, 7))
	for i := range blockVars {
		blockVars[i] = g.genVarname()
		g.pushVarDecl(blockVars[i])
	}
	numStatements := randutil.IntRange(g.rand, 3, 10)
	for i := 0; i < numStatements; i++ {
		g.pushStatement()
	}
	for _, name := range blockVars {
		v := g.scope.FindVarByName(name)
		varDump := ir.NewCall(ir.NewName("var_dump"), ir.NewVar(name, v.typ))
		g.currentBlock.Args = append(g.currentBlock.Args, varDump)
	}

	return fn
}

func (g *generator) genVarname() string {
	varname := "v" + strconv.Itoa(g.varNameSeq)
	g.varNameSeq++
	return varname
}

func (g *generator) pushVarDecl(name string) {
	typ := g.expr.PickType()
	lhs := ir.NewVar(name, typ)
	rhs := g.expr.GenerateValueOfType(typ)
	assign := ir.NewAssign(lhs, rhs)
	g.currentBlock.Args = append(g.currentBlock.Args, assign)
	g.scope.PushVar(name, typ)
}

func (g *generator) pushStatement() {
	switch randutil.IntRange(g.rand, 0, 4) {
	case 0, 1, 2:
		g.pushVarDecl(g.genVarname())
	case 3:
		g.pushBlockStmt()
	case 4:
		g.pushIfStmt()
	}
}

func (g *generator) pushBlockStmt() {
	newBlock := &ir.Node{Op: ir.OpBlock}
	oldBlock := g.currentBlock
	g.currentBlock = newBlock
	numStatements := randutil.IntRange(g.rand, 1, 3)
	for i := 0; i < numStatements; i++ {
		g.pushStatement()
	}
	oldBlock.Args = append(oldBlock.Args, newBlock)
	g.currentBlock = oldBlock
}

func (g *generator) pushIfStmt() {
	cond := g.expr.condValue()

	withoutBlock := randutil.IntRange(g.rand, 0, 1) == 0
	oldBlock := g.currentBlock
	if withoutBlock {
		ifStmt := &ir.Node{Op: ir.OpIf, Args: []*ir.Node{cond}}
		g.currentBlock = ifStmt
		g.pushStatement()
		oldBlock.Args = append(oldBlock.Args, ir.NewIf(cond, ifStmt))
	} else {
		newBlock := &ir.Node{Op: ir.OpBlock}
		g.currentBlock = newBlock
		g.pushStatement()
		oldBlock.Args = append(oldBlock.Args, ir.NewIf(cond, newBlock))
	}
	g.currentBlock = oldBlock
}
