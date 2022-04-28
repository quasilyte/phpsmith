package irprint

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/quasilyte/phpsmith/ir"
)

func TestPrintNodePretty(t *testing.T) {
	intType := &ir.ScalarType{Kind: ir.ScalarInt}

	tests := []struct {
		n    *ir.Node
		want string
	}{
		{ir.NewBoolLit(true), `true`},
		{ir.NewBoolLit(false), `false`},
		{ir.NewIntLit(0), `0`},
		{ir.NewIntLit(32), `32`},
		{ir.NewIntLit(-312), `-312`},
		{ir.NewFloatLit(0), `0.0`},
		{ir.NewFloatLit(-1.4), `-1.4`},
		{ir.NewAssignModify(ir.OpAdd, ir.NewVar("x", intType), ir.NewVar("y", intType)), "$x += $y"},
		// TODO: more string tests when printer handles them correctly.
		{ir.NewStringLit(""), `""`},
		{ir.NewStringLit("123"), `"123"`},
		{ir.NewStringLit("\\n"), `"\\n"`},

		{ir.NewEcho(ir.NewVar("foo", intType)), `echo $foo`},
		{ir.NewEcho(ir.NewVar("foo", intType), ir.NewBoolLit(false)), `echo $foo, false`},

		{ir.NewAdd(ir.NewIntLit(1), ir.NewIntLit(2)), `1 + 2`},
		{ir.NewSub(ir.NewIntLit(1), ir.NewIntLit(2)), `1 - 2`},

		{ir.NewReturn(ir.NewVar("x", intType)), "return $x"},
		{ir.NewReturnVoid(), "return"},

		{
			ir.NewBlock(ir.NewEcho(ir.NewStringLit("ok"))),
			`{
  echo "ok";
}
`,
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(fmt.Sprintf("test%d", i), func(t *testing.T) {
			var buf bytes.Buffer
			config := &Config{}
			FprintNode(&buf, test.n, config)
			have := buf.String()
			if have != test.want {
				t.Fatalf("print %s:\nhave: %q\nwant: %q", test.n.Op, have, test.want)
			}
		})
	}
}
