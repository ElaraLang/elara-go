package parserlegacy

import (
	"github.com/ElaraLang/elara/lexer"
	"regexp"
)

type NamespaceStmt struct {
	Namespace string
}

func (NamespaceStmt) stmtNode() {}

type ImportStmt struct {
	Imports []string
}

func (ImportStmt) stmtNode() {}

var namespaceRegex, _ = regexp.Compile(".+/.+")

func (p *Parser) parseFileMeta() (NamespaceStmt, ImportStmt) {
	p.consume(lexer.Namespace, "Expected file namespace declaration!")
	nsToken := p.consume(lexer.Identifier, "Expected valid namespace!")
	ns := string(nsToken.Data)
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
		impNs = string(importToken.Data)
		if !namespaceRegex.MatchString(impNs) {
			panic(ParseError{
				token:   nsToken,
				message: "Invalid namespace format to import",
			})
		}
		imports = append(imports, impNs)
		p.cleanNewLines()
	}
	return NamespaceStmt{
			Namespace: ns,
		}, ImportStmt{
			Imports: imports,
		}
}
