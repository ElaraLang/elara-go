package typer

import "github.com/ElaraLang/elara/parser"

type Typer struct {
	Input []parser.Stmt
}

func (t *Typer) HandleTyping() {
	// Pass 1 - Scanning for types
	types := t.scanForUserDefinedTypes()
	// Pass 2 - Scanning for function return types
	functionReturns := t.scanForFunctionReturns()
}

func (t *Typer) scanForUserDefinedTypes() []parser.Type {
	types := make([]parser.Type, 0)
	var cur parser.Stmt
	for index := range t.Input {
		cur = t.Input[index]
		switch cur.(type) {
		case parser.StructDefStmt:
			structDef := cur.(parser.StructDefStmt)
			appendType(&types, &structDef)
		}
	}
	return types
}

func (t *Typer) scanForFunctionReturns() []parser.Type {

}

func appendType(types *[]parser.Type, p *parser.StructDefStmt) {
	fields := make([]parser.DefinedType, 0)
	for fieldIndex := range p.StructFields {
		curField := p.StructFields[fieldIndex]
		field := parser.DefinedType{
			Identifier: curField.Identifier,
			DefType:    *curField.FieldType,
		}
		fields = append(fields, field)
	}
	*types = append(*types, parser.DefinedTypeContract{DefType: fields, Name: p.Identifier})
}
