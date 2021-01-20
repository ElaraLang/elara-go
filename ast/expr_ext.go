package ast

import (
	"github.com/ElaraLang/elara/util"
	"strconv"
)

func (e *BinaryExpression) expressionNode() {}
func (e *BinaryExpression) TokenValue() string {
	return e.Token.String()
}
func (e *BinaryExpression) ToString() string {
	return e.Left.ToString() + " " + e.Operator.TokenType.String() + " " + e.Right.ToString()
}

func (e *UnaryExpression) expressionNode() {}
func (e *UnaryExpression) TokenValue() string {
	return e.Token.String()
}
func (e *UnaryExpression) ToString() string {
	return e.Operator.TokenType.String() + " " + e.Right.ToString()
}

func (e *PropertyExpression) expressionNode() {}
func (e *PropertyExpression) TokenValue() string {
	return e.Token.String()
}
func (e *PropertyExpression) ToString() string {
	return "(" + e.Context.ToString() + ")." + e.Variable.name
}

func (e *IfExpression) expressionNode() {}
func (e *IfExpression) TokenValue() string {
	return e.Token.String()
}
func (e *IfExpression) ToString() string {
	return string(e.Token.Text) + " " +
		e.Condition.ToString() + " {\n" +
		e.MainBranch.ToString() + "\n} else {\n" +
		e.ElseBranch.ToString() + "\n}"
}

func (e *AccessExpression) expressionNode() {}
func (e *AccessExpression) TokenValue() string {
	return e.Token.String()
}
func (e *AccessExpression) ToString() string {
	return "(" + e.Expression.ToString() + ")[" + e.Index.ToString() + "]"
}

func (e *CallExpression) expressionNode() {}
func (e *CallExpression) TokenValue() string {
	return e.Token.String()
}
func (e *CallExpression) ToString() string {
	return "(" + e.Expression.ToString() + ")(" + util.JoinToString(e.Arguments, ", ") + ")"
}

func (e *TypeCastExpression) expressionNode() {}
func (e *TypeCastExpression) TokenValue() string {
	return e.Token.String()
}
func (e *TypeCastExpression) ToString() string {
	return "(" + e.Expression.ToString() + ") as (" + e.Type.ToString() + ")"
}

func (e *TypeCheckExpression) expressionNode() {}
func (e *TypeCheckExpression) TokenValue() string {
	return e.Token.String()
}
func (e *TypeCheckExpression) ToString() string {
	return "(" + e.Expression.ToString() + ") is (" + e.Type.ToString() + ")"
}

func (e *FunctionLiteral) expressionNode() {}
func (e *FunctionLiteral) TokenValue() string {
	return e.Token.String()
}
func (e *FunctionLiteral) ToString() string {
	return "(" + util.JoinToString(e.Parameters, ", ") + ") => " + e.ReturnType.ToString() + " {\n" + e.Body.ToString() + "\n}\n"
}

func (e *MapLiteral) expressionNode() {}
func (e *MapLiteral) TokenValue() string {
	return e.Token.String()
}
func (e *MapLiteral) ToString() string {
	return "{\n" + util.JoinToString(e.Entries, ",\n") + "\n}\n"
}

func (e *CollectionLiteral) expressionNode() {}
func (e *CollectionLiteral) TokenValue() string {
	return e.Token.String()
}
func (e *CollectionLiteral) ToString() string {
	return "[" + util.JoinToString(e.Elements, ", ") + "]"
}

func (e *BooleanLiteral) expressionNode() {}
func (e *BooleanLiteral) TokenValue() string {
	return e.Token.String()
}
func (e *BooleanLiteral) ToString() string {
	return strconv.FormatBool(e.Value)
}

func (e *IntegerLiteral) expressionNode() {}
func (e *IntegerLiteral) TokenValue() string {
	return e.Token.String()
}
func (e *IntegerLiteral) ToString() string {
	return strconv.FormatInt(e.Value, 10)
}

func (e *FloatLiteral) expressionNode() {}
func (e *FloatLiteral) TokenValue() string {
	return e.Token.String()
}
func (e *FloatLiteral) ToString() string {
	return strconv.FormatFloat(e.Value, 10, 4, 64)
}

func (e *CharLiteral) expressionNode() {}
func (e *CharLiteral) TokenValue() string {
	return e.Token.String()
}
func (e *CharLiteral) ToString() string {
	return string(e.Value)
}

func (e *StringLiteral) expressionNode() {}
func (e *StringLiteral) TokenValue() string {
	return e.Token.String()
}
func (e *StringLiteral) ToString() string {
	return e.Value
}
