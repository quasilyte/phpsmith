package ir

type RootNode interface {
	rootNode()
}

type RootRequire struct {
	Path string
}

type RootStmt struct {
	X *Node
}

type RootFuncDecl struct {
	Type *FuncType

	Body *Node
}

type RootClassDecl struct {
	Type *ClassType

	Methods []*RootFuncDecl
}

func (n *RootRequire) rootNode()   {}
func (n *RootStmt) rootNode()      {}
func (n *RootFuncDecl) rootNode()  {}
func (n *RootClassDecl) rootNode() {}
