package phpfunc

import (
	"github.com/quasilyte/phpsmith/ir"
)

func Add(dst map[string]*ir.FuncType) {
	for _, f := range funcList {
		dst[f.Name] = f
	}
}

func init() {
	for _, f := range funcList {
		minArgsNum := len(f.Params)
		for i := len(f.Params) - 1; i > 0; i-- {
			p := f.Params[i]
			if p.Init != nil {
				minArgsNum--
			}
		}
		f.MinArgsNum = minArgsNum
	}
}

var funcList = []*ir.FuncType{
	{
		Name: "strlen",
		Params: []ir.TypeField{
			{Name: "string", Type: ir.StringType},
		},
		Result: ir.IntType,
	},

	{
		Name: "ltrim",
		Params: []ir.TypeField{
			{Name: "string", Type: ir.StringType},
			{Name: "characters", Type: ir.StringType, Init: " \n\r\t\v\x00"},
		},
		Result: ir.StringType,
	},
}
