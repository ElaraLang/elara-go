package lexer

import "unicode"

func isWhitespace(ch rune) bool {
	return unicode.IsSpace(ch)
}

func isBracket(ch rune) bool {
	return ch == '(' || ch == ')' || ch == '{' || ch == '}' || ch == '<' || ch == '>'
}

func isStartOfSymbol(ch rune) bool {
	return ch == '.' || ch == '='
}

func isValidIdentifier(ch rune) bool {
	return !IllegalIdentifierChars[ch] && !isWhitespace(ch)
}

func isNumerical(ch rune) bool {
	return unicode.IsNumber(ch)
}

func isOperatorSymbol(ch rune) bool {
	return ch == '=' || ch == '+'
}

var eof = rune(-1)
