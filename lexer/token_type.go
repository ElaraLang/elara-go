package lexer

type TokenType int

var EOF_TOKEN = Token{
	TokenType: EOF,
	Data:      nil,
	Line:      -1,
	Col:       -1,
}

const (
	//special tokens
	Illegal TokenType = iota
	EOF
	NEWLINE // \n

	//Brackets
	LParen  // (
	RParen  // )
	LBrace  // {
	RBrace  // }
	LAngle  //<
	RAngle  //>
	LSquare // [
	RSquare // ]

	//Keywords
	Let
	Extend
	Return
	While
	Where
	Mut
	Lazy
	Open
	Struct
	Namespace
	Import
	Type
	If
	Else
	Match
	As
	Is
	Try

	//Operators
	Add          // +
	Subtract     // -
	Multiply     // *
	Slash        // /
	Mod          // %
	And          // &&
	Or           // ||
	Xor          // ^
	Equals       // ==
	NotEquals    // !=
	GreaterEqual // >=
	LesserEqual  // <=
	Not          // !

	TypeOr  // |
	TypeAnd // &

	//Symbol
	Equal // =
	Arrow // =>
	Dot   // .
	Hash

	//Literals
	BooleanTrue
	BooleanFalse
	String
	Char
	DecimalInt
	BinaryInt
	OctalInt
	HexadecimalInt
	Float

	Comma // ,
	Colon // :

	Identifier
	//Underscore
)

func (token TokenType) String() string {
	return tokenNames[token]
}

var tokenNames = map[TokenType]string{
	Illegal:        "Illegal",
	EOF:            "EOF",
	NEWLINE:        "NEWLINE",
	LParen:         "LParen",
	RParen:         "RParen",
	LBrace:         "LBrace",
	RBrace:         "RBrace",
	LAngle:         "LAngle",
	RAngle:         "RAngle",
	LSquare:        "LSquare",
	RSquare:        "RSquare",
	Type:           "TokenType",
	Let:            "Let",
	Extend:         "Extend",
	Return:         "Return",
	While:          "While",
	Where:          "Where",
	Mut:            "Mut",
	Lazy:           "Lazy",
	Open:           "Open",
	Struct:         "Struct",
	Namespace:      "Namespace",
	Import:         "Import",
	If:             "If",
	Else:           "Else",
	Match:          "Match",
	As:             "As",
	Is:             "Is",
	Try:            "Try",
	Add:            "Add",
	Subtract:       "Subtract",
	Multiply:       "Multiply",
	Slash:          "Slash",
	Mod:            "Mod",
	And:            "And",
	Or:             "Or",
	Xor:            "Xor",
	Equals:         "Equals",
	NotEquals:      "NotEquals",
	GreaterEqual:   "GreaterEqual",
	LesserEqual:    "LesserEqual",
	Not:            "Not",
	Equal:          "Equal",
	Arrow:          "Arrow",
	Dot:            "Dot",
	Hash:           "Hash",
	BooleanTrue:    "True",
	BooleanFalse:   "False",
	String:         "String",
	Char:           "Char",
	DecimalInt:     "DecimalInt",
	BinaryInt:      "BinaryInt",
	HexadecimalInt: "HexadecimalInt",
	OctalInt:       "OctalInt",
	Float:          "Float",

	Comma: "Comma",
	Colon: "Colon",

	Identifier: "Identifier",
	//Underscore: "Underscore",
}
