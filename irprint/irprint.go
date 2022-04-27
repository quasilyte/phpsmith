package irprint

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"

	"github.com/quasilyte/phpsmith/ir"
)

// TODO:
// - printing with random formatting (using externally provided rand object)
// - testing for both modes (random and normal)?

type Config struct {
	// Rand is used to add randomized formatting to the output.
	// If nil, no randomization will be used and the output will look like pretty-printed.
	Rand *rand.Rand
}

func FprintRootNode(w io.Writer, n ir.RootNode, config *Config) {
	p := &printer{
		config: config,
		w:      bufio.NewWriter(w),
	}
	p.printRootNode(n)
	p.w.Flush()
}

func FprintNode(w io.Writer, n *ir.Node, config *Config) {
	p := &printer{
		config: config,
		w:      bufio.NewWriter(w),
	}
	p.printNode(n)
	p.w.Flush()
}

type printer struct {
	config *Config
	w      *bufio.Writer
	depth  int
}

func (p *printer) indent() {
	for i := 0; i < p.depth; i++ {
		p.w.WriteByte(' ')
	}
}

func (p *printer) printRootNode(n ir.RootNode) {
	switch n := n.(type) {
	case *ir.RootFuncDecl:
		p.printFuncDecl(n)
	case *ir.RootStmt:
		p.printNode(n.X)
	}
}

func (p *printer) printFuncDecl(decl *ir.RootFuncDecl) {
	if len(decl.Tags) != 0 {
		p.w.WriteString("/**\n")
		for _, tag := range decl.Tags {
			fmt.Fprintf(p.w, " * %s %s\n", tag.Name(), tag.Value())
		}
		p.w.WriteString(" */\n")
	}

	p.w.WriteString("function " + decl.Name)
	p.w.WriteByte('(')
	for i, param := range decl.Params {
		if i != 0 {
			p.w.WriteString(", ")
		}
		if param.TypeHint != "" {
			p.w.WriteString(param.TypeHint)
			p.w.WriteByte(' ')
		}
		p.w.WriteString("$" + param.Name)
	}
	p.w.WriteString(") ")
	p.printNode(decl.Body)
	p.w.WriteByte('\n')
}

func (p *printer) printNode(n *ir.Node) {
	switch n.Op {
	case ir.OpBlock:
		p.depth += 2
		p.w.WriteString("{\n")
		for _, stmt := range n.Args {
			p.indent()
			p.printNode(stmt)
			p.w.WriteString(";\n")
		}
		p.w.WriteString("}\n")
		p.depth -= 2

	case ir.OpEcho:
		p.w.WriteString("echo ")
		p.printNodes(n.Args, ", ")

	case ir.OpReturn:
		p.w.WriteString("return ")
		p.printNode(n.Args[0])

	case ir.OpReturnVoid:
		p.w.WriteString("return")

	case ir.OpBoolLit:
		fmt.Fprintf(p.w, "%v", n.Value)
	case ir.OpIntLit:
		fmt.Fprintf(p.w, "%#v", n.Value)
	case ir.OpFloatLit:
		if n.Value.(float64) == 0 {
			p.w.WriteString("0.0")
		} else {
			fmt.Fprintf(p.w, "%#v", n.Value)
		}
	case ir.OpStringLit:
		p.printString(n)

	case ir.OpVar:
		p.w.WriteString("$" + n.Value.(string))

	case ir.OpAssign:
		p.printBinary(n.Args, "=")

	case ir.OpAdd:
		p.printBinary(n.Args, "+")
	case ir.OpSub:
		p.printBinary(n.Args, "-")
	case ir.OpConcat:
		p.printBinary(n.Args, ".")
	}
}

func (p *printer) printBinary(args []*ir.Node, op string) {
	p.printNode(args[0])
	p.w.WriteString(" " + op + " ")
	p.printNode(args[1])
}

func (p *printer) printNodes(nodes []*ir.Node, sep string) {
	for i, n := range nodes {
		if i != 0 {
			p.w.WriteString(sep)
		}
		p.printNode(n)
	}
}

func (p *printer) printString(s *ir.Node) {
	p.w.WriteByte('\'')
	p.w.WriteString(s.Value.(string))
	p.w.WriteByte('\'')
}
