package interpreter

import (
	"fmt"
	"github.com/ElaraLang/elara/ast"
	"github.com/ElaraLang/elara/lexer"
	"reflect"
)

type Type interface {
	Name() string
	//returns if *this* type accepts the other type
	Accepts(otherType Type, ctx *Context) bool
}

type StructType struct {
	TypeName          string
	Properties        []Property     //This preserves ordering of properties
	propertyPositions map[string]int //And this guarantees constant lookup still
	constructor       *Value         //*Function of the constructor
}

func (t *StructType) Name() string {
	return t.TypeName
}
func (t *StructType) Accepts(otherType Type, ctx *Context) bool {
	otherStruct, ok := otherType.(*StructType)
	if !ok {
		return false
	}
	for _, property := range t.Properties {
		byName, exists := otherStruct.GetProperty(property.Name)
		if !exists {
			return false //Must have all of the properties
		}
		if !property.Type.Accepts(byName.Type, ctx) {
			return false //And the types must be acceptable
		}
	}
	return true
}
func (t *StructType) GetProperty(identifier string) (Property, bool) {
	i, present := t.propertyPositions[identifier]
	if !present {
		return Property{}, false
	}
	return t.Properties[i], true
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
func NewSignatureFunctionType(signature Signature) *FunctionType {
	return &FunctionType{Signature: signature}
}
func (t *FunctionType) Name() string {
	return t.Signature.String()
}

/*
Function acceptance is defined by having the same number of parameters,
with all of A's parameters accepting the corresponding parameters for B
and A's return type accepting B's return type
*/
func (t *FunctionType) Accepts(otherType Type, ctx *Context) bool {
	otherFunc, ok := otherType.(*FunctionType)
	if !ok {
		return false
	}
	return t.Signature.Accepts(&otherFunc.Signature, ctx, false)
}

type EmptyType struct {
	name string
}

func (t *EmptyType) Name() string {
	return t.name
}
func (t *EmptyType) Accepts(otherType Type, _ *Context) bool {
	if *t == *(AnyType.(*EmptyType)) { //ew
		return true
	}
	//This is really trying to patch a deeper problem - this function relies on there only ever being 1 pointer to a type.
	asEmpty, isEmpty := otherType.(*EmptyType)

	if isEmpty {
		return t.name == asEmpty.name
	}
	return t == otherType
}
func NewEmptyType(name string) Type {
	return &EmptyType{name: name}
}

type UnionType struct {
	a Type
	b Type
}

func (t *UnionType) Name() string {
	return t.a.Name() + " | " + t.b.Name()
}
func (t *UnionType) Accepts(otherType Type, ctx *Context) bool {
	return t.a.Accepts(otherType, ctx) || t.b.Accepts(otherType, ctx)
}

type IntersectionType struct {
	a Type
	b Type
}

func (t *IntersectionType) Name() string {
	return t.a.Name() + " & " + t.b.Name()
}
func (t *IntersectionType) Accepts(otherType Type, ctx *Context) bool {
	return t.a.Accepts(otherType, ctx) && t.b.Accepts(otherType, ctx)
}

type DefinedType struct {
	parts map[string]Type
}

func (t *DefinedType) Name() string {
	return "TODO DefinedType name"
}
func (t *DefinedType) Accepts(other Type, ctx *Context) bool {
	asStruct, isStruct := other.(*StructType)
	for s, t2 := range t.parts {
		if isStruct {
			property, present := asStruct.GetProperty(s)
			if present && t2.Accepts(property.Type, ctx) {
				continue
			}
		}
		extension := ctx.FindExtension(other, s)
		if extension != nil && t2.Accepts(extension.Value.Type, ctx) {
			continue
		}
		return false
	}

	return true
}

func FromASTType(astType ast.Type, ctx *Context) Type {
	switch t := astType.(type) {
	case *ast.PrimaryType:
		found := ctx.FindType(t.Identifier.Name)
		if found != nil {
			return found
		}
		return NewEmptyType(t.Identifier.Name)

	case *ast.FunctionType:
		returned := FromASTType(t.ReturnType, ctx)
		args := make([]Parameter, len(t.ParamTypes))
		for i, arg := range t.ParamTypes {
			argType := FromASTType(arg, ctx)
			args[i] = Parameter{
				Name: fmt.Sprintf("arg%d", i),
				Type: argType,
			}
		}

		signature := Signature{
			Parameters: args,
			ReturnType: returned,
		}
		return NewSignatureFunctionType(signature)

	case *ast.CollectionType:
		elemType := FromASTType(t.Type, ctx)
		return &CollectionType{
			ElementType: elemType,
		}

	case *ast.AlgebraicType:
		switch t.Operation.TokenType {
		case lexer.TypeAnd:
			return &IntersectionType{
				a: FromASTType(t.Left, ctx),
				b: FromASTType(t.Right, ctx),
			}
		case lexer.TypeOr:
			return &UnionType{
				a: FromASTType(t.Left, ctx),
				b: FromASTType(t.Right, ctx),
			}
		}
	case *ast.ContractualType:
		parts := make(map[string]Type, len(t.Contracts))
		for _, definedType := range t.Contracts {
			parts[definedType.Identifier.Name] = FromASTType(definedType.Type, ctx)
		}
		return &DefinedType{
			parts: parts,
		}
	case *ast.MapType:
		keyType := FromASTType(t.KeyType, ctx)
		valueType := FromASTType(t.ValueType, ctx)
		return &MapType{
			KeyType: keyType, ValueType: valueType,
		}
	}
	panic("Cannot handle " + reflect.TypeOf(astType).Name())
	return nil
}
