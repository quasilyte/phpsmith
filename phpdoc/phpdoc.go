package phpdoc

type Tag interface {
	Name() string
	Value() string
}

type ParamTag struct {
	Type    string
	VarName string
}

func (t *ParamTag) Name() string { return "param" }

func (t *ParamTag) Value() string {
	return t.Type + " " + t.VarName
}
