package lexer

import (
	"strconv"
)

type Lexer struct {
	tape *RuneTape
	col  int
	line int
}

func NewLexer(tape *RuneTape) *Lexer {
	return &Lexer{
		tape: tape,
		col:  0,
		line: 0,
	}
}

const eof = rune(-1)

func (l *Lexer) next() rune {
	l.col++
	return l.tape.advance()
}

/*
	Lexes all Data provided by the input channel, sending tokens to the output channel
	This method will block until the channel is marked as empty by sending an eof character
*/
func Lex(input chan rune, output chan Token) {
	tape := NewRuneTape(input)
	lexer := NewLexer(tape)
	for {
		tokenType, data := lexer.readToken()
		token := Token{
			TokenType: tokenType,
			Data:      data,
			Line:      lexer.line,
			Col:       lexer.col,
		}
		output <- token
		if tokenType == EOF {
			break
		}
	}
}

func (l *Lexer) readToken() (TokenType, []rune) {
	char := l.next()
	switch char {
	case eof:
		return EOF, nil
	case '\t':
		l.col += 4
		return l.readToken()
	case ' ': //Skip whitespace
		return l.readToken()
	case '.':
		return Dot, nil
	case '(':
		return LParen, nil
	case ')':
		return RParen, nil
	case '{':
		return LBrace, nil
	case '}':
		return RBrace, nil
	case '<':
		next := l.tape.peek()
		if next == '=' {
			l.next()
			return LesserEqual, nil
		}
		return LAngle, nil
	case '>':
		next := l.tape.peek()
		if next == '=' {
			l.next()
			return GreaterEqual, nil
		}
		return RAngle, nil
	case '[':
		return LSquare, nil
	case ']':
		return RSquare, nil
	case '+':
		return Add, nil
	case '-':
		return Subtract, nil
	case '*':
		return Multiply, nil
	case '/':
		return Slash, nil
	case ':':
		return Colon, nil
	case ',':
		return Comma, nil
	case '!':
		next := l.tape.peek()
		if next == '=' {
			l.next()
			return NotEquals, nil
		}
		return Not, nil
	case '%':
		return Mod, nil
	case '^':
		return Xor, nil
	case '=':
		next := l.tape.peek()
		if next == '=' {
			l.next()
			return Equals, nil
		}
		if next == '>' {
			l.next()
			return Arrow, nil
		}
		return Equal, nil
	case '&':
		next := l.tape.peek()
		if next == '&' {
			l.next()
			return And, nil
		}
		return TypeAnd, nil
	case '|':
		next := l.tape.peek()
		if next == '|' {
			l.next()
			return Or, nil
		}
		return TypeOr, nil

	case '\n':
		l.col = 0
		l.line++
		return NEWLINE, nil
	case '\'':
		symbol := l.readCharLiteral()
		return Char, []rune{symbol}
	case '"':
		str := l.readStringLiteral()
		return String, str
	}
	if isDecimalDigit(char) {
		return l.readNumberLiteral(char)
	}
	return l.readKeywordOrIdentifier(char)
}

// Called after the lexer encounters a ' symbol, attempts to read a char literal
func (l *Lexer) readCharLiteral() rune {
	char := l.readChar()
	closing := l.next()
	if closing != '\'' {
		panic("Expected ' to close character literal, found " + string(closing))
	}
	return char
}

func (l *Lexer) readChar() rune {
	char := l.next()
	if char == '\\' {
		escape := l.next()
		l.col++
		switch escape {
		case 'n':
			char = '\n'
		case 'r':
			char = '\r'
		case '\\':
			char = '\\'
		case 't':
			char = '\t'
		case '\'':
			char = '\''
		default:
			panic("Unknown escape sequence \\" + strconv.QuoteRune(char))
		}
	}
	return char
}

// Called after the lexer encounters a " symbol, attempts to read a string literal
func (l *Lexer) readStringLiteral() []rune {
	runes := make([]rune, 0)

	for {
		c := l.tape.peek()
		if c == eof {
			panic("Unclosed string literal")
		}
		l.next()
		if c == '"' {
			break
		}
		runes = append(runes, c)
	}
	return runes
}

func (l *Lexer) readNumberLiteral(first rune) (numType TokenType, digits []rune) {
	numType = DecimalInt
	digits = make([]rune, 0)
	if first == '0' {
		//Potentially a hex or binary literal
		next := l.tape.peek()
		switch next {
		case 'x':
			l.next()
			return l.readHexInt()
		case 'b':
			l.next()
			return l.readBinaryInt()
		}
	} else {
		digits = append(digits, first)
	}

	for {
		next := l.tape.peek()
		if next == '.' {
			if numType == Float {
				break
			}
			numType = Float
			digits = append(digits, next)
			l.next()
			continue
		}
		if isWhitespace(next) {
			break
		}
		if !isDecimalDigit(next) {
			panic("Illegal symbol in decimal number literal " + string(next))
		}
		digits = append(digits, next)
		l.next()
	}
	return numType, digits
}
func (l *Lexer) readHexInt() (TokenType, []rune) {
	runes := make([]rune, 0)

	for {
		next := l.tape.peek()
		if next == '.' {
			panic("Hexadecimal float literals are not supported")
		}
		if isWhitespace(next) {
			break
		}
		if !isHexDigit(next) {
			panic("Illegal symbol in hexadecimal number literal " + string(next))
		}
		runes = append(runes, next)
		l.next()
	}

	return HexadecimalInt, runes
}
func (l *Lexer) readBinaryInt() (TokenType, []rune) {
	runes := make([]rune, 0)

	for {
		next := l.tape.peek()
		if next == '.' {
			panic("Binary float literals are not supported")
		}
		if isWhitespace(next) {
			break
		}
		if !isBinaryDigit(next) {
			panic("Illegal symbol in binary number literal " + string(next))
		}
		runes = append(runes, next)
		l.next()
	}

	return BinaryInt, runes
}

func (l *Lexer) readKeywordOrIdentifier(prev rune) (TokenType, []rune) {
	buffer := []rune{prev}
	for {
		next := l.tape.peek()
		if isWhitespace(next) || isIllegalIdentifierChar(next) {
			break
		}
		buffer = append(buffer, next)
		l.next()
	}
	if len(buffer) == 0 {
		return 0, nil
	}
	asString := string(buffer)
	switch asString {
	case "let":
		return Let, nil
	case "mut":
		return Mut, nil
	case "extend":
		return Extend, nil
	case "return":
		return Return, nil
	case "while":
		return While, nil
	case "lazy":
		return Lazy, nil
	case "struct":
		return Struct, nil
	case "namespace":
		return Namespace, nil
	case "import":
		return Import, nil
	case "type":
		return Type, nil
	case "if":
		return If, nil
	case "else":
		return Else, nil
	case "match":
		return Match, nil
	case "as":
		return As, nil
	case "is":
		return Is, nil
	case "open":
		return Open, nil
	case "try":
		return Try, nil
	case "true":
		return BooleanTrue, nil
	case "false":
		return BooleanFalse, nil
	default:
		return Identifier, buffer
	}
}
