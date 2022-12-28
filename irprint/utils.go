package irprint

import "github.com/quasilyte/phpsmith/ir"

func accessModifier(flags ir.TypeFlags) string {
	switch {
	case flags.IsPrivate():
		return "private"
	case flags.IsProtected():
		return "protected"
	default:
		return "public"
	}
}
