package lexer

import (
	"bytes"
	"unicode"
)

type TokenReader struct {
	runes  []rune
	cursor int
	line   int
	col    int
}

func NewTokenReader(runes []rune) *TokenReader {
	return &TokenReader{
		runes:  runes,
		cursor: 0,
		line:   0,
		col:    0,
	}
}

//Reads the current rune and moves the cursor to the next rune
func (s *TokenReader) read() rune {
	if s.cursor >= len(s.runes) {
		return eof
	}
	r := s.runes[s.cursor]
	s.cursor++
	return r
}

//Goes back to reading the previous rune
func (s *TokenReader) unread() {
	s.cursor--
}

func (s *TokenReader) peek() rune {
	if s.cursor >= len(s.runes)-1 {
		return eof
	}
	return s.runes[s.cursor]
}

//TODO this is pretty gross, could use a cleanup
func (s *TokenReader) Read() (tok TokenType, text []rune, line int, col int) {
	ch := s.read()

	if ch == eof {
		return EOF, []rune{ch}, s.line, s.col
	}

	if ch == '\r' {
		return s.Read()
	}

	if ch == '\n' {
		defer func() {
			s.col = 0
			s.line++
		}()
		return NEWLINE, []rune{ch}, s.line, s.col
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
		return Comma, []rune{ch}, s.line, s.col
	}
	if ch == ':' {
		defer func() {
			s.col++
		}()
		return Colon, []rune{ch}, s.line, s.col
	}

	if ch == '_' {
		next := s.peek()
		if isWhitespace(next) || next == eof {
			defer func() {
				s.col++
			}()
			return Underscore, []rune{ch}, s.line, s.col
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

	return Illegal, []rune{ch}, s.line, s.col
}

//Consume all whitespace until we reach an eof or a non-whitespace character
func (s *TokenReader) consumeWhitespace() int {
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

func (s *TokenReader) readIdentifier() (tok TokenType, text []rune) {
	i := s.cursor
	end := i + 1
	for {
		r := s.runes[end]
		if r == eof || !isValidIdentifier(r) {
			break
		}
		end++
		if end >= len(s.runes) {
			break
		}
	}
	s.cursor = end

	str := s.runes[i:end]
	if runeSliceEq(str, []rune("let")) {
		return Let, str
	}
	if runeSliceEq(str, []rune("type")) {
		return Type, str
	}
	if runeSliceEq(str, []rune("mut")) {
		return Mut, str
	}
	if runeSliceEq(str, []rune("restricted")) {
		return Restricted, str
	}
	if runeSliceEq(str, []rune("lazy")) {
		return Lazy, str
	}
	if runeSliceEq(str, []rune("extend")) {
		return Extend, str
	}
	if runeSliceEq(str, []rune("return")) {
		return Return, str
	}
	if runeSliceEq(str, []rune("while")) {
		return While, str
	}
	if runeSliceEq(str, []rune("struct")) {
		return Struct, str
	}
	if runeSliceEq(str, []rune("namespace")) {
		return Namespace, str
	}
	if runeSliceEq(str, []rune("import")) {
		return Import, str
	}
	if runeSliceEq(str, []rune("if")) {
		return If, str
	}
	if runeSliceEq(str, []rune("else")) {
		return Else, str
	}
	if runeSliceEq(str, []rune("match")) {
		return Match, str
	}
	if runeSliceEq(str, []rune("as")) {
		return As, str
	}
	if runeSliceEq(str, []rune("is")) {
		return Is, str
	}
	if runeSliceEq(str, []rune("true")) {
		return Boolean, str
	}
	if runeSliceEq(str, []rune("false")) {
		return Boolean, str
	}

	return Identifier, str
}

func (s *TokenReader) readBracket() (tok TokenType, text []rune) {
	str := s.read()
	switch str {
	case '(':
		return LParen, []rune{str}
	case ')':
		return RParen, []rune{str}
	case '{':
		return LBrace, []rune{str}
	case '}':
		return RBrace, []rune{str}
	case '<':
		return LAngle, []rune{str}
	case '>':
		return RAngle, []rune{str}
	case '[':
		return LSquare, []rune{str}
	case ']':
		return RSquare, []rune{str}
	}
	return Illegal, []rune{str}
}

func (s *TokenReader) readSymbol() (tok TokenType, text []rune) {
	ch := s.read()

	switch ch {
	case '.':
		return Dot, []rune{ch}
	case '=':
		peeked := s.peek()
		if peeked == '>' {
			s.read()
			return Arrow, []rune{ch, peeked}
		}
		if peeked == '=' {
			s.read()
			return Equals, []rune{ch, peeked}
		}
		return Equal, []rune{ch}
	}

	return Illegal, []rune{ch}
}
func (s *TokenReader) readAngleBracket() (tok TokenType, text []rune) {
	ch1 := s.read()
	ch := s.peek()
	if ch1 == '<' {
		switch ch {
		case '=':
			s.read()
			return LesserEqual, []rune{ch1, ch}
		}
		return LAngle, []rune{ch1}
	}
	if ch1 == '>' {
		switch ch {
		case '=':
			s.read()
			return GreaterEqual, []rune{ch1, ch}
		}
		return RAngle, []rune{ch1}
	}

	return Illegal, []rune{ch1}
}
func (s *TokenReader) readOperator() (tok TokenType, text []rune) {
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
		return Add, []rune(str)
	case "-":
		return Subtract, []rune(str)
	case "*":
		return Multiply, []rune(str)
	case "/":
		return Slash, []rune(str)
	case "%":
		return Mod, []rune(str)
	case "&&":
		return And, []rune(str)
	case "||":
		return Or, []rune(str)
	case "^":
		return Xor, []rune(str)
	case "==":
		return Equals, []rune(str)
	case "!=":
		return NotEquals, []rune(str)
	case ">=":
		return GreaterEqual, []rune(str)
	case "<=":
		return LesserEqual, []rune(str)
	case "!":
		return Not, []rune(str)

		//Dirty hack, these 2 should probably be in readBracket but oh well...
	case ">":
		return LAngle, []rune(str)
	case "<":
		return RAngle, []rune(str)
	}

	return Illegal, []rune(str)
}

func (s *TokenReader) readString() (tok TokenType, text []rune) {
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
	return String, []rune(buf.String())
}

func (s *TokenReader) readNumber() (tok TokenType, text []rune) {
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

	return numType, []rune(buf.String())
}

func (s *TokenReader) readNumberNew() (tok TokenType, text []rune) {
	i := s.cursor
	end := i + 1
	numType := Int

	for {
		r := s.runes[end]
		if r == eof {
			break
		}
		if r == '\n' {
			s.unread()
			break
		}
		if r == '.' {
			numType = Float
		} else if !unicode.IsNumber(r) {
			s.unread()
			break
		}
		end++
		if end >= len(s.runes) {
			break
		}
	}
	s.cursor = end

	return numType, s.runes[i:end]
}

func runeSliceEq(a []rune, b []rune) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
