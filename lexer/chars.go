package lexer

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func isValidIdentifier(ch rune) bool {
	return !IllegalIdentifierChars[ch] && !isWhitespace(ch)
}

var eof = rune(0)
