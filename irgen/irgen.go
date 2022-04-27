package irgen

import (
	"math/rand"

	"github.com/quasilyte/phpsmith/ir"
)

type Config struct {
	Rand *rand.Rand
}

type Program struct {
	Files []*File
}

type File struct {
	Name string

	Nodes []ir.RootNode
}

func CreateProgram(config *Config) *Program {
	g := newGenerator(config)
	return g.CreateProgram()
}
