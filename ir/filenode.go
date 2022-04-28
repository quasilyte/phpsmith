package ir

import (
	"github.com/quasilyte/phpsmith/phpdoc"
)

type RootNode interface {
	rootNode()
}

type RootStmt struct {
	X *Node
}

type RootFuncDecl struct {
	Name string

	Type *FuncType

	Tags []phpdoc.Tag

	Body *Node
}

func (n *RootStmt) rootNode()     {}
func (n *RootFuncDecl) rootNode() {}
