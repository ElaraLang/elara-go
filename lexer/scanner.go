package lexer

import (
	"bufio"
	"bytes"
	"io"
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

func (s *Scanner) Scan() (tok Token, text string) {
	ch := s.read()

	if ch == eof {
		return EOF, string(ch)
	}
	if isWhitespace(ch) {
		s.unread()
		return s.readWhitespace()
	}
	if isValidIdentifier(ch) {
		s.unread()
		return s.readIdentifier()
	}

	return ILLEGAL, string(ch)
}

func (s *Scanner) readWhitespace() (tok Token, text string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return WS, buf.String()
}

func (s *Scanner) readIdentifier() (tok Token, text string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isValidIdentifier(ch) {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	str := buf.String()
	switch str {
	case "let":
		return LET, str
	case "mut":
		return MUT, str
	case "=":
		return EQUAL, str
	}

	return IDENTIFIER, str
}
