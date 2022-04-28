package phpfunc

import (
	"github.com/quasilyte/phpsmith/ir"
)

func GetList() []*ir.FuncType {
	return funcList
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
		Name: "json_encode",
		Params: []ir.TypeField{
			{Name: "value", Type: ir.MixedType},
		},
		Result:   ir.StringType,
		NeedCast: true,
	},
	{
		Name: "strtolower",
		Params: []ir.TypeField{
			{Name: "str", Type: ir.StringType},
		},
		Result: ir.StringType,
	},
	{
		Name: "rtrim",
		Params: []ir.TypeField{
			{Name: "s", Type: ir.StringType},
			{Name: "what", Type: ir.StringType, Init: " \x0a\x0d\x09\x0b\x00"},
		},
		Result: ir.StringType,
	},
	{
		Name: "strcasecmp",
		Params: []ir.TypeField{
			{Name: "str1", Type: ir.StringType},
			{Name: "str2", Type: ir.StringType},
		},
		Result: ir.IntType,
	},
	{
		Name: "strlen",
		Params: []ir.TypeField{
			{Name: "str", Type: ir.StringType},
		},
		Result: ir.IntType,
	},
	{
		Name: "ord",
		Params: []ir.TypeField{
			{Name: "c", Type: ir.StringType},
		},
		Result: ir.IntType,
	},
	{
		Name: "ltrim",
		Params: []ir.TypeField{
			{Name: "s", Type: ir.StringType},
			{Name: "what", Type: ir.StringType, Init: " \x0a\x0d\x09\x0b\x00"},
		},
		Result: ir.StringType,
	},
	{
		Name: "strnatcmp",
		Params: []ir.TypeField{
			{Name: "str1", Type: ir.StringType},
			{Name: "str2", Type: ir.StringType},
		},
		Result: ir.IntType,
	},
	{
		Name: "ucfirst",
		Params: []ir.TypeField{
			{Name: "str", Type: ir.StringType},
		},
		Result: ir.StringType,
	},
	{
		Name: "strtoupper",
		Params: []ir.TypeField{
			{Name: "str", Type: ir.StringType},
		},
		Result: ir.StringType,
	},
	{
		Name: "lcfirst",
		Params: []ir.TypeField{
			{Name: "str", Type: ir.StringType},
		},
		Result: ir.StringType,
	},
	// {
	// 	Name: "str_starts_with",
	// 	Params: []ir.TypeField{
	// 		{Name: "haystack", Type: ir.StringType},
	// 		{Name: "needle", Type: ir.StringType},
	// 	},
	// 	Result: ir.BoolType,
	// },
	// {
	// 	Name: "vprintf",
	// 	Params: []ir.TypeField{
	// 		{Name: "format", Type: ir.StringType},
	// 		{Name: "args", Type: &ir.ArrayType{Elem: ir.MixedType}},
	// 	},
	// 	Result: ir.IntType,
	// },
	{
		Name: "ucwords",
		Params: []ir.TypeField{
			{Name: "str", Type: ir.StringType},
		},
		Result: ir.StringType,
	},
	{
		Name: "strrev",
		Params: []ir.TypeField{
			{Name: "str", Type: ir.StringType},
		},
		Result: ir.StringType,
	},
	{
		Name: "substr_replace",
		Params: []ir.TypeField{
			{Name: "str", Type: ir.StringType},
			{Name: "replacement", Type: ir.StringType},
			{Name: "start", Type: ir.IntType, Strict: true},
			{Name: "length", Type: ir.IntType, Init: 9223372036854775807, Strict: true},
		},
		Result: ir.StringType,
	},
	{
		Name: "trim",
		Params: []ir.TypeField{
			{Name: "s", Type: ir.StringType},
			{Name: "what", Type: ir.StringType, Init: " \x0a\x0d\x09\x0b\x00"},
		},
		Result: ir.StringType,
	},
	// {
	// 	Name: "vsprintf",
	// 	Params: []ir.TypeField{
	// 		{Name: "format", Type: ir.StringType},
	// 		{Name: "args", Type: &ir.ArrayType{Elem: ir.MixedType}},
	// 	},
	// 	Result: ir.StringType,
	// },
	{
		Name: "str_split",
		Params: []ir.TypeField{
			{Name: "str", Type: ir.StringType},
			{Name: "split_length", Type: ir.IntType, Init: 1},
		},
		Result: &ir.ArrayType{Elem: ir.StringType},
	},
	{
		Name: "is_dir",
		Params: []ir.TypeField{
			{Name: "name", Type: ir.StringType},
		},
		Result: ir.BoolType,
	},
	{
		Name: "getimagesize",
		Params: []ir.TypeField{
			{Name: "name", Type: ir.StringType},
		},
		Result: ir.MixedType,
	},
	// {
	// 	Name: "str_ends_with",
	// 	Params: []ir.TypeField{
	// 		{Name: "haystack", Type: ir.StringType},
	// 		{Name: "needle", Type: ir.StringType},
	// 	},
	// 	Result: ir.BoolType,
	// },
	{
		Name: "is_writeable",
		Params: []ir.TypeField{
			{Name: "name", Type: ir.StringType},
		},
		Result: ir.BoolType,
	},
	// {
	// 	Name: "str_repeat",
	// 	Params: []ir.TypeField{
	// 		{Name: "s", Type: ir.StringType},
	// 		{Name: "multiplier", Type: ir.IntType},
	// 	},
	// 	Result: ir.StringType,
	// },
	{
		Name: "basename",
		Params: []ir.TypeField{
			{Name: "name", Type: ir.StringType},
			{Name: "suffix", Type: ir.StringType, Init: ""},
		},
		Result: ir.StringType,
	},
	{
		Name: "dirname",
		Params: []ir.TypeField{
			{Name: "name", Type: ir.StringType},
		},
		Result: ir.StringType,
	},
	{
		Name: "levenshtein",
		Params: []ir.TypeField{
			{Name: "str1", Type: ir.StringType},
			{Name: "str2", Type: ir.StringType},
		},
		Result: ir.IntType,
	},
	{
		Name: "file_exists",
		Params: []ir.TypeField{
			{Name: "name", Type: ir.StringType},
		},
		Result: ir.BoolType,
	},
	{
		Name: "chr",
		Params: []ir.TypeField{
			{Name: "v", Type: ir.IntType, Strict: true},
		},
		Result: ir.StringType,
	},
	{
		Name: "is_readable",
		Params: []ir.TypeField{
			{Name: "name", Type: ir.StringType},
		},
		Result: ir.BoolType,
	},
	{
		Name: "parse_str",
		Params: []ir.TypeField{
			{Name: "str", Type: ir.StringType},
			{Name: "arr", Type: ir.MixedType},
		},
		Result: ir.VoidType,
	},
	{
		Name:   "pi",
		Params: []ir.TypeField{},
		Result: ir.FloatType,
	},
	// {
	// 	Name: "htmlspecialchars",
	// 	Params: []ir.TypeField{
	// 		{Name: "str", Type: ir.StringType},
	// 		{Name: "flags", Type: ir.IntType, Init: 0},
	// 	},
	// 	Result: ir.StringType,
	// },
	{
		Name: "strcmp",
		Params: []ir.TypeField{
			{Name: "str1", Type: ir.StringType},
			{Name: "str2", Type: ir.StringType},
		},
		Result: ir.IntType,
	},
	{
		Name: "is_file",
		Params: []ir.TypeField{
			{Name: "name", Type: ir.StringType},
		},
		Result: ir.BoolType,
	},
	{
		Name: "long2ip",
		Params: []ir.TypeField{
			{Name: "ip", Type: ir.IntType, Strict: true},
		},
		Result: ir.StringType,
	},
	{
		Name: "stripslashes",
		Params: []ir.TypeField{
			{Name: "str", Type: ir.StringType},
		},
		Result: ir.StringType,
	},
	// {
	// 	Name: "log",
	// 	Params: []ir.TypeField{
	// 		{Name: "v", Type: ir.FloatType},
	// 		{Name: "base", Type: ir.FloatType, Init: 2.7182818284590452353602874713527},
	// 	},
	// 	Result: ir.FloatType,
	// },
	{
		Name: "htmlentities",
		Params: []ir.TypeField{
			{Name: "str", Type: ir.StringType},
		},
		Result: ir.StringType,
	},
	// {
	// 	Name: "wordwrap",
	// 	Params: []ir.TypeField{
	// 		{Name: "str", Type: ir.StringType},
	// 		{Name: "width", Type: ir.IntType, Init: 75},
	// 		{Name: "break", Type: ir.StringType, Init: "\n"},
	// 		{Name: "cut", Type: ir.BoolType, Init: false},
	// 	},
	// 	Result: ir.StringType,
	// },
	{
		Name: "deg2rad",
		Params: []ir.TypeField{
			{Name: "v", Type: ir.FloatType},
		},
		Result: ir.FloatType,
	},
	{
		Name: "cosh",
		Params: []ir.TypeField{
			{Name: "v", Type: ir.FloatType},
		},
		Result: ir.FloatType,
	},
	// {
	// 	Name: "hexdec",
	// 	Params: []ir.TypeField{
	// 		{Name: "number", Type: ir.StringType},
	// 	},
	// 	Result: ir.IntType,
	// },
	// {
	// 	Name: "bindec",
	// 	Params: []ir.TypeField{
	// 		{Name: "number", Type: ir.StringType},
	// 	},
	// 	Result: ir.IntType,
	// },
	// {
	// 	Name: "dechex",
	// 	Params: []ir.TypeField{
	// 		{Name: "number", Type: ir.IntType},
	// 	},
	// 	Result: ir.StringType,
	// },
	{
		Name: "sqrt",
		Params: []ir.TypeField{
			{Name: "v", Type: ir.FloatType},
		},
		Result: ir.FloatType,
	},
	{
		Name: "addslashes",
		Params: []ir.TypeField{
			{Name: "str", Type: ir.StringType},
		},
		Result: ir.StringType,
	},
	{
		Name: "tan",
		Params: []ir.TypeField{
			{Name: "v", Type: ir.FloatType},
		},
		Result: ir.FloatType,
	},
	{
		Name: "sin",
		Params: []ir.TypeField{
			{Name: "v", Type: ir.FloatType},
		},
		Result: ir.FloatType,
	},
	{
		Name: "rad2deg",
		Params: []ir.TypeField{
			{Name: "v", Type: ir.FloatType},
		},
		Result: ir.FloatType,
	},
	{
		Name: "cos",
		Params: []ir.TypeField{
			{Name: "v", Type: ir.FloatType},
		},
		Result: ir.FloatType,
	},
	{
		Name: "addcslashes",
		Params: []ir.TypeField{
			{Name: "str", Type: ir.StringType},
			{Name: "what", Type: ir.StringType},
		},
		Result: ir.StringType,
	},
	{
		Name: "ceil",
		Params: []ir.TypeField{
			{Name: "v", Type: ir.FloatType},
		},
		Result: ir.FloatType,
	},
	{
		Name: "bin2hex",
		Params: []ir.TypeField{
			{Name: "str", Type: ir.StringType},
		},
		Result: ir.StringType,
	},
	{
		Name: "exp",
		Params: []ir.TypeField{
			{Name: "v", Type: ir.FloatType},
		},
		Result: ir.FloatType,
	},
	// {
	// 	Name: "html_entity_decode",
	// 	Params: []ir.TypeField{
	// 		{Name: "str", Type: ir.StringType},
	// 		{Name: "flags", Type: ir.IntType, Init: 0, Strict: true},
	// 		{Name: "encoding", Type: ir.StringType, Init: "1251"},
	// 	},
	// 	Result: ir.StringType,
	// },
	{
		Name: "acosh",
		Params: []ir.TypeField{
			{Name: "v", Type: ir.FloatType},
		},
		Result: ir.FloatType,
	},
	{
		Name: "atan",
		Params: []ir.TypeField{
			{Name: "v", Type: ir.FloatType},
		},
		Result: ir.FloatType,
	},
	{
		Name: "floor",
		Params: []ir.TypeField{
			{Name: "v", Type: ir.FloatType},
		},
		Result: ir.FloatType,
	},
	{
		Name: "round",
		Params: []ir.TypeField{
			{Name: "v", Type: ir.FloatType},
			{Name: "precision", Type: ir.IntType, Init: 0, Strict: true},
		},
		Result: ir.FloatType,
	},
	{
		Name: "fmod",
		Params: []ir.TypeField{
			{Name: "x", Type: ir.FloatType},
			{Name: "y", Type: ir.FloatType},
		},
		Result: ir.FloatType,
	},
	{
		Name: "count_chars",
		Params: []ir.TypeField{
			{Name: "str", Type: ir.StringType},
			{Name: "mode", Type: ir.IntType, Init: 0},
		},
		Result: ir.MixedType,
	},
	{
		Name: "decbin",
		Params: []ir.TypeField{
			{Name: "number", Type: ir.IntType, Strict: true},
		},
		Result: ir.StringType,
	},
	{
		Name: "asinh",
		Params: []ir.TypeField{
			{Name: "v", Type: ir.FloatType},
		},
		Result: ir.FloatType,
	},
	// {
	// 	Name: "htmlspecialchars_decode",
	// 	Params: []ir.TypeField{
	// 		{Name: "str", Type: ir.StringType},
	// 		{Name: "flags", Type: ir.IntType, Init: 0, Strict: true},
	// 	},
	// 	Result: ir.StringType,
	// },
	// {
	// 	Name: "base_convert",
	// 	Params: []ir.TypeField{
	// 		{Name: "number", Type: ir.StringType},
	// 		{Name: "frombase", Type: ir.IntType},
	// 		{Name: "tobase", Type: ir.IntType},
	// 	},
	// 	Result: ir.StringType,
	// },
	{
		Name: "parse_url",
		Params: []ir.TypeField{
			{Name: "str", Type: ir.StringType},
			{Name: "component", Type: ir.IntType, Init: -1},
		},
		Result: ir.MixedType,
	},
	{
		Name: "sinh",
		Params: []ir.TypeField{
			{Name: "v", Type: ir.FloatType},
		},
		Result: ir.FloatType,
	},
	{
		Name: "asin",
		Params: []ir.TypeField{
			{Name: "v", Type: ir.FloatType},
		},
		Result: ir.FloatType,
	},
	{
		Name: "urlencode",
		Params: []ir.TypeField{
			{Name: "str", Type: ir.StringType},
		},
		Result: ir.StringType,
	},
	// {
	// 	Name: "floatval",
	// 	Params: []ir.TypeField{
	// 		{Name: "v", Type: ir.MixedType},
	// 	},
	// 	Result: ir.FloatType,
	// },
	{
		Name: "rawurldecode",
		Params: []ir.TypeField{
			{Name: "str", Type: ir.StringType},
		},
		Result: ir.StringType,
	},
	// {
	// 	Name: "intval",
	// 	Params: []ir.TypeField{
	// 		{Name: "v", Type: ir.MixedType},
	// 	},
	// 	Result: ir.IntType,
	// },
	// {
	// 	Name: "boolval",
	// 	Params: []ir.TypeField{
	// 		{Name: "v", Type: ir.MixedType},
	// 	},
	// 	Result: ir.BoolType,
	// },
	{
		Name: "acos",
		Params: []ir.TypeField{
			{Name: "v", Type: ir.FloatType},
		},
		Result: ir.FloatType,
	},
	{
		Name: "rawurlencode",
		Params: []ir.TypeField{
			{Name: "str", Type: ir.StringType},
		},
		Result: ir.StringType,
	},
	{
		Name: "atan2",
		Params: []ir.TypeField{
			{Name: "y", Type: ir.FloatType},
			{Name: "x", Type: ir.FloatType},
		},
		Result: ir.FloatType,
	},
	{
		Name: "sha1",
		Params: []ir.TypeField{
			{Name: "s", Type: ir.StringType},
			{Name: "raw_output", Type: ir.BoolType, Init: false},
		},
		Result: ir.StringType,
	},
	{
		Name: "base64_encode",
		Params: []ir.TypeField{
			{Name: "str", Type: ir.StringType},
		},
		Result: ir.StringType,
	},
	{
		Name: "urldecode",
		Params: []ir.TypeField{
			{Name: "str", Type: ir.StringType},
		},
		Result: ir.StringType,
	},
	{
		Name: "md5",
		Params: []ir.TypeField{
			{Name: "s", Type: ir.StringType},
			{Name: "raw_output", Type: ir.BoolType, Init: false},
		},
		Result: ir.StringType,
	},
	{
		Name: "natsort",
		Params: []ir.TypeField{
			{Name: "a", Type: &ir.ArrayType{Elem: ir.MixedType}},
		},
		Result: ir.VoidType,
	},
	{
		Name: "shuffle",
		Params: []ir.TypeField{
			{Name: "a", Type: &ir.ArrayType{Elem: ir.MixedType}},
		},
		Result: ir.VoidType,
	},
	{
		Name: "crc32",
		Params: []ir.TypeField{
			{Name: "s", Type: ir.StringType},
		},
		Result: ir.IntType,
	},
	{
		Name: "preg_quote",
		Params: []ir.TypeField{
			{Name: "str", Type: ir.StringType},
			{Name: "delimiter", Type: ir.StringType, Init: ""},
		},
		Result: ir.StringType,
	},
	{
		Name: "is_object",
		Params: []ir.TypeField{
			{Name: "v", Type: ir.MixedType},
		},
		Result: ir.BoolType,
	},
	{
		Name: "checkdate",
		Params: []ir.TypeField{
			{Name: "month", Type: ir.IntType},
			{Name: "day", Type: ir.IntType},
			{Name: "year", Type: ir.IntType},
		},
		Result: ir.BoolType,
	},
	// {
	// 	Name: "is_double",
	// 	Params: []ir.TypeField{
	// 		{Name: "v", Type: ir.MixedType},
	// 	},
	// 	Result: ir.BoolType,
	// },
	// {
	// 	Name: "is_float",
	// 	Params: []ir.TypeField{
	// 		{Name: "v", Type: ir.MixedType},
	// 	},
	// 	Result: ir.BoolType,
	// },
	// {
	// 	Name: "is_int",
	// 	Params: []ir.TypeField{
	// 		{Name: "v", Type: ir.MixedType},
	// 	},
	// 	Result: ir.BoolType,
	// },
	// {
	// 	Name: "is_null",
	// 	Params: []ir.TypeField{
	// 		{Name: "v", Type: ir.MixedType},
	// 	},
	// 	Result: ir.BoolType,
	// },
	// {
	// 	Name: "is_numeric",
	// 	Params: []ir.TypeField{
	// 		{Name: "v", Type: ir.MixedType},
	// 	},
	// 	Result: ir.BoolType,
	// },
	// {
	// 	Name: "is_string",
	// 	Params: []ir.TypeField{
	// 		{Name: "v", Type: ir.MixedType},
	// 	},
	// 	Result: ir.BoolType,
	// },
	// {
	// 	Name: "is_integer",
	// 	Params: []ir.TypeField{
	// 		{Name: "v", Type: ir.MixedType},
	// 	},
	// 	Result: ir.BoolType,
	// },
	// {
	// 	Name: "is_long",
	// 	Params: []ir.TypeField{
	// 		{Name: "v", Type: ir.MixedType},
	// 	},
	// 	Result: ir.BoolType,
	// },
	{
		Name: "is_nan",
		Params: []ir.TypeField{
			{Name: "v", Type: ir.FloatType},
		},
		Result: ir.BoolType,
	},
	// {
	// 	Name: "is_array",
	// 	Params: []ir.TypeField{
	// 		{Name: "v", Type: ir.MixedType},
	// 	},
	// 	Result: ir.BoolType,
	// },
	{
		Name: "is_infinite",
		Params: []ir.TypeField{
			{Name: "v", Type: ir.FloatType},
		},
		Result: ir.BoolType,
	},
	{
		Name: "is_bool",
		Params: []ir.TypeField{
			{Name: "v", Type: ir.MixedType},
		},
		Result: ir.BoolType,
	},
	{
		Name: "gettype",
		Params: []ir.TypeField{
			{Name: "v", Type: ir.MixedType},
		},
		Result: ir.StringType,
	},
	{
		Name: "is_finite",
		Params: []ir.TypeField{
			{Name: "v", Type: ir.FloatType},
		},
		Result: ir.BoolType,
	},
	{
		Name: "count",
		Params: []ir.TypeField{
			{Name: "val", Type: &ir.ArrayType{Elem: ir.MixedType}},
		},
		Result: ir.IntType,
	},
	{
		Name: "is_scalar",
		Params: []ir.TypeField{
			{Name: "v", Type: ir.MixedType},
		},
		Result: ir.BoolType,
	},
	{
		Name: "array_sum",
		Params: []ir.TypeField{
			{Name: "a", Type: &ir.ArrayType{Elem: ir.MixedType}},
		},
		Result: ir.FloatType,
	},
	{
		Name: "array_count_values",
		Params: []ir.TypeField{
			{Name: "a", Type: &ir.ArrayType{Elem: ir.MixedType}},
		},
		Result: &ir.ArrayType{Elem: ir.IntType},
	},
	{
		Name: "array_rand",
		Params: []ir.TypeField{
			{Name: "a", Type: &ir.ArrayType{Elem: ir.MixedType}},
			{Name: "num", Type: ir.IntType, Init: 1},
		},
		Result: ir.MixedType,
	},
	{
		Name: "array_flip",
		Params: []ir.TypeField{
			{Name: "a", Type: &ir.ArrayType{Elem: ir.MixedType}},
		},
		Result: &ir.ArrayType{Elem: ir.MixedType},
	},
	{
		Name: "in_array",
		Params: []ir.TypeField{
			{Name: "value", Type: ir.MixedType},
			{Name: "a", Type: &ir.ArrayType{Elem: ir.MixedType}},
			{Name: "strict", Type: ir.BoolType, Init: false},
		},
		Result: ir.BoolType,
	},
	{
		Name: "sizeof",
		Params: []ir.TypeField{
			{Name: "val", Type: &ir.ArrayType{Elem: ir.MixedType}},
		},
		Result: ir.IntType,
	},
	{
		Name: "array_search",
		Params: []ir.TypeField{
			{Name: "val", Type: ir.MixedType},
			{Name: "a", Type: &ir.ArrayType{Elem: ir.MixedType}},
			{Name: "strict", Type: ir.BoolType, Init: false},
		},
		Result: ir.MixedType,
	},
	{
		Name: "array_keys",
		Params: []ir.TypeField{
			{Name: "a", Type: &ir.ArrayType{Elem: ir.MixedType}},
		},
		Result: &ir.ArrayType{Elem: ir.MixedType},
	},
	{
		Name: "explode",
		Params: []ir.TypeField{
			{Name: "delimiter", Type: ir.StringType},
			{Name: "str", Type: ir.StringType},
			{Name: "limit", Type: ir.IntType, Init: 9223372036854775807},
		},
		Result: &ir.ArrayType{Elem: ir.StringType},
	},
	{
		Name: "array_key_exists",
		Params: []ir.TypeField{
			{Name: "v", Type: ir.StringType},
			{Name: "a", Type: &ir.ArrayType{Elem: ir.MixedType}},
		},
		Result: ir.BoolType,
	},
	{
		Name: "implode",
		Params: []ir.TypeField{
			{Name: "s", Type: ir.StringType},
			{Name: "v", Type: &ir.ArrayType{Elem: ir.StringType}},
		},
		Result: ir.StringType,
	},
}
