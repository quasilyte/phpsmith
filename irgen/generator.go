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

	stmtDepth int

	varNameSeq int

	currentBlock *ir.Node

	insideLoop bool

	scope *scope

	symtab *symbolTable
	expr   *exprGenerator
}

type fileTemplate struct {
	name       string
	classTypes []*ir.ClassType
	funcTypes  []*ir.FuncType
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

	// We collect all the declarations before generating any IR.
	// This is needed to finalize the types information.
	var fileTemplates []fileTemplate

	numClasses := randutil.IntRange(g.rand, 7, 10)
	// First, declare all the classes without setting their fields or methods.
	for i := 0; i < numClasses; i++ {
		className := fmt.Sprintf("Class%d", i)
		g.symtab.DeclareClass(className)
	}
	// Now that all classes can reference each other, generate their types.
	// We have class types available for the fields at this point.
	for _, c := range g.symtab.classes {
		fileName := c.Name + ".php"
		fileTemplates = append(fileTemplates, g.createClassFileTemplate(c.Name, fileName))
	}

	numLibs := randutil.IntRange(g.rand, 3, 5)
	for i := 0; i < numLibs; i++ {
		fileName := fmt.Sprintf("lib%d.php", i)
		fileTemplates = append(fileTemplates, g.createLibFileTemplate(fileName))
	}

	// Define/add all symbols.
	for _, ft := range fileTemplates {
		for _, c := range ft.classTypes {
			g.symtab.DefineClass(c)
		}
		for _, f := range ft.funcTypes {
			g.symtab.AddFunc(f)
		}
	}

	// We can finalize the symtab now.
	// No new symbols should be added after this point.
	g.symtab.Sort()

	// Next generate the actual IR for the file templates.
	for _, ft := range fileTemplates {
		g.files = append(g.files, g.createFile(ft))
		mainFileRequires = append(mainFileRequires, &ir.RootRequire{Path: ft.name})
	}

	mainFile := g.createMainFile(mainFileRequires)
	g.files = append(g.files, mainFile)
	return &Program{
		Files:        g.files,
		RuntimeFiles: runtimeFiles,
	}
}

func (g *generator) createClassFileTemplate(className, fileName string) fileTemplate {
	return fileTemplate{
		name: fileName,
		classTypes: []*ir.ClassType{
			g.createClassType(className),
		},
	}
}

func (g *generator) createFile(ft fileTemplate) *File {
	file := &File{Name: ft.name}

	for _, c := range ft.classTypes {
		decl := &ir.RootClassDecl{
			Type: c,
		}
		for _, m := range c.Methods {
			decl.Methods = append(decl.Methods, g.createFunc(m))
		}
		file.Nodes = append(file.Nodes, decl)
	}

	for _, f := range ft.funcTypes {
		file.Nodes = append(file.Nodes, g.createFunc(f))
	}

	return file
}

func (g *generator) createClassType(classname string) *ir.ClassType {
	numFields := randutil.IntRange(g.rand, 3, 8)
	numMethods := randutil.IntRange(g.rand, 3, 5)
	c := &ir.ClassType{
		Name:    classname,
		Fields:  make([]ir.TypeField, numFields),
		Methods: make([]*ir.FuncType, numMethods),
	}
	for i := range c.Fields {
		field := &c.Fields[i]
		field.Name = fmt.Sprintf("field%d", i)
		field.Type = g.expr.PickType()
		if canConstexprInitialize(field.Type) && randutil.Chance(g.rand, 0.8) {
			field.Init = g.expr.GenerateConstValueOfType(field.Type)
		}
	}
	for i := range c.Methods {
		fn := g.createFuncType(fmt.Sprintf("method%d", i), true, c)
		c.Methods[i] = fn
		g.symtab.AddFunc(fn)
	}
	return c
}

func (g *generator) createLibFileTemplate(fileName string) fileTemplate {
	ft := fileTemplate{
		name: fileName,
	}
	funcPrefix := strings.TrimSuffix(fileName, ".php")
	numLibFuncs := randutil.IntRange(g.rand, 3, 5)
	for i := 0; i < numLibFuncs; i++ {
		funcName := fmt.Sprintf("%s_func%d", funcPrefix, i)
		funcType := g.createFuncType(funcName, true, nil)
		ft.funcTypes = append(ft.funcTypes, funcType)
	}
	return ft
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
		funcType := g.createFuncType("func"+strconv.Itoa(i), false, nil)
		funcs[i] = g.createFunc(funcType)
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

func (g *generator) createFuncType(name string, isLibFunc bool, classType *ir.ClassType) *ir.FuncType {
	var fn *ir.FuncType

	if isLibFunc {
		fn = &ir.FuncType{
			Name:      name,
			Result:    g.expr.PickType(),
			Class:     classType,
			IsLibFunc: true,
		}
		maxParams := 10
		if classType != nil {
			maxParams = 6
		}
		numParams := randutil.IntRange(g.rand, 0, maxParams)
		for i := 0; i < numParams; i++ {
			paramName := fmt.Sprintf("p%d", i)
			param := ir.TypeField{Name: paramName, Type: g.expr.PickType()}
			fn.Tags = append(fn.Tags, &phpdoc.ParamTag{
				VarName: "$" + paramName,
				Type:    param.Type.String(),
			})
			fn.Params = append(fn.Params, param)
		}
		fn.MinArgsNum = len(fn.Params)
		fn.Tags = append(fn.Tags, &phpdoc.ReturnTag{Type: fn.Result.String()})
	} else {
		fn = &ir.FuncType{
			Name:   name,
			Result: ir.VoidType,
		}
	}
	return fn
}

func (g *generator) createFunc(funcType *ir.FuncType) *ir.RootFuncDecl {
	fn := &ir.RootFuncDecl{
		Body: ir.NewBlock(),
		Type: funcType,
	}

	g.scope.Enter()
	if funcType.Class != nil {
		g.scope.PushParam("this", funcType.Class)
	}
	for _, param := range fn.Type.Params {
		g.scope.PushParam(param.Name, param.Type)
	}
	g.scope.Enter()
	defer func() {
		g.scope.Leave()
		g.scope.Leave()
		if len(g.scope.depths) != 0 {
			panic("corrupted scope stack?")
		}
	}()

	g.varNameSeq = 0
	g.currentBlock = fn.Body

	if funcType.IsLibFunc {
		call := ir.NewCall(ir.NewName("_visit_function"), ir.NewStringLit(fn.Type.FullName()))
		ret := ir.NewReturn(g.expr.GenerateConstValueOfType(fn.Type.Result))
		funcCallGuard := ir.NewIf(ir.NewNot(call), ret)
		g.currentBlock.Args = append(g.currentBlock.Args, funcCallGuard)
	}

	numBlockVars := 0
	if funcType.IsLibFunc {
		numBlockVars = randutil.IntRange(g.rand, 0, 2)
	} else {
		numBlockVars = randutil.IntRange(g.rand, 3, 7)
	}
	blockVars := make([]string, numBlockVars)
	for i := range blockVars {
		blockVars[i] = g.genVarname(false)
		g.pushVarDecl(blockVars[i])
	}
	numStatements := 0
	if funcType.IsLibFunc {
		numStatements = randutil.IntRange(g.rand, 1, 3)
	} else {
		numStatements = randutil.IntRange(g.rand, 3, 10)
	}
	for i := 0; i < numStatements; i++ {
		g.pushStatement()
	}

	if funcType.IsLibFunc {
		ret := ir.NewReturn(g.expr.GenerateValueOfType(fn.Type.Result))
		g.currentBlock.Args = append(g.currentBlock.Args, ret)
	} else {
		for _, name := range blockVars {
			v := g.scope.FindVarByName(name)
			if canDump(v.typ) {
				varDump := g.varDumpCall(ir.NewVar(name, v.typ))
				g.currentBlock.Args = append(g.currentBlock.Args, varDump)
			}
		}
	}

	return fn
}

func (g *generator) genVarname(internal bool) string {
	var varname string
	if internal {
		varname = "_iv" + strconv.Itoa(g.varNameSeq)
	} else {
		varname = "v" + strconv.Itoa(g.varNameSeq)
	}
	g.varNameSeq++
	return varname
}

func (g *generator) pushVarDecl(name string) {
	typ := g.expr.PickType()
	lhs := ir.NewVar(name, typ)
	rhs := g.expr.GenerateValueOfType(typ)
	if scalarType, ok := typ.(*ir.ScalarType); ok {
		switch scalarType.Kind {
		case ir.ScalarFloat, ir.ScalarInt:
			rhs = &ir.Node{Op: ir.OpCast, Args: []*ir.Node{rhs}, Type: typ}
		}
	}
	assign := ir.NewAssign(lhs, rhs)
	if scalarType, ok := typ.(*ir.ScalarType); ok && scalarType.Kind == ir.ScalarBool {
		assign.Value = &phpdoc.VarTag{VarName: "$" + name, Type: "bool"}
	}
	g.currentBlock.Args = append(g.currentBlock.Args, assign)
	g.scope.PushVar(name, typ)
}

func (g *generator) pushStatement() {
	g.stmtDepth++
	defer func() {
		g.stmtDepth--
	}()

	switch randutil.IntRange(g.rand, 0, 10+(g.stmtDepth*2)) {
	case 0:
		if g.insideLoop {
			g.currentBlock.Args = append(g.currentBlock.Args, ir.NewBreak(0))
		} else {
			g.pushBlockStmt()
		}
	case 1:
		if g.insideLoop {
			g.currentBlock.Args = append(g.currentBlock.Args, ir.NewContinue(0))
		} else {
			g.pushIfStmt()
		}
	case 2, 3, 4:
		if !g.pushVarDump() {
			g.pushAssignStmt()
		}
	case 5, 6:
		g.pushAssignStmt()
	case 7:
		g.pushLoopStmt()
	case 8:
		g.pushSwitchStmt()
	default:
		g.pushVarDecl(g.genVarname(false))
	}
}

func (g *generator) pushSwitchStmt() {
	var tagType ir.Type
	if randutil.Chance(g.rand, 0.3) {
		tagType = g.expr.PickEnumType()
	} else {
		tagType = g.expr.PickScalarTypeNoBool()
	}
	numCases := randutil.IntRange(g.rand, 0, 10)
	hasDefault := randutil.Bool(g.rand)
	prevCurrentBlock := g.currentBlock

	g.scope.Enter()
	defer g.scope.Leave()

	tagExpr := g.expr.GenerateValueOfType(tagType)
	switchNode := &ir.Node{Op: ir.OpSwitch, Args: []*ir.Node{tagExpr}}
	caseSet := make(map[any]struct{})
	for i := 0; i < numCases; i++ {
		x := g.expr.GenerateValueOfType(tagType)
		caseValue := extractValue(x)
		if _, ok := caseSet[caseValue]; ok {
			continue
		}
		caseSet[caseValue] = struct{}{}

		g.scope.Enter()
		caseNode := &ir.Node{Op: ir.OpCase, Args: []*ir.Node{x}}
		caseSize := randutil.IntRange(g.rand, 0, 2)
		g.currentBlock = caseNode
		for j := 0; j < caseSize; j++ {
			g.pushStatement()
		}
		if randutil.IntRange(g.rand, 0, 8) != 0 {
			g.currentBlock.Args = append(g.currentBlock.Args, ir.NewBreak(0))
		}
		switchNode.Args = append(switchNode.Args, caseNode)

		g.scope.Leave()
	}
	if hasDefault {
		g.scope.Enter()
		caseNode := &ir.Node{Op: ir.OpDefaultCase}
		caseSize := randutil.IntRange(g.rand, 0, 2)
		g.currentBlock = caseNode
		for j := 0; j < caseSize; j++ {
			g.pushStatement()
		}
		switchNode.Args = append(switchNode.Args, caseNode)
		g.scope.Leave()
	}

	g.currentBlock = prevCurrentBlock
	g.currentBlock.Args = append(g.currentBlock.Args, switchNode)
}

func (g *generator) pushLoopStmt() {
	prevInLoop := g.insideLoop
	prevCurrentBlock := g.currentBlock
	g.insideLoop = true
	g.scope.Enter()

	iterVarName := g.genVarname(true)
	iterVar := ir.NewVar(iterVarName, ir.IntType)
	iterVarAssign := ir.NewAssign(iterVar, ir.NewIntLit(0))
	g.currentBlock.Args = append(g.currentBlock.Args, iterVarAssign)
	loopCond := ir.NewLess(ir.NewPostInc(iterVar), ir.NewIntLit(int64(randutil.IntRange(g.rand, 1, 10))))
	whileNode := &ir.Node{Op: ir.OpWhile}
	whileNode.Args = append(whileNode.Args, loopCond)

	g.currentBlock = whileNode
	g.pushBlockStmt()

	g.scope.Leave()
	g.insideLoop = prevInLoop
	g.currentBlock = prevCurrentBlock
	g.currentBlock.Args = append(g.currentBlock.Args, whileNode)
}

func (g *generator) pushAssignStmt() {
	lhs, typ := g.expr.PickLvalue()
	if lhs == nil {
		g.pushVarDecl(g.genVarname(false))
		return
	}
	var op ir.Op
	if typ, ok := typ.(*ir.ScalarType); ok && randutil.Bool(g.rand) {
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
	rhs := g.expr.GenerateValueOfType(typ)
	if op != ir.OpInvalid {
		assign = ir.NewAssignModify(op, lhs, rhs)
	} else {
		assign = ir.NewAssign(lhs, rhs)
	}
	g.currentBlock.Args = append(g.currentBlock.Args, assign)
}

func (g *generator) pushVarDump() bool {
	for attempts := 0; attempts < 5; attempts++ {
		typ := g.expr.PickType()
		if !canDump(typ) {
			continue
		}
		arg := g.expr.GenerateValueOfType(typ)
		varDump := g.varDumpCall(arg)
		g.currentBlock.Args = append(g.currentBlock.Args, varDump)
		return true
	}
	return false
}

func (g *generator) varDumpCall(arg *ir.Node) *ir.Node {
	file := ir.NewName("__FILE__")
	line := ir.NewName("__LINE__")
	return ir.NewCall(ir.NewName("dump_with_pos"), file, line, arg)
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

	oldBlock := g.currentBlock
	g.scope.Enter()

	newBlock := &ir.Node{Op: ir.OpBlock}
	g.currentBlock = newBlock
	g.pushStatement()
	oldBlock.Args = append(oldBlock.Args, ir.NewIf(cond, newBlock))

	g.scope.Leave()
	g.currentBlock = oldBlock
}
