package parser

import (
	"github.com/ElaraLang/elara/ast"
)

func collectParserResult(p *Parser) ([]ast.Statement, []ParseError) {
	resultChannel := make(chan interface{})
	go collectParserOutput(p, resultChannel)
	go collectParserErrors(p, resultChannel)

	var resStmts []ast.Statement
	var resErrors []ParseError

	resA := <-resultChannel
	resB := <-resultChannel

	switch r := resA.(type) {
	case []ast.Statement:
		resStmts = r
		resErrors = resB.([]ParseError)
	case []ParseError:
		resErrors = r
		resStmts = resB.([]ast.Statement)
	default:
		close(resultChannel)
		panic("Invalid result collected by result channel!")
	}

	close(resultChannel)
	return resStmts, resErrors
}

func collectParserOutput(p *Parser, resChan chan interface{}) {
	res := make([]ast.Statement, 0)
loop:
	for !(p.Tape.IsClosed() && p.Tape.index == len(p.Tape.tokens)-1) {
		select {
		case stmt, ok := <-p.OutputChannel:
			if ok {
				res = append(res, stmt)
			} else {
				break loop
			}
		default:
		}
	}
	resChan <- res
}

func collectParserErrors(p *Parser, resChan chan interface{}) {
	errors := make([]ParseError, 0)
loop:
	for !(p.Tape.IsClosed() && p.Tape.index == len(p.Tape.tokens)-1) {
		select {
		case err, ok := <-p.ErrorChannel:
			if ok {
				errors = append(errors, err)
			} else {
				break loop
			}
		default:
		}
	}
	resChan <- errors
}
