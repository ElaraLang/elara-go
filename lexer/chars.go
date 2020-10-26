package lexer

import "unicode"

func IsWhitespace(ch rune) bool {
	return unicode.IsSpace(ch)
}

func isBracket(ch rune) bool {
	return ch == '(' || ch == ')' || ch == '{' || ch == '}' || ch == '[' || ch == ']' || isAngleBracket(ch)
}

func isStartOfSymbol(ch rune) bool {
	return ch == '.' || ch == '='
}

func isAngleBracket(ch rune) bool {
	return ch == '<' || ch == '>'
}

//This function is a bit of a hotspot, mostly due to how often it's called. Not much to be done here though - map access is pretty fast :/
func isValidIdentifier(ch rune) bool {
	s := IllegalIdentifierCharSlice[ch]
	return !s
}

func isNumerical(ch rune) bool {
	return unicode.IsNumber(ch)
}

func isOperatorSymbol(ch rune) bool {
	return ch == '=' || ch == '+' || ch == '-' || ch == '*' || ch == '/' || ch == '%' || ch == '&' || ch == '|' || ch == '^' || ch == '!' || ch == '>' || ch == '<'
}

var eof = rune(-1)
