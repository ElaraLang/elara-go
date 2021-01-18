package parser

import "github.com/ElaraLang/elara/lexer"

type Parser struct {
	Tape TokenTape
}

func NewParser(tokens []lexer.Token) Parser {
	return Parser{Tape: NewTokenTape(tokens)}
}

func NewReplParser(tokens []lexer.Token) Parser {
	return Parser{Tape: NewReplTokenTape()}
}
