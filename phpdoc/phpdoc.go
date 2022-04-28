package phpdoc

type Tag interface {
	Name() string
	Value() string
}

type VarTag struct {
	Type    string
	VarName string
}

type ParamTag struct {
	Type    string
	VarName string
}

func (t *VarTag) Name() string { return "var" }

func (t *VarTag) Value() string {
	if t.VarName == "" {
		return t.Type
	}
	return t.Type + " " + t.VarName
}

func (t *ParamTag) Name() string { return "param" }

func (t *ParamTag) Value() string {
	return t.Type + " " + t.VarName
}
