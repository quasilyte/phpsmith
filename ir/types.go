package ir

var (
	VoidType   = &ScalarType{Kind: ScalarVoid}
	BoolType   = &ScalarType{Kind: ScalarBool}
	IntType    = &ScalarType{Kind: ScalarInt}
	FloatType  = &ScalarType{Kind: ScalarFloat}
	StringType = &ScalarType{Kind: ScalarString}
	MixedType  = &ScalarType{Kind: ScalarMixed}
)

type TypeField struct {
	Name   string
	Type   Type
	Strict bool
	Init   interface{}
}

type ScalarKind int

const (
	ScalarUnknown ScalarKind = iota
	ScalarVoid
	ScalarBool
	ScalarInt
	ScalarFloat
	ScalarString
	ScalarMixed
)

func (k ScalarKind) String() string {
	switch k {
	case ScalarVoid:
		return "void"
	case ScalarBool:
		return "bool"
	case ScalarInt:
		return "int"
	case ScalarFloat:
		return "float"
	case ScalarString:
		return "string"
	default:
		return "?"
	}
}

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

type FuncType struct {
	Name       string
	Params     []TypeField
	MinArgsNum int
	Result     Type
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
