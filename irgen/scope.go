package irgen

import (
	"github.com/quasilyte/phpsmith/ir"
)

type scope struct {
	vars   []scopeVar
	depths []int
}

type scopeVar struct {
	name string
	typ  ir.Type
}

func newScope() *scope {
	return &scope{}
}

func (s *scope) Enter() {
	s.depths = append(s.depths, 0)
}

func (s *scope) Leave() {
	depth := s.depths[len(s.depths)-1]
	s.depths = s.depths[:len(s.depths)-1]
	s.vars = s.vars[:len(s.vars)-depth]
}

func (s *scope) PushVar(name string, typ ir.Type) {
	s.vars = append(s.vars, scopeVar{name: name, typ: typ})
	s.depths[len(s.depths)-1]++
}

func (s *scope) CurrentBlockVars() []scopeVar {
	depth := s.depths[len(s.depths)-1]
	return s.vars[len(s.vars)-depth:]
}

func (s *scope) FindVarOfType(typ ir.Type) *scopeVar {
	return s.FindVar(func(v *scopeVar) bool {
		return typesIdentical(typ, v.typ)
	})
}

func (s *scope) FindVarByName(name string) *scopeVar {
	return s.FindVar(func(v *scopeVar) bool {
		return v.name == name
	})
}

func (s *scope) FindVar(predicate func(*scopeVar) bool) *scopeVar {
	seen := make(map[string]struct{})
	for i := len(s.vars) - 1; i >= 0; i-- {
		v := &s.vars[i]
		if _, ok := seen[v.name]; ok {
			continue
		}
		seen[v.name] = struct{}{}
		if predicate(v) {
			return v
		}
	}
	return nil
}
