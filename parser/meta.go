package parser

import "elara/lexer"

type MetaInfoStmt struct {
	Package string
	Imports []string
}

func (MetaInfoStmt) stmtNode() {}

func (p *Parser) parseFileMeta() MetaInfoStmt {
	p.consume(lexer.Namespace, "Expected file namespace declaration!")
	ns := string(p.consume(lexer.String, "Expected valid namespace!").Text)
	p.cleanNewLines()
	imports := make([]string, 0)
	var impNs string
	for p.match(lexer.Import) {
		impNs = string(p.consume(lexer.String, "Expected valid namespace to import!").Text)
		imports = append(imports, impNs)
		p.cleanNewLines()
	}
	return MetaInfoStmt{
		Package: ns,
		Imports: imports,
	}
}
