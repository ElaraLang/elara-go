package types

import (
	"elara/parser"
)

var IntType Type = &SimpleType{name: "Int"}
var FloatType Type = &SimpleType{name: "Float"}
var StringType Type = &SimpleType{name: "String"}
var UnitType Type = &SimpleType{name: "Unit"}
var AnyType Type = &SimpleAnonymousType{name: "Any", accepts: func(other Type) bool {
	return true
}}

var types = map[string]Type{
	"Int":     IntType,
	"Float":   FloatType,
	"String":  StringType,
	"Unit":    UnitType,
	"AnyType": AnyType,
}

func FromASTType(p parser.Type) Type {
	switch p.(type) {
	case parser.ElementaryTypeContract:
		identifier := p.(parser.ElementaryTypeContract).Identifier
		existing := types[identifier]
		if existing != nil {
			return existing
		}

		newSimpleType := &SimpleType{name: identifier}
		types[identifier] = newSimpleType
		return newSimpleType
		//TODO This will break things if 2 types have the same name!
	}

	return nil
}
