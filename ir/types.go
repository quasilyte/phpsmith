package ir

import (
	"strings"

	"github.com/quasilyte/phpsmith/phpdoc"
)

type Type interface {
	String() string
	Tag() int
}

const (
	TypeTagScalar = iota
	TypeTagClass
	TypeTagUnion
	TypeTagNullable
	TypeTagArray
	TypeTagTuple
	TypeTagFunc
	TypeTagEnum
)

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
	Flags  TypeFlags
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
	case ScalarMixed:
		return "mixed"
	default:
		return "?"
	}
}

type ScalarType struct {
	Kind ScalarKind
}

type ClassType struct {
	Name string

	Fields []TypeField

	Methods []*FuncType
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

type TupleType struct {
	Elems []Type
}

type FuncType struct {
	Name       string
	Params     []TypeField
	Tags       []phpdoc.Tag
	MinArgsNum int
	Result     Type
	NeedCast   bool
	IsLibFunc  bool
	Class      *ClassType
}

func (typ *FuncType) FullName() string {
	if typ.Class == nil {
		return typ.Name
	}
	return typ.Class.Name + "::" + typ.Name
}

type EnumType struct {
	ValueType *ScalarType
	Values    []interface{}
}

func (typ *ScalarType) Tag() int   { return TypeTagScalar }
func (typ *ClassType) Tag() int    { return TypeTagClass }
func (typ *UnionType) Tag() int    { return TypeTagUnion }
func (typ *NullableType) Tag() int { return TypeTagNullable }
func (typ *ArrayType) Tag() int    { return TypeTagArray }
func (typ *EnumType) Tag() int     { return TypeTagEnum }
func (typ *TupleType) Tag() int    { return TypeTagTuple }
func (typ *FuncType) Tag() int     { return TypeTagFunc }

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

func (typ *EnumType) String() string {
	return typ.ValueType.String()
}

func (typ *TupleType) String() string {
	parts := make([]string, len(typ.Elems))
	for i, e := range typ.Elems {
		parts[i] = e.String()
	}
	return "tuple(" + strings.Join(parts, ",") + ")"
}

func (typ *FuncType) String() string {
	args := make([]string, len(typ.Params))
	for i, e := range typ.Params {
		args[i] = e.Type.String()
	}
	return "callable(" + strings.Join(args, ",") + "):" + typ.Result.String()
}
