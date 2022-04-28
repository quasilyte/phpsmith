package phpfunc

import (
	"github.com/quasilyte/phpsmith/ir"
)

func Add(dst map[string]*ir.FuncType) {
	for _, f := range funcList {
		dst[f.name] = f.typ
	}
}

type funcEntry struct {
	name string
	typ  *ir.FuncType
}

var funcList = []funcEntry{
	{
		name: "strlen",
		typ: &ir.FuncType{
			Name: "strlen",
			Params: []ir.TypeField{
				{Name: "string", Type: ir.StringType},
			},
			Result: ir.IntType,
		},
	},
}
