package typesystem

import "elara/interpreter"

var AnyType = SimpleType("Any", []interpreter.Function{})
var UnitType = SimpleType("Unit", []interpreter.Function{})

var IntType = SimpleType("Int", []interpreter.Function{})
var FloatType = SimpleType("Float", []interpreter.Function{})
var StringType = SimpleType("String", []interpreter.Function{})
