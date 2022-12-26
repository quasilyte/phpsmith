package irgen

import (
	"math"
	"math/rand"
	"strings"
	"unicode"

	"github.com/quasilyte/phpsmith/randutil"
)

type valueGenerator struct {
	rand *rand.Rand
}

func newValueGenerator(r *rand.Rand) *valueGenerator {
	return &valueGenerator{rand: r}
}

func toEfaceSlice[T any](xs []T) []any {
	result := make([]any, len(xs))
	for i, x := range xs {
		result[i] = x
	}
	return result
}

func generateUniqueValues[T comparable](n int, f func() T) []T {
	set := make(map[T]struct{}, n)
	for len(set) < n {
		x := f()
		if _, ok := set[x]; ok {
			continue
		}
		set[x] = struct{}{}
	}
	slice := make([]T, 0, len(set))
	for x := range set {
		slice = append(slice, x)
	}
	return slice
}

func (g *valueGenerator) IntValue() int64 {
	switch g.rand.Intn(8) {
	case 0, 1:
		return int64(g.rand.Intn(0xffff))
	case 2, 3:
		return -int64(g.rand.Intn(0xffff))
	case 4:
		return int64(randutil.IntRange(g.rand, 100000, 19438420511))
	default:
		return intLitValues[g.rand.Intn(len(intLitValues))]
	}
}

func (g *valueGenerator) FloatValue() float64 {
	switch g.rand.Intn(8) {
	case 0:
		return g.rand.Float64()
	case 2, 3:
		return g.rand.Float64() * float64(g.rand.Intn(1000))
	case 4:
		return g.rand.Float64() * float64(g.rand.Intn(10000000))
	default:
		return floatLitValues[g.rand.Intn(len(floatLitValues))]
	}
}

func (g *valueGenerator) StringValue() string {
	if randutil.Chance(g.rand, 0.2) {
		return randutil.Elem(g.rand, stringLitValues)
	}

	var s strings.Builder
	count := randutil.IntRange(g.rand, 1, 6)
	for i := 0; i < count; i++ {
		ch := g.rand.Intn(unicode.MaxASCII)
		if !unicode.IsPrint(rune(ch)) || ch == '$' {
			s.WriteString(stringLitValues[g.rand.Intn(len(stringLitValues))])
		} else {
			s.WriteByte(byte(ch))
		}
	}
	return s.String()
}

func (b *valueGenerator) BoolValue() bool {
	return randutil.Chance(b.rand, 0.5)
}

var intLitValues = []int64{
	0,
	-1,
	0xff,
	9284128,
	128412288,
	-9284120,
	-0xff,
}

var floatLitValues = []float64{
	0,
	-1,
	2.51,
	329.5,
	0.00043,
	21948.293242,
	-2222.9999,
	2842.6378,
	math.NaN(),
	math.Inf(1),
	math.Inf(-1),
}

var stringLitValues = []string{
	"",
	",",
	" ",
	"``",
	"''",
	"0x1f",
	"000",
	"24",
	"-123",
	"\x00",
	"simple string",
	"ハロー・ワールド",
	"1\n2",
	"<div/>",
	"<h1>ok</h1>",
	"<p>",
	"</p>",
	`{"key":1}`,
	`["val"]`,
}
