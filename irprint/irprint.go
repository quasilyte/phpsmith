package irprint

import (
	"bufio"
	"bytes"
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

type printFlags int

const (
	flagNeedSemicolon printFlags = 1 << iota
	flagNeedNewline
)

func (flags printFlags) NeedSemicolon() bool {
	return flags&flagNeedSemicolon != 0
}

func (flags printFlags) NeedNewline() bool {
	return flags&flagNeedNewline != 0
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

	p.w.WriteString("function " + decl.Type.Name)
	p.w.WriteByte('(')
	for i, param := range decl.Type.Params {
		if i != 0 {
			p.w.WriteString(", ")
		}
		// TODO: print a type hint for some types, sometimes?
		p.w.WriteString("$" + param.Name)
	}
	p.w.WriteString(") ")
	p.printNode(decl.Body)
	p.w.WriteByte('\n')
}

//nolint:gocyclo
func (p *printer) printNode(n *ir.Node) printFlags {
	switch n.Op {
	case ir.OpBlock:
		p.depth += 2
		p.w.WriteString("{\n")
		for _, stmt := range n.Args {
			p.indent()
			flags := p.printNode(stmt)
			if flags.NeedSemicolon() {
				p.w.WriteByte(';')
			}
			if flags.NeedNewline() {
				p.w.WriteString("\n")
			}
		}
		p.depth -= 2
		p.indent()
		p.w.WriteString("}\n")
		return 0

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
	case ir.OpName:
		p.w.WriteString(n.Value.(string))

	case ir.OpAssign:
		p.printBinary(n, "=")

	case ir.OpAdd:
		p.printBinary(n, "+")
	case ir.OpSub:
		p.printBinary(n, "-")
	case ir.OpConcat:
		p.printBinary(n, ".")

	case ir.OpAndWord:
		p.printBinary(n, "and")
	case ir.OpAnd:
		p.printBinary(n, "&&")
	case ir.OpXorWord:
		p.printBinary(n, "xor")
	case ir.OpOrWord:
		p.printBinary(n, "or")
	case ir.OpOr:
		p.printBinary(n, "||")

	case ir.OpEqual2:
		p.printBinary(n, "==")
	case ir.OpEqual3:
		p.printBinary(n, "===")
	case ir.OpLess:
		p.printBinary(n, "<")
	case ir.OpLessOrEqual:
		p.printBinary(n, "<=")
	case ir.OpGreater:
		p.printBinary(n, ">")
	case ir.OpGreaterOrEqual:
		p.printBinary(n, ">=")

	case ir.OpPreInc:
		p.printUnaryPrefix(n, "++")
	case ir.OpPreDec:
		p.printUnaryPrefix(n, "--")

	case ir.OpPostInc:
		p.printUnaryPostfix(n, "++")
	case ir.OpPostDec:
		p.printUnaryPostfix(n, "--")

	case ir.OpNot:
		p.printUnaryPrefix(n, "!")

	case ir.OpParens:
		p.w.WriteByte('(')
		p.printNode(n.Args[0])
		p.w.WriteByte(')')

	case ir.OpTernary:
		p.printNode(n.Args[0])
		p.w.WriteString(" ? ")
		p.printNode(n.Args[1])
		p.w.WriteString(" : ")
		p.printNode(n.Args[2])

	case ir.OpArrayLit:
		if len(n.Args) == 0 {
			p.w.WriteString("[]")
		} else {
			p.w.WriteString("[\n")
			p.depth += 2
			for _, elem := range n.Args {
				p.indent()
				p.printNode(elem)
				p.w.WriteString(",\n")
			}
			p.depth -= 2
			p.indent()
			p.w.WriteString("]")
		}

	case ir.OpCall:
		p.printNode(n.Args[0])
		p.w.WriteByte('(')
		for i, arg := range n.Args[1:] {
			if i != 0 {
				p.w.WriteString(", ")
			}
			p.printNode(arg)
		}
		p.w.WriteByte(')')

	case ir.OpCast:
		p.w.WriteByte('(')
		p.w.WriteString(n.Type.String())
		p.w.WriteByte(')')
		p.printNode(n.Args[0])

	case ir.OpIf:
		p.w.WriteString("if (")
		p.printNode(n.Args[0])
		p.w.WriteString(") ")
		return p.printNode(n.Args[1])
	}

	return flagNeedNewline | flagNeedSemicolon
}

func (p *printer) printUnaryPrefix(n *ir.Node, op string) {
	p.w.WriteString(op)
	p.printNode(n.Args[0])
}

func (p *printer) printUnaryPostfix(n *ir.Node, op string) {
	p.printNode(n.Args[0])
	p.w.WriteString(op)
}

func (p *printer) printBinary(n *ir.Node, op string) {
	p.printNode(n.Args[0])
	p.w.WriteString(" " + op + " ")
	p.printNode(n.Args[1])
}

func (p *printer) printNodes(nodes []*ir.Node, sep string) {
	for i, n := range nodes {
		if i != 0 {
			p.w.WriteString(sep)
		}
		p.printNode(n)
	}
}

func (p *printer) printString(n *ir.Node) {
	s := n.Value.(string)
	var buf bytes.Buffer
	buf.Grow(len(s))
	for i := 0; i < len(s); i++ {
		ch := s[i]
		switch ch {
		case '\r':
			buf.WriteString(`\r`)
		case '\n':
			buf.WriteString(`\n`)
		case '"':
			buf.WriteString(`\"`)
		case '\\':
			buf.WriteString(`\\`)
		case 0:
			buf.WriteString(`\000`)
		case '\a':
			buf.WriteString(`\a`)
		case '\b':
			buf.WriteString(`\b`)
		case '\f':
			buf.WriteString(`\f`)
		case '\t':
			buf.WriteString(`\t`)
		case '\v':
			buf.WriteString(`\v`)
		default:
			if ch < 32 {
				buf.WriteString(`\0`)
				buf.WriteByte('0' + ch/8)
				buf.WriteByte('0' + ch%8)
			} else {
				buf.WriteByte(ch)
			}
		}
	}

	quote := byte('"')
	p.w.WriteByte(quote)
	p.w.Write(buf.Bytes())
	p.w.WriteByte(quote)
}
