package irgen

import (
	"math/rand"
	"strconv"

	"github.com/quasilyte/phpsmith/ir"
	"github.com/quasilyte/phpsmith/phpdoc"
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
	g.files = append(g.files, g.createMainFile())
	return &Program{Files: g.files}
}

func (g *generator) createMainFile() *File {
	f := &File{Name: "main.php"}

	funcs := make([]*ir.RootFuncDecl, randutil.IntRange(g.rand, 2, 4))
	for i := range funcs {
		funcs[i] = g.createFunc("func" + strconv.Itoa(i))
	}

	// Create a main func.
	mainFunc := &ir.RootFuncDecl{
		Type: &ir.FuncType{
			Name:   "main",
			Result: ir.VoidType,
		},
		Body: &ir.Node{Op: ir.OpBlock},
	}
	for _, fn := range funcs {
		funcNode := ir.NewName(fn.Type.Name)
		call := &ir.Node{Op: ir.OpCall, Args: []*ir.Node{funcNode}}
		mainFunc.Body.Args = append(mainFunc.Body.Args, call)
	}

	for _, fn := range funcs {
		f.Nodes = append(f.Nodes, fn)
	}
	f.Nodes = append(f.Nodes, mainFunc)

	f.Nodes = append(f.Nodes, &ir.RootStmt{
		X: ir.NewCall(ir.NewName("main")),
	})

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
	defer func() {
		g.scope.Leave()
		if len(g.scope.depths) != 0 {
			panic("corrupted scope stack?")
		}
	}()

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
	if scalarType, ok := typ.(*ir.ScalarType); ok && scalarType.Kind == ir.ScalarBool {
		assign.Value = &phpdoc.VarTag{VarName: "$" + name, Type: "bool"}
	}
	g.currentBlock.Args = append(g.currentBlock.Args, assign)
	g.scope.PushVar(name, typ)
}

func (g *generator) pushStatement() {
	switch randutil.IntRange(g.rand, 0, 7) {
	case 0, 1, 2:
		g.pushVarDecl(g.genVarname())
	case 3:
		g.pushBlockStmt()
	case 4:
		g.pushIfStmt()
	case 5:
		g.pushVarDump()
	case 6, 7:
		g.pushAssignStmt()
	}
}

func (g *generator) pushAssignStmt() {
	v := g.pickVar()
	if v == nil {
		g.pushVarDecl(g.genVarname())
		return
	}
	var op ir.Op
	if typ, ok := v.typ.(*ir.ScalarType); ok && randutil.Bool(g.rand) {
		var opChoice []ir.Op
		switch typ.Kind {
		case ir.ScalarInt, ir.ScalarFloat:
			opChoice = []ir.Op{ir.OpAdd, ir.OpSub}
		case ir.ScalarString:
			opChoice = []ir.Op{ir.OpConcat}
		case ir.ScalarMixed:
			opChoice = []ir.Op{ir.OpAdd, ir.OpSub, ir.OpConcat}
		}
		if len(opChoice) != 0 {
			op = opChoice[g.rand.Intn(len(opChoice))]
		}
	}
	var assign *ir.Node
	lhs := ir.NewVar(v.name, v.typ)
	rhs := g.expr.GenerateValueOfType(v.typ)
	if op != ir.OpInvalid {
		assign = ir.NewAssignModify(op, lhs, rhs)
	} else {
		assign = ir.NewAssign(lhs, rhs)
	}
	g.currentBlock.Args = append(g.currentBlock.Args, assign)
}

func (g *generator) pushVarDump() {
	typ := g.expr.PickType()
	arg := g.expr.GenerateValueOfType(typ)
	switch typ.(type) {
	case *ir.ScalarType, *ir.ArrayType:
		varDump := ir.NewCall(ir.NewName("var_dump"), arg)
		g.currentBlock.Args = append(g.currentBlock.Args, varDump)
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
	g.scope.Enter()
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
	g.scope.Leave()
	g.currentBlock = oldBlock
}

func (g *generator) pickVar() *scopeVar {
	blockVars := g.scope.CurrentBlockVars()
	if len(blockVars) == 0 {
		return nil
	}
	return &blockVars[g.rand.Intn(len(blockVars))]
}
