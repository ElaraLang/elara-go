package lexer

func isHexDigit(n rune) bool {
	return isDecimalDigit(n) || n >= 'A' && n <= 'F' || n >= 'a' && n <= 'f'
}
func isBinaryDigit(n rune) bool {
	return n == '0' || n == '1'
}

func isDecimalDigit(n rune) bool {
	return n >= '0' && n <= '9'
}

func isIllegalIdentifierChar(n rune) bool {
	_, present := illegalIdentifiers[n]
	return present
}

var illegalIdentifiers = map[rune]struct{}{
	' ':  {},
	'(':  {},
	')':  {},
	'{':  {},
	'}':  {},
	'<':  {},
	'>':  {},
	'[':  {},
	']':  {},
	',':  {},
	'.':  {},
	'\n': {},
	'\'': {},
	'"':  {},
}
