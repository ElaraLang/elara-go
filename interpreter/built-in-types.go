package interpreter

import (
	"fmt"
	"github.com/ElaraLang/elara/util"
)

var AnyType = NewEmptyType("Any")
var UnitType = NewEmptyType("Unit")

var FloatType = NewEmptyType("Float")
var BooleanType = NewEmptyType("Boolean")
var StringType = NewEmptyType("String")
var OutputType = NewEmptyType("Output")

var types = []Type{
	AnyType,
	UnitType,
	IntType,
	FloatType,
	BooleanType,
	StringType,
	OutputType,
}

func Init(context *Context) {
	for _, t := range types {
		context.types[t.Name()] = t
	}
	InitInts(context)

	stringPlusName := "plus"
	stringPlus := &Function{
		Signature: Signature{
			Parameters: []Parameter{
				{
					Name: "this",
					Type: StringType,
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
			concatenated := this.Value.(string) + util.Stringify(otherParam.Value)
			return NonReturningValue(&Value{
				Type:  StringType,
				Value: concatenated,
			})
		}),
		name: &stringPlusName,
	}
	stringPlusType := NewFunctionType(stringPlus)
	context.DefineVariable(stringPlusName, &Variable{
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

			concatenated := ctx.Stringify(this) + otherParam.Value.(string)
			return NonReturningValue(&Value{
				Type:  StringType,
				Value: concatenated,
			})
		}),
		name: &anyPlusName,
	}
	anyPlusType := NewFunctionType(anyPlus)
	context.DefineVariable(anyPlusName, &Variable{
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
	context.DefineVariable(colPlusName, &Variable{
		Name:    colPlusName,
		Mutable: false,
		Type:    colPlusType,
		Value: &Value{
			Type:  colPlusType,
			Value: colPlus,
		},
	})
	//	"to-int": {
	//		Signature: Signature{
	//			Parameters: []Parameter{},
	//			ReturnType: *IntType,
	//		},
	//		Body: NewAbstractCommand(func(ctx *Context) *Value {
	//			value, err := strconv.ParseInt(ctx.receiver.Value.(String), 10, 64)
	//			if err != nil {
	//				panic(err)
	//			}
	//			return &Value{
	//				Type:  IntType,
	//				Value: value,
	//			}
	//		}),
	//	},
	//	"equals": {
	//		Signature: Signature{
	//			Parameters: []Parameter{
	//				{
	//					Name: "value",
	//					Type: *StringType,
	//				},
	//			},
	//			ReturnType: *BooleanType,
	//		},
	//		Body: NewAbstractCommand(func(ctx *Context) *Value {
	//			parameter := ctx.FindParameter("value")
	//			eq := ctx.receiver.Value.(String) == parameter.Value
	//			return &Value{
	//				Type:  BooleanType,
	//				Value: eq,
	//			}
	//		}),
	//	},
	//})
	//
	//BooleanType.variables = convert(map[String]Function{
	//	"and": {
	//		Signature: Signature{
	//			Parameters: []Parameter{
	//				{
	//					Name: "value",
	//					Type: *BooleanType,
	//				},
	//			},
	//			ReturnType: *BooleanType,
	//		},
	//		Body: NewAbstractCommand(func(ctx *Context) *Value {
	//			parameter := ctx.FindParameter("value")
	//			and := ctx.receiver.Value.(bool) && parameter.Value.(bool)
	//			return &Value{
	//				Type:  BooleanType,
	//				Value: and,
	//			}
	//		}),
	//	},
	//	"not": {
	//		Signature: Signature{
	//			Parameters: []Parameter{},
	//			ReturnType: *BooleanType,
	//		},
	//		Body: NewAbstractCommand(func(ctx *Context) *Value {
	//			return &Value{
	//				Type:  BooleanType,
	//				Value: !ctx.receiver.Value.(bool),
	//			}
	//		}),
	//	},
	//	"plus": {
	//		Signature: Signature{
	//			Parameters: []Parameter{
	//				{
	//					Name: "value",
	//					Type: *AnyType,
	//				},
	//			},
	//			ReturnType: *StringType,
	//		},
	//		Body: NewAbstractCommand(func(ctx *Context) *Value {
	//			parameter := ctx.FindParameter("value")
	//			thisStr := util.Stringify(ctx.receiver.Value)
	//			otherStr := util.Stringify(parameter.Value)
	//			return StringValue(thisStr + otherStr)
	//		}),
	//	},
	//})
	//
	//intAdd := Function{
	//	Signature: Signature{
	//		Parameters: []Parameter{
	//			{
	//				Name: "value",
	//				Type: *IntType,
	//			},
	//		},
	//		ReturnType: *IntType,
	//	},
	//	Body: NewAbstractCommand(intAdd),
	//}
	//floatAdd := Function{
	//	Signature: Signature{
	//		Parameters: []Parameter{
	//			{
	//				Name: "value",
	//				Type: *IntType,
	//			},
	//		},
	//		ReturnType: *FloatType,
	//	},
	//	Body: NewAbstractCommand(func(ctx *Context) *Value {
	//		parameter := ctx.FindParameter("value")
	//		asInt, isInt := parameter.Value.(int64)
	//		if isInt {
	//			result := ctx.receiver.Value.(float64) + float64(asInt)
	//			return &Value{
	//				Type:  FloatType,
	//				Value: result,
	//			}
	//		} else {
	//			asFloat, isFloat := parameter.Value.(float64)
	//			if isFloat {
	//				result := ctx.receiver.Value.(float64) + asFloat
	//				return &Value{
	//					Type:  FloatType,
	//					Value: result,
	//				}
	//			} else {
	//				//TODO
	//				//While this might work, it ignores the fact that values won't be "cast" if passed. An Int passed as Any will still try and use Int functions
	//				result := util.Stringify(ctx.receiver.Value) + util.Stringify(parameter.Value)
	//				return &Value{
	//					Type:  StringType,
	//					Value: result,
	//				}
	//			}
	//		}
	//	}),
	//}
	//floatAdd := Function{
	//	Signature: Signature{
	//		Parameters: []Parameter{
	//			{
	//				Name: "value",
	//				Type: *IntType,
	//			},
	//		},
	//		ReturnType: *FloatType,
	//	},
	//	Body: NewAbstractCommand(func(ctx *Context) *Value {
	//		parameter := ctx.FindParameter("value")
	//		asInt, isInt := parameter.Value.(int64)
	//		if isInt {
	//			result := ctx.receiver.Value.(float64) + float64(asInt)
	//			return &Value{
	//				Type:  FloatType,
	//				Value: result,
	//			}
	//		} else {
	//			asFloat, isFloat := parameter.Value.(float64)
	//			if isFloat {
	//				result := ctx.receiver.Value.(float64) + asFloat
	//				return &Value{
	//					Type:  FloatType,
	//					Value: result,
	//				}
	//			} else {
	//				//TODO
	//				//While this might work, it ignores the fact that values won't be "cast" if passed. An Int passed as Any will still try and use Int functions
	//				result := util.Stringify(ctx.receiver.Value) + util.Stringify(parameter.Value)
	//				return &Value{
	//					Type:  StringType,
	//					Value: result,
	//				}
	//			}
	//		}
	//	}),
	//}
	//
	//IntType.variables = convert(map[String]Function{
	//	"plus": intAdd,
	//	"add":  intAdd,
	//	"minus": {
	//		Signature: Signature{
	//			Parameters: []Parameter{
	//				{
	//					Name: "value",
	//					Type: *IntType,
	//				},
	//			},
	//			ReturnType: *IntType,
	//		},
	//		Body: NewAbstractCommand(func(ctx *Context) *Value {
	//			parameter := ctx.FindParameter("value")
	//			result := ctx.receiver.Value.(int64) - parameter.Value.(int64)
	//			return &Value{
	//				Type:  IntType,
	//				Value: result,
	//			}
	//		}),
	//	},
	//	"times": {
	//		Signature: Signature{
	//			Parameters: []Parameter{
	//				{
	//					Name: "value",
	//					Type: *IntType,
	//				},
	//			},
	//			ReturnType: *IntType,
	//		},
	//		Body: NewAbstractCommand(func(ctx *Context) *Value {
	//			parameter := ctx.FindParameter("value")
	//			result := ctx.receiver.Value.(int64) * parameter.Value.(int64)
	//			return &Value{
	//				Type:  IntType,
	//				Value: result,
	//			}
	//		}),
	//	},
	//	"divide": {
	//		Signature: Signature{
	//			Parameters: []Parameter{
	//				{
	//					Name: "value",
	//					Type: *IntType,
	//				},
	//			},
	//			ReturnType: *IntType,
	//		},
	//		Body: NewAbstractCommand(func(ctx *Context) *Value {
	//			parameter := ctx.FindParameter("value")
	//			result := ctx.receiver.Value.(int64) / parameter.Value.(int64)
	//			return &Value{
	//				Type:  IntType,
	//				Value: result,
	//			}
	//		}),
	//	},
	//	"equals": {
	//		Signature: Signature{
	//			Parameters: []Parameter{
	//				{
	//					Name: "value",
	//					Type: *IntType,
	//				},
	//			},
	//			ReturnType: *BooleanType,
	//		},
	//		Body: NewAbstractCommand(func(ctx *Context) *Value {
	//			parameter := ctx.FindParameter("value")
	//			result := ctx.receiver.Value.(int64) == parameter.Value.(int64)
	//			return &Value{
	//				Type:  BooleanType,
	//				Value: result,
	//			}
	//		}),
	//	},
	//})
	//
	//FloatType.variables = convert(map[String]Function{
	//	"plus": floatAdd,
	//	"add":  floatAdd,
	//})
	//
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
	context.DefineVariable(outputWriteName, &Variable{
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
			this := c.FindParameter(0)
			other := c.FindParameter(1)
			return NonReturningValue(BooleanValue(this.Value == other.Value))
		}),
	}
	anyEqualsType := NewFunctionType(anyEquals)
	context.DefineVariable(anyEqualsName, &Variable{
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
