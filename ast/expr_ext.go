package ast

func (b *BinaryExpression) expressionNode() {}
func (b *BinaryExpression) TokenValue() string {
	return b.Token.String()
}
func (b *BinaryExpression) toString() string {
	return b.Left.toString() + " " + b.Operator.TokenType.String() + " " + b.Right.toString()
}
