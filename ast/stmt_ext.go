package ast

import "github.com/ElaraLang/elara/util"

func (s *ImportStatement) statementNode() {}
func (s *ImportStatement) TokenValue() string {
	return s.Token.String()
}
func (s *ImportStatement) ToString() string {
	return s.TokenValue() + " " + s.Module.ToString()
}

func (s *NamespaceStatement) statementNode() {}
func (s *NamespaceStatement) TokenValue() string {
	return s.Token.String()
}
func (s *NamespaceStatement) ToString() string {
	return s.TokenValue() + " " + s.Module.ToString()
}

func (s *ExpressionStatement) statementNode() {}
func (s *ExpressionStatement) TokenValue() string {
	return s.Token.String()
}
func (s *ExpressionStatement) ToString() string {
	return s.ToString()
}

func (s *DeclarationStatement) statementNode() {}
func (s *DeclarationStatement) TokenValue() string {
	return s.Token.String()
}
func (s *DeclarationStatement) ToString() string {
	var typ string
	if s.Type == nil {
		typ = "INFER"
	} else {
		typ = s.Type.ToString()
	}
	return "let " + util.JoinStringConditionally(map[string]bool{
		"mut":  s.Mutable,
		"lazy": s.Lazy,
		"open": s.Open,
	}, " ") + s.Identifier.Name + ":" + typ + " = " + s.Value.ToString()
}

func (s *StructDefStatement) statementNode() {}
func (s *StructDefStatement) TokenValue() string {
	return s.Token.String()
}
func (s *StructDefStatement) ToString() string {
	return s.TokenValue() + " " + s.Id.Name +
		" {\n" + joinToString(s.Fields, " ") + "\n}\n"
}

func (s *WhileStatement) statementNode() {}
func (s *WhileStatement) TokenValue() string {
	return s.Token.String()
}
func (s *WhileStatement) ToString() string {
	return s.TokenValue() + " " + s.Condition.ToString() + " " + s.Body.ToString()
}

func (s *ExtendStatement) statementNode() {}
func (s *ExtendStatement) TokenValue() string {
	return s.Token.String()
}
func (s *ExtendStatement) ToString() string {
	return s.TokenValue() + " " +
		s.Identifier.Name + " as " + s.Alias.Name + " " +
		s.Body.ToString()
}

func (s *BlockStatement) statementNode() {}
func (s *BlockStatement) TokenValue() string {
	return s.Token.String()
}
func (s *BlockStatement) ToString() string {
	result := "{\n"
	for _, v := range s.Block {
		result += v.ToString()
		result += "\n"
	}
	return result + "}\n"
}

func (s *TypeStatement) statementNode() {}
func (s *TypeStatement) TokenValue() string {
	return s.Token.String()
}
func (s *TypeStatement) ToString() string {
	return s.TokenValue() + " " + s.Identifier.Name + " = " + s.Contract.ToString()
}

func (s *GenerifiedStatement) statementNode() {}
func (s *GenerifiedStatement) TokenValue() string {
	return s.Token.String()
}
func (s *GenerifiedStatement) ToString() string {
	return "<" + joinToString(s.Contracts, " ") + ">\n" + s.Statement.ToString()
}

func (s *ReturnStatement) statementNode() {}
func (s *ReturnStatement) TokenValue() string {
	return s.Token.String()
}
func (s *ReturnStatement) ToString() string {
	return s.TokenValue() + " " + s.Value.ToString()
}
