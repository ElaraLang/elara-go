package ast

type Parameter struct {
	Type       Type
	Identifier Identifier
}

func (p *Parameter) ToString() string {
	res := ""
	if p.Type != nil {
		res += p.Type.ToString() + " "
	}
	res += p.Identifier.name
	return res
}

type Entry struct {
	Key   Expression
	Value Expression
}

func (p *Entry) ToString() string {
	return "(" + p.Key.ToString() + " : " + p.Value.ToString() + ")"
}

type StructField struct {
	Type       Type
	Identifier Identifier
}

func (p *StructField) ToString() string {
	return p.Type.ToString() + " " + p.Identifier.name
}
