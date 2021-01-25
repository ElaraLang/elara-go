package ast

func (p *PrimaryType) typeNode() {}
func (p *PrimaryType) TokenValue() string {
	return p.Token.String()
}
func (p *PrimaryType) ToString() string {
	return p.Identifier.Name
}

func (p *MapType) typeNode() {}
func (p *MapType) TokenValue() string {
	return p.Token.String()
}
func (p *MapType) ToString() string {
	return "{" + p.KeyType.ToString() + ":" + p.ValueType.ToString() + "}"
}

func (p *CollectionType) typeNode() {}
func (p *CollectionType) TokenValue() string {
	return p.Token.String()
}
func (p *CollectionType) ToString() string {
	return "[" + p.Type.ToString() + "]"
}

func (p *FunctionType) typeNode() {}
func (p *FunctionType) TokenValue() string {
	return p.Token.String()
}
func (p *FunctionType) ToString() string {
	return "(" + joinToString(p.ParamTypes, ", ") + ") =>" + p.ReturnType.ToString()
}

func (p *AlgebraicType) typeNode() {}
func (p *AlgebraicType) TokenValue() string {
	return p.Token.String()
}
func (p *AlgebraicType) ToString() string {
	return "(" + p.Left.ToString() + " " + string(p.Operation.Text) + " " + p.Right.ToString() + ")"
}

func (p *ContractualType) typeNode() {}
func (p *ContractualType) TokenValue() string {
	return p.Token.String()
}
func (p *ContractualType) ToString() string {
	return "type { " + joinToString(p.Contracts, ", ") + " }"
}
