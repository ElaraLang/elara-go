package lexer

type TokenType int

const (
	//special tokens
	Illegal TokenType = iota
	EOF

	//Brackets
	LParen
	RParen
	LBrace
	RBrace
	LAngle
	RAngle
	LSquare
	RSquare

	//Keywords
	Let
	Extend
	Return
	While
	Mut
	Struct
	Namespace
	Import
	If
	Else
	Match

	//Operators
	Add
	Subtract
	Multiply
	Slash
	Mod
	And
	Or
	Xor
	Equals
	NotEquals
	GreaterEqual
	Greater
	Lesser
	LesserEqual
	Not

	//Symbol
	Equal
	Arrow
	Dot

	//Literals
	Boolean
	String
	Int
	Float

	Comma
	Colon

	Identifier
	Underscore
)

func (token *TokenType) string() string {
	return TokenNames[*token]
}

var TokenNames = map[TokenType]string{
	Illegal: "Illegal",
	EOF:     "EOF",

	LParen:       "LParen",
	RParen:       "RParen",
	LBrace:       "LBrace",
	RBrace:       "RBrace",
	LAngle:       "LAngle",
	RAngle:       "RAngle",
	LSquare:      "LSquare",
	RSquare:      "RSquare",
	Let:          "Let",
	Extend:       "Extend",
	Return:       "Return",
	While:        "While",
	Mut:          "Mut",
	Struct:       "Struct",
	Namespace:    "Namespace",
	Import:       "Import",
	If:           "If",
	Else:         "Else",
	Match:        "Match",
	Add:          "Add",
	Subtract:     "Subtract",
	Multiply:     "Multiply",
	Slash:        "Slash",
	Mod:          "Mod",
	And:          "And",
	Or:           "Or",
	Xor:          "Xor",
	Equals:       "Equals",
	NotEquals:    "NotEquals",
	GreaterEqual: "GreaterEqual",
	Greater:      "Greater",
	Lesser:       "Lesser",
	LesserEqual:  "LesserEqual",
	Not:          "Not",
	Equal:        "Equal",
	Arrow:        "Arrow",
	Dot:          "Dot",
	Boolean:      "Boolean",
	String:       "String",
	Int:          "Int",
	Float:        "Float",

	Comma: "Comma",
	Colon: "Colon",

	Identifier: "Identifier",
	Underscore: "Underscore",
}
var IllegalIdentifierChars = map[rune]bool{
	',': true,
	'.': true,
	':': true,
	'#': true,
	'[': true,
	']': true,
	'(': true,
	')': true,
	'{': true,
	'}': true,
	'"': true,
}
