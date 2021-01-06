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

var types = map[string]Type{
	"Any":  AnyType,
	"Unit": UnitType,

	"Int":     IntType,
	"Float":   FloatType,
	"Boolean": BooleanType,
	"String":  StringType,
	"Output":  OutputType,
}

func Init(context *Context) {
	for s, t := range types {
		context.types[s] = t
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
					Name: "other",
					Type: AnyType,
				}},
			ReturnType: StringType,
		},
		Body: NewAbstractCommand(func(ctx *Context) *Value {
			this := ctx.FindParameter("this")
			otherParam := ctx.FindParameter("other")
			concatenated := this.Value.(string) + util.Stringify(otherParam.Value)
			return &Value{
				Type:  StringType,
				Value: concatenated,
			}
		}),
		name: &stringPlusName,
	}
	stringPlusType := NewFunctionType(stringPlus)
	context.DefineVariable(stringPlusName, Variable{
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
					Name: "other",
					Type: StringType,
				}},
			ReturnType: StringType,
		},
		Body: NewAbstractCommand(func(ctx *Context) *Value {
			this := ctx.FindParameter("this")
			otherParam := ctx.FindParameter("other")
			concatenated := util.Stringify(this.Value) + otherParam.Value.(string)
			return &Value{
				Type:  StringType,
				Value: concatenated,
			}
		}),
		name: &anyPlusName,
	}
	anyPlusType := NewFunctionType(anyPlus)
	context.DefineVariable(anyPlusName, Variable{
		Name:    anyPlusName,
		Mutable: false,
		Type:    anyPlusType,
		Value: &Value{
			Type:  anyPlusType,
			Value: anyPlus,
		},
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
					Name: "other",
					Type: NewCollectionTypeOf(AnyType),
				},
			},
			ReturnType: NewCollectionTypeOf(AnyType),
		},
		Body: NewAbstractCommand(func(ctx *Context) *Value {
			this := ctx.FindParameter("this").Value.(*Collection)
			other := ctx.FindParameter("other").Value.(*Collection)

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
			return &Value{
				Type:  NewCollectionType(newCol),
				Value: newCol,
			}
		}),
		name: &colPlusName,
	}
	colPlusType := NewFunctionType(colPlus)
	context.DefineVariable(colPlusName, Variable{
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
	//			value, err := strconv.ParseInt(ctx.receiver.Value.(string), 10, 64)
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
	//			eq := ctx.receiver.Value.(string) == parameter.Value
	//			return &Value{
	//				Type:  BooleanType,
	//				Value: eq,
	//			}
	//		}),
	//	},
	//})
	//
	//BooleanType.variables = convert(map[string]Function{
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
	//IntType.variables = convert(map[string]Function{
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
	//FloatType.variables = convert(map[string]Function{
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
					Name: "value",
					Type: AnyType,
				},
			},
			ReturnType: UnitType,
		},
		Body: NewAbstractCommand(func(ctx *Context) *Value {
			parameter := ctx.FindParameter("value")
			asString := ctx.Stringify(parameter)
			fmt.Printf("%s", asString)
			return UnitValue()
		}),
		name: &outputWriteName,
	}
	outputWriteType := NewFunctionType(stringPlus)
	context.DefineVariable(outputWriteName, Variable{
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
					Name: "other",
					Type: AnyType,
				},
			},
			ReturnType: BooleanType,
		},
		name: &anyEqualsName,
		Body: NewAbstractCommand(func(c *Context) *Value {
			this := c.FindParameter("this")
			other := c.FindParameter("other")
			return BooleanValue(this.Value == other.Value)
		}),
	}
	anyEqualsType := NewFunctionType(anyEquals)
	context.DefineVariable(anyEqualsName, Variable{
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
