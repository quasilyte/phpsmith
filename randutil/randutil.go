package randutil

import "math/rand"

func Bool(r *rand.Rand) bool {
	return r.Intn(2) == 1
}
