package parser

import (
	"elara/lexer"
	"regexp"
)

type MetaInfoStmt struct {
	Namespace string
	Imports   []string
}

func (MetaInfoStmt) stmtNode() {}

var namespaceRegex, _ = regexp.Compile(".+/.+")

func (p *Parser) parseFileMeta() MetaInfoStmt {
	p.consume(lexer.Namespace, "Expected file namespace declaration!")
	nsToken := p.consume(lexer.Identifier, "Expected valid namespace!")
	ns := string(nsToken.Text)
	if !namespaceRegex.MatchString(ns) {
		panic(ParseError{
			token:   nsToken,
			message: "Invalid namespace format",
		})
	}
	p.cleanNewLines()
	imports := make([]string, 0)
	var impNs string
	for p.match(lexer.Import) {
		importToken := p.consume(lexer.Identifier, "Expected valid namespace to import!")
		impNs = string(importToken.Text)
		if !namespaceRegex.MatchString(impNs) {
			panic(ParseError{
				token:   nsToken,
				message: "Invalid namespace format to import",
			})
		}
		imports = append(imports, impNs)
		p.cleanNewLines()
	}
	return MetaInfoStmt{
		Namespace: ns,
		Imports:   imports,
	}
}
