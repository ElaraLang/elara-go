package parser

import (
	"github.com/ElaraLang/elara/ast"
	"github.com/ElaraLang/elara/lexer"
)

type resolver = func(p *Parser) bool

func (p *Parser) resolvingPrefixParslet(resolutionConditions map[*resolver]*prefixParslet) prefixParslet {
	return func() ast.Expression {
		for resolver, parslet := range resolutionConditions {
			if (*resolver)(p) {
				return (*parslet)()
			}
		}
		return nil
	}
}

func (p *Parser) resolvingTypePrefixParslet(resolutionConditions map[*resolver]*prefixTypeParslet) prefixTypeParslet {
	return func() ast.Type {
		for resolver, parslet := range resolutionConditions {
			if (*resolver)(p) {
				return (*parslet)()
			}
		}
		return nil
	}
}

func (p *Parser) functionGroupTypeResolver() map[*resolver]*prefixTypeParslet {
	checkFuncDef := isFunctionDefinition
	checkGroupExpression := not(isFunctionDefinition)

	var functionParslet prefixTypeParslet = p.parseFunctionType
	var groupParslet prefixTypeParslet = p.parseGroupedType
	return map[*resolver]*prefixTypeParslet{
		&(checkFuncDef):         &functionParslet,
		&(checkGroupExpression): &groupParslet,
	}
}

func (p *Parser) functionGroupResolver() map[*resolver]*prefixParslet {
	checkFuncDef := isFunctionDefinition
	checkGroupExpression := not(isFunctionDefinition)

	var functionParslet prefixParslet = p.parseFunction
	var groupParslet prefixParslet = p.parseGroupExpression

	return map[*resolver]*prefixParslet{
		&(checkFuncDef):         &functionParslet,
		&(checkGroupExpression): &groupParslet,
	}
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
