package interpreter

import (
	"fmt"
	"github.com/ElaraLang/elara/util"
)

type Value struct {
	Type  Type
	Value interface{}
}

type ReturnedValue struct {
	Value       *Value
	IsReturning bool
}

func NonReturningValue(value *Value) ReturnedValue {
	return ReturnedValue{
		Value:       value,
		IsReturning: false,
	}
}

func ReturningValue(value *Value) ReturnedValue {
	return ReturnedValue{
		Value:       value,
		IsReturning: true,
	}
}

func NilValue() ReturnedValue {
	return ReturnedValue{
		Value:       nil,
		IsReturning: false,
	}
}

func (r ReturnedValue) Unwrap() *Value {
	if r.IsReturning {
		panic("Value should return")
	}
	return r.Value
}
func (v *Value) String() *string {
	if v == nil {
		return nil
	}
	formatted := fmt.Sprintf("%s (%s)", util.Stringify(v.Value), v.Type.Name())
	return &formatted
}

var unitValue = &Value{
	Type:  UnitType,
	Value: "Unit",
}

func UnitValue() *Value {
	return unitValue
}
