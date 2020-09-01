package lexer

type TokenType int

const (
	//special tokens
	ILLEGAL TokenType = iota
	EOF

	LPAREN
	RPAREN
	LBRACE
	RBRACE
	LANGLE
	RANGLE

	//Keywords
	LET
	EXTEND
	RETURN
	WHILE
	MUT
	STRUCT
	NAMESPACE
	IMPORT
	IF
	ELSE
	MATCH

	//Operators
	ADD
	SUBTRACT
	MULTIPLY
	DIVIDE
	MOD
	AND
	OR
	XOR
	EQUALS
	NOT_EQUALS
	GREATER_EQUAL
	GREATER
	LESSER
	LESSER_EQUAL
	NOT
	EQUAL
	ARROW

	DOT

	//Literals
	BOOLEAN
	STRING
	INT
	FLOAT

	COMMA
	COLON
	SLASH

	IDENTIFIER
	UNDERSCORE
)

func (token *TokenType) string() string {
	return TokenNames[*token]
}

var TokenNames = map[TokenType]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",

	LPAREN:        "LPAREN",
	RPAREN:        "RPAREN",
	LBRACE:        "LBRACE",
	RBRACE:        "RBRACE",
	LANGLE:        "LANGLE",
	RANGLE:        "RANGLE",
	LET:           "LET",
	EXTEND:        "EXTEND",
	RETURN:        "RETURN",
	WHILE:         "WHILE",
	MUT:           "MUT",
	STRUCT:        "STRUCT",
	NAMESPACE:     "NAMESPACE",
	IMPORT:        "IMPORT",
	IF:            "IF",
	ELSE:          "ELSE",
	MATCH:         "MATCH",
	ADD:           "ADD",
	SUBTRACT:      "SUBTRACT",
	MULTIPLY:      "MULTIPLY",
	DIVIDE:        "DIVIDE",
	MOD:           "MOD",
	AND:           "AND",
	OR:            "OR",
	XOR:           "XOR",
	EQUALS:        "EQUALS",
	NOT_EQUALS:    "NOT_EQUALS",
	GREATER_EQUAL: "GREATER_EQUAL",
	GREATER:       "GREATER",
	LESSER:        "LESSER",
	LESSER_EQUAL:  "LESSER_EQUAL",
	NOT:           "NOT",
	EQUAL:         "EQUAL",
	ARROW:         "ARROW",
	DOT:           "DOT",
	BOOLEAN:       "BOOLEAN",
	STRING:        "STRING",
	INT:           "INT",
	FLOAT:         "FLOAT",

	COMMA: "COMMA",
	COLON: "COLON",
	SLASH: "SLASH",

	IDENTIFIER: "IDENTIFIER",
	UNDERSCORE: "UNDERSCORE",
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
