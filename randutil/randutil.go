package randutil

import "math/rand"

func Bool(r *rand.Rand) bool {
	return r.Intn(2) == 1
}

func Chance(r *rand.Rand, probability float64) bool {
	return r.Float64() <= probability
}

func IntRange(r *rand.Rand, min, max int) int {
	return min + r.Intn(max-min+1)
}

func Elem[T any](r *rand.Rand, xs []T) T {
	index := r.Intn(len(xs))
	return xs[index]
}
