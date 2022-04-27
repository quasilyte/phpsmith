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

	Params []ParamInfo

	ResultType     Type
	ResultTypeHint string

	Tags []phpdoc.Tag

	Body *Node
}

type ParamInfo struct {
	Name string

	Type Type

	TypeHint string
}

func (n *RootStmt) rootNode()     {}
func (n *RootFuncDecl) rootNode() {}
