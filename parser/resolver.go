package parser

import (
	"github.com/ElaraLang/elara/ast"
	"github.com/ElaraLang/elara/lexer"
)

type resolver = func(p *Parser) bool

func (p *Parser) resolvingParslet(resolutionConditions map[*resolver]*prefixParslet) prefixParslet {
	return func() ast.Expression {
		for resolver, parslet := range resolutionConditions {
			if (*resolver)(p) {
				return (*parslet)()
			}
		}
		return nil
	}
}

func (p *Parser) functionGroupResolver() map[*resolver]*prefixParslet {
	checkFuncDef := isFunctionDefinition
	checkGroupExpression := not(isFunctionDefinition)

	var functionParslet prefixParslet = p.parseFunction
	var groupParslet prefixParslet = p.parseGroupExpression

	functionGroupResolver := map[*resolver]*prefixParslet{
		&(checkFuncDef):         &functionParslet,
		&(checkGroupExpression): &groupParslet,
	}
	return functionGroupResolver
}

func isFunctionDefinition(p *Parser) bool {
	closingIndex := p.Tape.FindDepthClosingIndex(lexer.LParen, lexer.RParen)
	return p.Tape.ValidationPeek(closingIndex+1, lexer.Arrow)
}

func not(predicate func(p *Parser) bool) func(p *Parser) bool {
	return func(p *Parser) bool {
		return !predicate(p)
	}
}
