package parser

import (
	"github.com/ElaraLang/elara/ast"
	"github.com/ElaraLang/elara/lexer"
)

func (p *Parser) parseMapEntries() []ast.Entry {
	entries := make([]ast.Entry, 0)
	for !p.Tape.ValidationPeek(0, lexer.RBrace) {
		if len(entries) > 0 {
			p.Tape.Expect(lexer.Comma)
		}
		p.Tape.skipLineBreaks()
		key := p.parseExpression(Lowest)
		p.Tape.Expect(lexer.Equal)
		value := p.parseExpression(Lowest)
		entries = append(entries, ast.Entry{Key: key, Value: value})
	}
	return entries
}
