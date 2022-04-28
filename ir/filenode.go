package ir

import (
	"github.com/quasilyte/phpsmith/phpdoc"
)

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

	Tags []phpdoc.Tag

	Body *Node
}

func (n *RootRequire) rootNode()  {}
func (n *RootStmt) rootNode()     {}
func (n *RootFuncDecl) rootNode() {}
