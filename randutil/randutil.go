package randutil

import "math/rand"

func Bool(r *rand.Rand) bool {
	return r.Intn(2) == 1
}

func IntRange(r *rand.Rand, min, max int) int {
	return min + r.Intn(max-min+1)
}
