package parser

import "github.com/ElaraLang/elara/lexer"

type ParseError struct {
	ErrorToken lexer.Token
	Message    string
}

func NewParseError(errorTok lexer.Token, message string) ParseError {
	return ParseError{
		ErrorToken: errorTok,
		Message:    message,
	}
}

func (p *Parser) error(errorTok lexer.Token, message string) {
	panic(NewParseError(errorTok, message))
}

func expect(expectation string, location string) string {
	return "Expecting " + expectation + " at " + location + "."
}

func (p *Parser) handleParseError() {
	if r := recover(); r != nil {
		switch err := r.(type) {
		case ParseError:
			p.ErrorChannel <- err
			break
		case []ParseError:
			for _, v := range err {
				p.ErrorChannel <- v
			}
		default:
			p.ErrorChannel <- NewParseError(p.Tape.Current(), "Invalid error thrown by parser")
			break
		}
		p.syncWithError()
	}
}

func (p *Parser) syncWithError() {
	for !p.Tape.ValidateHead(lexer.NEWLINE, lexer.EOF) {
		p.Tape.advance()
	}
	p.Tape.skipLineBreaks()
}
