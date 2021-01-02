package interpreter

import (
	"github.com/ElaraLang/elara/parser"
	"reflect"
)

type Type interface {
	Name() string
	//returns if *this* type accepts the other type
	Accepts(otherType Type) bool
}

type StructType struct {
	TypeName   string
	Properties map[string]Property
}

func (t *StructType) Name() string {
	return t.TypeName
}

func (t *StructType) Accepts(otherType Type) bool {
	elem, ok := otherType.(*StructType)
	if !ok {
		return false
	}
	for s, property := range t.Properties {
		other, exists := elem.Properties[s]
		if !exists {
			return false
		}
		if !property.Type.Accepts(other.Type) {
			return false
		}
	}
	return true
}

type Property struct {
	Name string
	Type Type
	//bitmask (base/modifiers.go)
	Modifiers    uint
	DefaultValue *Value
}

type FunctionType struct {
	Signature Signature
}

func NewFunctionType(function *Function) *FunctionType {
	return &FunctionType{Signature: function.Signature}
}

func (t *FunctionType) Name() string {
	return t.Signature.String()
}

/*
Function acceptance is defined by having the same number of parameters,
with all of A's parameters accepting the corresponding parameters for B
and A's return type accepting B's return type
*/
func (t *FunctionType) Accepts(other Type) bool {
	otherFunc, ok := other.(*FunctionType)
	if !ok {
		return false
	}

	if len(t.Signature.Parameters) != len(otherFunc.Signature.Parameters) {
		return false
	}
	for i, parameter := range t.Signature.Parameters {
		otherParam := otherFunc.Signature.Parameters[i]
		if !parameter.Type.Accepts(otherParam.Type) {
			return false
		}
	}
	return t.Signature.ReturnType.Accepts(otherFunc.Signature.ReturnType)
}

type CollectionType struct {
	ElementType Type
}

func (t *CollectionType) Name() string {
	return t.ElementType.Name() + "[]" //Eg String[]
}

func (t *CollectionType) Accepts(other Type) bool {
	otherColl, ok := other.(*CollectionType)
	if !ok {
		return false
	}

	return t.ElementType.Accepts(otherColl.ElementType)
}

//TODO mapType

type EmptyType struct {
	name string
}

func (t *EmptyType) Name() string {
	return t.name
}

func (t *EmptyType) Accepts(other Type) bool {
	if t == AnyType {
		return true
	}
	return t == other
}

func NewEmptyType(name string) Type {
	return &EmptyType{name: name}
}

func FromASTType(astType parser.Type, ctx *Context) Type {
	switch t := astType.(type) {
	case parser.ElementaryTypeContract:
		found := ctx.FindType(t.Identifier)
		if found != nil {
			return found
		}
		return NewEmptyType(t.Identifier)
	}
	println("Cannot handle " + reflect.TypeOf(astType).Name())
	return nil
}
