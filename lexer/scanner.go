package lexer

import (
	"unicode"
)

//TODO remove the amount of `defer` declarations as they seem to have a cost, and more optimisations of readIdentifier
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
		oldCol := s.col
		s.col = 0
		s.line++
		return NEWLINE, []rune{ch}, s.line - 1, oldCol
	}

	if IsWhitespace(ch) {
		s.consumeWhitespace()
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
	count := 1 //We know this function will be called when the lexer has already encountered at least 1 whitespace
	for {
		ch := s.read()
		if ch == eof {
			break
		}
		if ch == '\n' {
			s.line++
			s.col = 0
		} else if ch == '\t' || ch == ' ' {
			count++
		} else {
			s.unread()
			break
		}
	}
	s.col += count
	return count
}

func (s *TokenReader) readIdentifier() (tok TokenType, text []rune) {
	i := s.cursor
	end := i
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
	length := end - i //possibly slightly faster than len()

	switch str[0] {
	case 'l':
		{
			if length == 3 && str[1] == 'e' && str[2] == 't' {
				return Let, str
			}
			if length == 4 && str[1] == 'a' && str[2] == 'z' && str[3] == 'y' {
				return Lazy, str
			}
			return Identifier, str
		}
	case 't':
		{
			if length == 4 {
				if str[1] == 'y' && str[2] == 'p' && str[3] == 'e' {
					return Type, str
				}
				if str[1] == 'r' && str[2] == 'u' && str[3] == 'e' {
					return BooleanTrue, str
				}
			}
			return Identifier, str
		}
	case 'i':
		{
			if length == 2 {
				if str[1] == 'f' {
					return If, str
				}
				if str[1] == 's' {
					return Is, str
				}
			}
			if length == 6 && str[1] == 'm' && str[2] == 'p' && str[3] == 'o' && str[4] == 'r' && str[5] == 't' {
				return Import, str
			}
			return Identifier, str
		}
	}
	//TODO optimise other comparisons
	if runeSliceEq(str, []rune("mut")) {
		return Mut, str
	}
	if runeSliceEq(str, []rune("restricted")) {
		return Restricted, str
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
	if runeSliceEq(str, []rune("else")) {
		return Else, str
	}
	if runeSliceEq(str, []rune("match")) {
		return Match, str
	}
	if runeSliceEq(str, []rune("as")) {
		return As, str
	}
	if runeSliceEq(str, []rune("false")) {
		return BooleanFalse, str
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
	start := s.cursor
	end := start
	for {
		r := s.runes[end]
		if r == eof {
			break
		}
		if !isOperatorSymbol(r) {
			s.unread()
			break
		}
		end++
		if end >= len(s.runes) {
			break
		}
	}
	s.cursor = end

	str := s.runes[start:end]
	switch str[0] {
	case '+':
		return Add, str
	case '-':
		return Subtract, str
	case '*':
		return Multiply, str
	case '/':
		return Slash, str
	case '%':
		return Mod, str
	case '^':
		return Xor, str

	case '>':
		{
			l := len(str)
			if l == 1 {
				return LAngle, str
			}
			n := str[1]
			if l > 2 || n != '=' {
				panic("Unknown operator " + string(str))
			}
			return GreaterEqual, str
		}
	case '<':
		{
			l := len(str)
			if l == 1 {
				return RAngle, str
			}
			n := str[1]
			if l > 2 || n != '=' {
				panic("Unknown operator " + string(str))
			}
			return LesserEqual, str
		}
	case '!':
		{
			l := len(str)
			if l == 1 {
				return Not, str
			}
			n := str[1]
			if l > 2 || n != '=' {
				panic("Unknown operator " + string(str))
			}
			return NotEquals, str
		}
	}
	if len(str) != 2 {
		panic("Unknown operator " + string(str))
	}
	if runeSliceEq(str, []rune("&&")) {
		return And, str
	}
	if runeSliceEq(str, []rune("||")) {
		return Or, str
	}
	if runeSliceEq(str, []rune("==")) {
		return Equals, str
	}
	return Illegal, str
}

//This function is called with the assumption that the beginning " has ALREADY been read.
func (s *TokenReader) readString() (tok TokenType, text []rune) {
	start := s.cursor
	end := start

	for {
		r := s.runes[end]
		if r == eof {
			break
		}
		end++
		if r == '"' {
			break
		}
		if end >= len(s.runes) {
			break
		}
	}
	s.cursor = end
	return String, s.runes[start : end-1]
}

func (s *TokenReader) readNumber() (tok TokenType, text []rune) {
	start := s.cursor
	end := start
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

	return numType, s.runes[start:end]
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
