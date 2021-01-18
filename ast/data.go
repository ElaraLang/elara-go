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
