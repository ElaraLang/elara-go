package lexer

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"unicode"
)

type Scanner struct {
	r    *bufio.Reader
	line int
	col  int
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
func (s *Scanner) unread() {
	_ = s.r.UnreadRune()
}

func (s *Scanner) peek() rune {
	peeked := s.read()
	s.unread()
	return peeked
}

func (s *Scanner) Read() (tok TokenType, text string, line int, col int) {
	ch := s.read()

	if ch == eof {
		return EOF, string(ch), s.line, s.col
	}

	if ch == '\r' {
		return s.Read()
	}

	if ch == '\n' {
		defer func() {
			s.col = 0
			s.line++
		}()
		return NEWLINE, string(ch), s.line, s.col
	}

	if isWhitespace(ch) {
		s.col++ //consumeWhitespace will add any *further* whitespace counters
		s.col += s.consumeWhitespace()
		return s.Read()
	}

	if ch == ',' {
		defer func() {
			s.col++
		}()
		return Comma, string(ch), s.line, s.col
	}
	if ch == ':' {
		defer func() {
			s.col++
		}()
		return Colon, string(ch), s.line, s.col
	}

	if ch == '_' {
		next := s.peek()
		if isWhitespace(next) || next == eof {
			defer func() {
				s.col++
			}()
			return Underscore, string(ch), s.line, s.col
		}
		s.unread()
	}

	if isAngleBracket(ch) {
		s.unread()
		bracket, t := s.readAngleBracket()
		defer func() {
			s.col += len(t)
		}()
		return bracket, t, s.line, s.col
	}

	if isStartOfSymbol(ch) {
		s.unread()
		symbol, t := s.readSymbol()
		defer func() {
			s.col += len(t)
		}()
		return symbol, t, s.line, s.col
	}

	if isOperatorSymbol(ch) {
		s.unread()
		op, t := s.readOperator()
		defer func() {
			s.col += len(t)
		}()
		if op != Illegal && op != EOF {
			return op, t, s.line, s.col
		}
	}

	if isBracket(ch) {
		s.unread()
		bracket, t := s.readBracket()
		defer func() {
			s.col += len(t)
		}()
		return bracket, t, s.line, s.col
	}

	if isNumerical(ch) {
		s.unread()
		number, t := s.readNumber()
		defer func() {
			s.col += len(t)
		}()
		return number, t, s.line, s.col
	}

	if ch == '"' {
		str, t := s.readString()
		defer func() {
			s.col += len(t)
		}()
		return str, t, s.line, s.col
	}

	if isValidIdentifier(ch) {
		s.unread()
		identifier, t := s.readIdentifier()
		defer func() {
			s.col += len(t)
		}()
		return identifier, t, s.line, s.col
	}

	return Illegal, string(ch), s.line, s.col
}

//Consume all whitespace until we reach an eof or a non-whitespace character
func (s *Scanner) consumeWhitespace() int {
	count := 0
	for {
		ch := s.read()
		if ch == eof {
			break
		}
		if !isWhitespace(ch) {
			s.unread()
			break
		}
		count++
	}
	return count
}

func (s *Scanner) readIdentifier() (tok TokenType, text string) {
	var builder strings.Builder
	builder.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isValidIdentifier(ch) {
			s.unread()
			break
		} else {
			builder.WriteRune(ch)
		}
	}

	str := builder.String()
	switch str {
	case "let":
		return Let, str
	case "type":
		return Type, str
	case "mut":
		return Mut, str
	case "restricted":
		return Restricted, str
	case "lazy":
		return Lazy, str
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
		if ch == '\n' {
			s.unread()
			break
		}
		if ch == '.' {
			numType = Float
		} else if !unicode.IsNumber(ch) {
			s.unread()
			break
		}

		buf.WriteRune(ch)
	}

	return numType, buf.String()
}
