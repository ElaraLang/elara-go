package parser

import (
	"github.com/ElaraLang/elara/ast"
	"github.com/ElaraLang/elara/lexer"
)

func (p *Parser) parseMapEntries() []ast.Entry {
	entries := make([]ast.Entry, 0)
	for !p.Tape.ValidationPeek(0, lexer.RBrace) {
		key := p.parseExpression(LOWEST)
		p.Tape.Expect(lexer.Colon)
		value := p.parseExpression(LOWEST)
		entries = append(entries, ast.Entry{Key: key, Value: value})
	}
	return entries
}
