package parser

import (
	"github.com/ElaraLang/elara/ast"
	"github.com/ElaraLang/elara/lexer"
)

type Resolver = func(p *Parser) bool

func (p *Parser) createResolvingParslet(resolutionConditions map[*Resolver]*parsePrefix) parsePrefix {

	return func() ast.Expression {
		for resolver, parslet := range resolutionConditions {
			if (*resolver)(p) {
				return (*parslet)()
			}
		}
		return nil
	}
}

func (p *Parser) createFunctionGroupResolver() map[*Resolver]*parsePrefix {
	checkFuncDef := isFunctionDefinition
	checkGroupExpression := not(isFunctionDefinition)
	var functionParslet parsePrefix = p.parseFunction
	var groupParslet parsePrefix = p.parseGroupExpression
	functionGroupResolver := map[*Resolver]*parsePrefix{
		&checkFuncDef:           &functionParslet,
		&(checkGroupExpression): &groupParslet,
	}
	return functionGroupResolver
}

func isFunctionDefinition(p *Parser) bool {
	closingIndex := p.Tape.FindDepthClosingIndex(lexer.LParen, lexer.RParen)
	return p.Tape.ValidationPeek(closingIndex+1, lexer.Arrow)
}

func not(fn func(p *Parser) bool) func(p *Parser) bool {
	return func(p *Parser) bool {
		return !fn(p)
	}
}
