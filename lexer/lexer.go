package lexer

import "strconv"

type Lexer struct {
	col  int
	line int
}

const eof = rune(-1)

func (l *Lexer) next(input chan rune) rune {
	char, ok := <-input
	if !ok {
		return eof
	}
	l.col++
	return char
}

/*
	Lexes all data provided by the input channel, sending tokens to the output channel
	This method will block until the channel is marked as empty
*/
func Lex(input chan rune, output chan Token) {
	lexer := &Lexer{}
	for {
		tokenType, data := lexer.readToken(input)
		token := Token{
			TokenType: tokenType,
			data:      data,
			line:      lexer.line,
			col:       lexer.col,
		}
		output <- token
		if tokenType == EOF {
			break
		}
	}
}

func (l *Lexer) readToken(input chan rune) (TokenType, []rune) {
	char := l.next(input)
	switch char {
	case eof:
		return EOF, nil
	case ' ': //Skip whitespace
		return l.readToken(input)
	case '(':
		return LParen, nil
	case ')':
		return RParen, nil
	case '{':
		return LBrace, nil
	case '}':
		return RBrace, nil
	case '<':
		return LAngle, nil
	case '>':
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
		return Not, nil
	case '%':
		return Mod, nil
	case '^':
		return Xor, nil
	case '\n':
		l.col = 0
		l.line++
		return NEWLINE, nil

	case '\'':
		symbol := l.readCharLiteral(input)
		return Char, []rune{symbol}
	case '"':
		str := l.readStringLiteral(input)
		return String, str
	}
	return 0, nil
}

// Called after the lexer encounters a ' symbol, attempts to read a char literal
func (l *Lexer) readCharLiteral(input chan rune) rune {
	char := l.readChar(input)
	closing := l.next(input)
	if closing != '\'' {
		panic("Expected ' to close character literal, found " + string(closing))
	}
	return char
}
func (l *Lexer) readChar(input chan rune) rune {
	char := l.next(input)
	if char == '\\' {
		escape := l.next(input)
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
func (l *Lexer) readStringLiteral(input chan rune) []rune {
	runes := make([]rune, 0)
	c := l.next(input)
	for c != '"' {
		runes = append(runes, c)
		if c == eof {
			panic("Unclosed string literal")
		}
	}
	return runes
}
