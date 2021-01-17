package interpreter

import (
	"fmt"
	"github.com/ElaraLang/elara/util"
)

var AnyType = NewEmptyType("Any")
var UnitType = NewEmptyType("Unit")

var FloatType = NewEmptyType("Float")
var BooleanType = NewEmptyType("Boolean")

var CharType = NewEmptyType("Char")

var StringType = NewCollectionTypeOf(CharType)
var OutputType = NewEmptyType("Output")

var types = []Type{
	AnyType,
	UnitType,
	IntType,
	FloatType,
	BooleanType,
	StringType,
	CharType,
	OutputType,
}

func Init(context *Context) {
	for _, t := range types {
		context.types[t.Name()] = t
	}
	context.types["String"] = StringType

	InitInts(context)

	stringPlusName := "plus"
	stringPlus := &Function{
		Signature: Signature{
			Parameters: []Parameter{
				{
					Name:     "this",
					Type:     StringType,
					Position: 0,
				},
				{
					Name:     "other",
					Type:     AnyType,
					Position: 1,
				}},
			ReturnType: StringType,
		},
		Body: NewAbstractCommand(func(ctx *Context) *ReturnedValue {
			this := ctx.FindParameter(0)
			otherParam := ctx.FindParameter(1)
			concatenated := this.Value.(*Collection).elemsAsString() + util.Stringify(otherParam.Value)
			return NonReturningValue(&Value{
				Type:  StringType,
				Value: concatenated,
			})
		}),
		name: &stringPlusName,
	}
	stringPlusType := NewFunctionType(stringPlus)
	context.DefineVariable(&Variable{
		Name:    stringPlusName,
		Mutable: false,
		Type:    stringPlusType,
		Value: &Value{
			Type:  stringPlusType,
			Value: stringPlus,
		},
	})

	anyPlusName := "plus"
	anyPlus := &Function{
		Signature: Signature{
			Parameters: []Parameter{
				{
					Name: "this",
					Type: AnyType,
				},
				{
					Name:     "other",
					Type:     StringType,
					Position: 1,
				}},
			ReturnType: StringType,
		},
		Body: NewAbstractCommand(func(ctx *Context) *ReturnedValue {
			this := ctx.FindParameter(0)
			otherParam := ctx.FindParameter(1)

			concatenated := ctx.Stringify(this) + otherParam.Value.(*Collection).elemsAsString()
			return NonReturningValue(&Value{
				Type:  StringType,
				Value: concatenated,
			})
		}),
		name: &anyPlusName,
	}
	anyPlusType := NewFunctionType(anyPlus)
	context.DefineVariable(&Variable{
		Name:    anyPlusName,
		Mutable: false,
		Type:    anyPlusType,
		Value: &Value{
			Type:  anyPlusType,
			Value: anyPlus,
		},
	})

	define(context, "toString", &Function{
		Signature: Signature{
			Parameters: []Parameter{
				{
					Name: "this",
					Type: AnyType,
				},
			},
			ReturnType: StringType,
		},
		Body: NewAbstractCommand(func(ctx *Context) *ReturnedValue {
			this := ctx.FindParameter(0)
			return NonReturningValue(StringValue(ctx.Stringify(this)))
		}),
	})

	colPlusName := "plus"
	colPlus := &Function{
		Signature: Signature{
			Parameters: []Parameter{
				{
					Name: "this",
					Type: NewCollectionTypeOf(AnyType),
				},
				{
					Name:     "other",
					Type:     NewCollectionTypeOf(AnyType),
					Position: 1,
				},
			},
			ReturnType: NewCollectionTypeOf(AnyType),
		},
		Body: NewAbstractCommand(func(ctx *Context) *ReturnedValue {
			this := ctx.FindParameter(0).Value.(*Collection)
			other := ctx.FindParameter(1).Value.(*Collection)

			elems := make([]*Value, len(this.Elements)+len(other.Elements))
			for i, element := range this.Elements {
				elems[i] = element
			}
			for i, element := range other.Elements {
				elems[i+len(this.Elements)] = element
			}

			newCol := &Collection{
				ElementType: this.ElementType,
				Elements:    elems,
			}
			return NonReturningValue(&Value{
				Type:  NewCollectionType(newCol),
				Value: newCol,
			})
		}),
		name: &colPlusName,
	}
	colPlusType := NewFunctionType(colPlus)
	context.DefineVariable(&Variable{
		Name:    colPlusName,
		Mutable: false,
		Type:    colPlusType,
		Value: &Value{
			Type:  colPlusType,
			Value: colPlus,
		},
	})

	define(context, "times", &Function{
		Signature: Signature{
			Parameters: []Parameter{
				{
					Name: "this",
					Type: NewCollectionTypeOf(AnyType),
				},
				{
					Name:     "other",
					Type:     IntType,
					Position: 1,
				},
			},
			ReturnType: NewCollectionTypeOf(AnyType),
		},
		Body: NewAbstractCommand(func(ctx *Context) *ReturnedValue {
			thisParam := ctx.FindParameter(0)
			this := thisParam.Value.(*Collection)
			amount := ctx.FindParameter(1).Value.(int64)
			if amount == 1 {
				return NonReturningValue(thisParam)
			}

			currentElemLength := len(this.Elements)
			newSize := int64(currentElemLength) * amount
			newColl := make([]*Value, newSize)
			for i := int64(0); i < newSize; i++ {
				index := i % amount
				if currentElemLength == 1 {
					index = 0
				}
				newColl[i] = this.Elements[index]
			}

			collection := &Collection{
				ElementType: this.ElementType,
				Elements:    newColl,
			}
			return NonReturningValue(NewValue(NewCollectionType(collection), collection))
		}),
	})

	outputWriteName := "write"
	outputWrite := &Function{
		Signature: Signature{
			Parameters: []Parameter{
				{
					Name: "this",
					Type: OutputType,
				},
				{
					Name:     "value",
					Type:     AnyType,
					Position: 1,
				},
			},
			ReturnType: UnitType,
		},
		Body: NewAbstractCommand(func(ctx *Context) *ReturnedValue {
			value := ctx.FindParameter(1)
			asString := ctx.Stringify(value)

			fmt.Printf("%s", asString)
			return NonReturningValue(UnitValue())
		}),
		name: &outputWriteName,
	}
	outputWriteType := NewFunctionType(stringPlus)
	context.DefineVariable(&Variable{
		Name:    outputWriteName,
		Mutable: false,
		Type:    outputWriteType,
		Value: &Value{
			Type:  outputWriteType,
			Value: outputWrite,
		},
	})

	anyEqualsName := "equals"
	anyEquals := &Function{
		Signature: Signature{
			Parameters: []Parameter{
				{
					Name: "this",
					Type: AnyType,
				},
				{
					Name:     "other",
					Type:     AnyType,
					Position: 1,
				},
			},
			ReturnType: BooleanType,
		},
		name: &anyEqualsName,
		Body: NewAbstractCommand(func(c *Context) *ReturnedValue {
			this := c.FindParameter(0).Value
			other := c.FindParameter(1).Value
			switch a := this.(type) {
			case *Collection:
				otherAsCol, otherIsCol := other.(*Collection)
				if !otherIsCol {
					break
				}
				if len(a.Elements) != len(otherAsCol.Elements) {
					break
				}
				for i, element := range a.Elements {
					otherElem := otherAsCol.Elements[i]
					if !element.Equals(c, otherElem) {
						break
					}
				}
				return NonReturningValue(BooleanValue(true))
			}
			return NonReturningValue(BooleanValue(this == other))
		}),
	}
	anyEqualsType := NewFunctionType(anyEquals)
	context.DefineVariable(&Variable{
		Name:    anyEqualsName,
		Mutable: false,
		Type:    anyEqualsType,
		Value: &Value{
			Type:  anyEqualsType,
			Value: anyEquals,
		},
	})
}

//func intAdd(ctx *Context) *Value {
//	parameter := ctx.FindParameter("value")
//	asInt, isInt := parameter.Value.(int64)
//	if isInt {
//		result := ctx.receiver.Value.(int64) + asInt
//		return &Value{
//			Type:  IntType,
//			Value: result,
//		}
//	} else {
//		asFloat, isFloat := parameter.Value.(float64)
//		if isFloat {
//			result := float64(ctx.receiver.Value.(int64)) + asFloat
//			return &Value{
//				Type:  FloatType,
//				Value: result,
//			}
//		} else {
//			//TODO
//			//While this might work, it ignores the fact that values won't be "cast" if passed. An Int passed as Any will still try and use Int functions
//			result := util.Stringify(ctx.receiver.Value) + util.Stringify(parameter.Value)
//			return &Value{
//				Type:  StringType,
//				Value: result,
//			}
//		}
//	}
//}
