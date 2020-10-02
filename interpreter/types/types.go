package types

import (
	"elara/parser"
)

var IntType Type = &SimpleType{name: "Int"}
var FloatType Type = &SimpleType{name: "Float"}
var StringType Type = &SimpleType{name: "String"}
var UnitType Type = &SimpleType{name: "Unit"}
var AnyType Type = &SimpleType{name: "Any"}

func FromASTType(p parser.Type) Type {
	switch p.(type) {
	case parser.ElementaryTypeContract:
		return &SimpleType{name: p.(parser.ElementaryTypeContract).Identifier}

	}

	return nil
}
