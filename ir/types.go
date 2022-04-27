package ir

//go:generate stringer -type ScalarKind -trimprefix Scalar
type ScalarKind int

const (
	ScalarUnknown ScalarKind = iota
	ScalarVoid
	ScalarBool
	ScalarInt
	ScalarFloat
	ScalarString
)

type Type interface {
	String() string
}

type ScalarType struct {
	Kind ScalarKind
}

type ClassType struct {
	Name string
}

type UnionType struct {
	X Type
	Y Type
}

type NullableType struct {
	X Type
}

type ArrayType struct {
	Elem Type
}

func (typ *ScalarType) String() string {
	return typ.Kind.String()
}

func (typ *ClassType) String() string {
	return typ.Name
}

func (typ *UnionType) String() string {
	return "(" + typ.X.String() + "|" + typ.Y.String() + ")"
}

func (typ *NullableType) String() string {
	return "(?" + typ.X.String() + ")"
}

func (typ *ArrayType) String() string {
	return "(" + typ.Elem.String() + "[])"
}
