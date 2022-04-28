package irgen

import (
	_ "embed"
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/quasilyte/phpsmith/ir"
	"github.com/quasilyte/phpsmith/phpdoc"
	"github.com/quasilyte/phpsmith/phpfunc"
	"github.com/quasilyte/phpsmith/randutil"
)

//go:embed _php/fuzzlib.php
var phpFuzzlib []byte

type generator struct {
	config *Config

	rand *rand.Rand

	files []*File

	varNameSeq int

	currentBlock *ir.Node

	insideLoop bool

	scope *scope

	symtab *symbolTable
	expr   *exprGenerator
}

func newGenerator(config *Config) *generator {
	symtab := newSymbolTable()
	{
		coreFuncs := phpfunc.GetList()
		for _, fn := range coreFuncs {
			symtab.AddFunc(fn)
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
	var mainFileRequires []*ir.RootRequire

	mainFileRequires = append(mainFileRequires, &ir.RootRequire{Path: "fuzzlib.php"})
	runtimeFiles := []*RuntimeFile{
		{Name: "fuzzlib.php", Contents: phpFuzzlib},
	}

	numLibs := randutil.IntRange(g.rand, 3, 5)
	for i := 0; i < numLibs; i++ {
		filename := fmt.Sprintf("lib%d.php", i)
		g.files = append(g.files, g.createLibFile(filename))
		mainFileRequires = append(mainFileRequires, &ir.RootRequire{Path: filename})
	}
	mainFile := g.createMainFile(mainFileRequires)
	g.files = append(g.files, mainFile)
	return &Program{
		Files:        g.files,
		RuntimeFiles: runtimeFiles,
	}
}

func (g *generator) createLibFile(filename string) *File {
	file := &File{Name: filename}

	funcPrefix := strings.TrimSuffix(filename, ".php")

	numLibFuncs := randutil.IntRange(g.rand, 2, 4)
	for i := 0; i < numLibFuncs; i++ {
		funcName := fmt.Sprintf("%s_func%d", funcPrefix, i)
		fn := g.createFunc(funcName, true)
		file.Nodes = append(file.Nodes, fn)
		g.symtab.AddFunc(fn.Type)
	}

	return file
}

func (g *generator) createMainFile(requires []*ir.RootRequire) *File {
	file := &File{
		Name: "main.php",
	}

	for _, r := range requires {
		file.Nodes = append(file.Nodes, r)
	}

	funcs := make([]*ir.RootFuncDecl, randutil.IntRange(g.rand, 2, 4))
	for i := range funcs {
		funcs[i] = g.createFunc("func"+strconv.Itoa(i), false)
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
		file.Nodes = append(file.Nodes, fn)
	}
	file.Nodes = append(file.Nodes, mainFunc)

	file.Nodes = append(file.Nodes, &ir.RootStmt{
		X: ir.NewCall(ir.NewName("main")),
	})

	return file
}

func (g *generator) createFunc(name string, isLibFunc bool) *ir.RootFuncDecl {
	fn := &ir.RootFuncDecl{
		Body: ir.NewBlock(),
	}

	if isLibFunc {
		fn.Type = &ir.FuncType{
			Result: g.expr.PickType(),
		}
		numParams := randutil.IntRange(g.rand, 0, 10)
		for i := 0; i < numParams; i++ {
			paramName := fmt.Sprintf("p%d", i)
			param := ir.TypeField{Name: paramName, Type: g.expr.PickType()}
			fn.Tags = append(fn.Tags, &phpdoc.ParamTag{
				VarName: "$" + paramName,
				Type:    param.Type.String(),
			})
			fn.Type.Params = append(fn.Type.Params, param)
		}
		fn.Type.MinArgsNum = len(fn.Type.Params)
		fn.Tags = append(fn.Tags, &phpdoc.ReturnTag{Type: fn.Type.Result.String()})
	} else {
		fn.Type = &ir.FuncType{
			Name:   name,
			Result: ir.VoidType,
		}
	}
	fn.Type.Name = name

	g.scope.Enter()
	for _, param := range fn.Type.Params {
		g.scope.PushVar(param.Name, param.Type)
	}
	defer func() {
		g.scope.Leave()
		if len(g.scope.depths) != 0 {
			panic("corrupted scope stack?")
		}
	}()

	g.varNameSeq = 0
	g.currentBlock = fn.Body

	numBlockVars := 0
	if isLibFunc {
		numBlockVars = randutil.IntRange(g.rand, 0, 2)
	} else {
		numBlockVars = randutil.IntRange(g.rand, 3, 7)
	}
	blockVars := make([]string, numBlockVars)
	for i := range blockVars {
		blockVars[i] = g.genVarname()
		g.pushVarDecl(blockVars[i])
	}
	numStatements := 0
	if isLibFunc {
		numStatements = randutil.IntRange(g.rand, 1, 3)
	} else {
		numStatements = randutil.IntRange(g.rand, 3, 10)
	}
	for i := 0; i < numStatements; i++ {
		g.pushStatement()
	}

	if isLibFunc {
		ret := ir.NewReturn(g.expr.GenerateValueOfType(fn.Type.Result))
		g.currentBlock.Args = append(g.currentBlock.Args, ret)
	} else {
		for _, name := range blockVars {
			v := g.scope.FindVarByName(name)
			varDump := ir.NewCall(ir.NewName("var_dump"), ir.NewVar(name, v.typ))
			g.currentBlock.Args = append(g.currentBlock.Args, varDump)
		}
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

	switch randutil.IntRange(g.rand, 0, 8) {
	case 0, 1, 2:
		g.pushVarDecl(g.genVarname())
	case 3:
		if g.insideLoop {
			g.currentBlock.Args = append(g.currentBlock.Args, ir.NewBreak(0))
		} else {
			g.pushBlockStmt()
		}
	case 4:
		if g.insideLoop {
			g.currentBlock.Args = append(g.currentBlock.Args, ir.NewContinue(0))
		} else {
			g.pushIfStmt()
		}
	case 5:
		g.pushVarDump()
	case 6, 7:
		g.pushAssignStmt()
	case 8:
		g.pushLoop()
	}
}

func (g *generator) pushLoop() {
	prevInLoop := g.insideLoop
	prevCurrentBlock := g.currentBlock
	g.insideLoop = true
	g.scope.Enter()

	iterVarName := g.genVarname()
	iterVar := ir.NewVar(iterVarName, ir.IntType)
	iterVarAssign := ir.NewAssign(iterVar, ir.NewIntLit(0))
	g.currentBlock.Args = append(g.currentBlock.Args, iterVarAssign)
	loopCond := ir.NewLess(ir.NewPostInc(iterVar), ir.NewIntLit(int64(randutil.IntRange(g.rand, 1, 10))))
	whileNode := &ir.Node{Op: ir.OpWhile}
	whileNode.Args = append(whileNode.Args, loopCond)
	g.currentBlock = whileNode
	switch randutil.IntRange(g.rand, 0, 3) {
	case 0, 1, 2:
		g.pushBlockStmt()
	case 3:
		g.pushStatement()
	}

	g.scope.Leave()
	g.insideLoop = prevInLoop
	g.currentBlock = prevCurrentBlock
	g.currentBlock.Args = append(g.currentBlock.Args, whileNode)
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
