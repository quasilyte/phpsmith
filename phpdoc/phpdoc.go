package phpdoc

type Tag interface {
	Name() string
	Value() string
}

type ReturnTag struct {
	Type string
}

type VarTag struct {
	Type    string
	VarName string
}

type ParamTag struct {
	Type    string
	VarName string
}

func (t *ReturnTag) Name() string { return "return" }

func (t *ReturnTag) Value() string { return t.Type }

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
