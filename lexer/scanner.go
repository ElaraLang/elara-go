package lexer

import (
	"bufio"
	"bytes"
	"io"
	"unicode"
)

type Scanner struct {
	r *bufio.Reader
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

//places the previously read rune back on the reader
func (s *Scanner) unread() { _ = s.r.UnreadRune() }

func (s *Scanner) peek() rune {
	peeked := s.read()
	s.unread()
	return peeked
}

func (s *Scanner) Read() (tok TokenType, text string) {
	ch := s.read()

	if ch == eof {
		return EOF, string(ch)
	}

	if isWhitespace(ch) {
		s.consumeWhitespace()
		return s.Read()
	}

	if ch == ',' {
		return Comma, string(ch)
	}
	if ch == ':' {
		return Colon, string(ch)
	}
	if ch == '_' {
		next := s.peek()
		if isWhitespace(next) || next == eof {
			return Underscore, string(ch)
		}
		s.unread()
	}

	if isAngleBracket(ch) {
		s.unread()
		return s.readAngleBracket()
	}

	if isStartOfSymbol(ch) {
		s.unread()
		return s.readSymbol()
	}

	if isOperatorSymbol(ch) {
		s.unread()
		op, txt := s.readOperator()
		if op != Illegal && op != EOF {
			return op, txt
		}
	}

	if isBracket(ch) {
		s.unread()
		return s.readBracket()
	}

	if isNumerical(ch) {
		s.unread()
		return s.readNumber()
	}

	if ch == '"' {
		return s.readString()
	}

	if isValidIdentifier(ch) {
		s.unread()
		return s.readIdentifier()
	}

	return Illegal, string(ch)
}

//Consume all whitespace until we reach an eof or a non-whitespace character
func (s *Scanner) consumeWhitespace() uint {
	count := uint(0)
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		}
		count++
	}
	return count
}

func (s *Scanner) readIdentifier() (tok TokenType, text string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isValidIdentifier(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	str := buf.String()
	switch str {
	case "let":
		return Let, str
	case "mut":
		return Mut, str
	case "extend":
		return Extend, str
	case "return":
		return Return, str
	case "while":
		return While, str
	case "struct":
		return Struct, str
	case "namespace":
		return Namespace, str
	case "import":
		return Import, str
	case "if":
		return If, str
	case "else":
		return Else, str
	case "match":
		return Match, str
	case "as":
		return As, str
	case "is":
		return Is, str
	case "true":
		return Boolean, str
	case "false":
		return Boolean, str
	}

	return Identifier, str
}

func (s *Scanner) readBracket() (tok TokenType, text string) {
	str := s.read()
	switch str {
	case '(':
		return LParen, string(str)
	case ')':
		return RParen, string(str)
	case '{':
		return LBrace, string(str)
	case '}':
		return RBrace, string(str)
	case '<':
		return LAngle, string(str)
	case '>':
		return RAngle, string(str)
	case '[':
		return LSquare, string(str)
	case ']':
		return RSquare, string(str)
	}
	return Illegal, string(str)
}

func (s *Scanner) readSymbol() (tok TokenType, text string) {
	ch := s.read()

	switch ch {
	case '.':
		return Dot, string(ch)
	case '=':
		peeked := s.peek()
		if peeked == '>' {
			s.read()
			return Arrow, string(ch) + string(peeked)
		}
		if peeked == '=' {
			s.read()
			return Equals, string(ch) + string(peeked)
		}
		return Equal, string(ch)
	}

	return Illegal, string(ch)
}
func (s *Scanner) readAngleBracket() (tok TokenType, text string) {
	ch1 := s.read()
	ch := s.peek()
	if ch1 == '<' {
		switch ch {
		case '=':
			s.read()
			return LesserEqual, string(ch1) + string(ch)
		}
		return LAngle, string(ch1)
	}
	if ch1 == '>' {
		switch ch {
		case '=':
			s.read()
			return GreaterEqual, string(ch1) + string(ch)
		}
		return RAngle, string(ch1)
	}

	return Illegal, string(ch1)
}
func (s *Scanner) readOperator() (tok TokenType, text string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isOperatorSymbol(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}
	str := buf.String()
	switch str {
	case "+":
		return Add, str
	case "-":
		return Subtract, str
	case "*":
		return Multiply, str
	case "/":
		return Slash, str
	case "%":
		return Mod, str
	case "&&":
		return And, str
	case "||":
		return Or, str
	case "^":
		return Xor, str
	case "==":
		return Equals, str
	case "!=":
		return NotEquals, str
	case ">=":
		return GreaterEqual, str
	case "<=":
		return LesserEqual, str
	case "!":
		return Not, str

		//Dirty hack, these 2 should probably be in readBracket but oh well...
	case ">":
		return LAngle, str
	case "<":
		return RAngle, str
	}

	return Illegal, str
}

func (s *Scanner) readString() (tok TokenType, text string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			break
		} else if ch == '"' {
			break
		} else {
			buf.WriteRune(ch)
		}
	}
	return String, buf.String()
}

func (s *Scanner) readNumber() (tok TokenType, text string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	numType := Int
	for {
		ch := s.read()
		if ch == eof {
			break
		}
		if ch == '.' {
			numType = Float
		} else if !unicode.IsNumber(ch) {
			break
		}
		buf.WriteRune(ch)
	}

	return numType, buf.String()
}
