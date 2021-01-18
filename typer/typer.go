package typer

import "github.com/ElaraLang/elara/parserlegacy"

type Typer struct {
	Input []parserlegacy.Stmt
}

func (t *Typer) HandleTyping() {
	// Pass 1 - Scanning for types
	types := t.scanForUserDefinedTypes()
	// Pass 2 - Scanning for function return types
	functionReturns := t.scanForFunctionReturns()
}

func (t *Typer) scanForUserDefinedTypes() []parserlegacy.Type {
	types := make([]parserlegacy.Type, 0)
	var cur parserlegacy.Stmt
	for index := range t.Input {
		cur = t.Input[index]
		switch cur.(type) {
		case parserlegacy.StructDefStmt:
			structDef := cur.(parserlegacy.StructDefStmt)
			appendType(&types, &structDef)
		}
	}
	return types
}

func (t *Typer) scanForFunctionReturns() []parserlegacy.Type {

}

func appendType(types *[]parserlegacy.Type, p *parserlegacy.StructDefStmt) {
	fields := make([]parserlegacy.DefinedType, 0)
	for fieldIndex := range p.StructFields {
		curField := p.StructFields[fieldIndex]
		field := parserlegacy.DefinedType{
			Identifier: curField.Identifier,
			DefType:    *curField.FieldType,
		}
		fields = append(fields, field)
	}
	*types = append(*types, parserlegacy.DefinedTypeContract{DefType: fields, Name: p.Identifier})
}
