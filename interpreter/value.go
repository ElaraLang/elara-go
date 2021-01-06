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
func (r ReturnedValue) UnwrapNotNil() *Value {
	if r.IsReturning {
		panic("Value should return")
	}
	if r.Value == nil {
		panic("Value must not be nil")
	}
	return r.Value
}
func (v *Value) String() string {
	if v == nil {
		return ""
	}
	formatted := fmt.Sprintf("%s (%s)", util.Stringify(v.Value), v.Type.Name())
	return formatted
}

func (v *Value) Copy() *Value {
	if v == nil {
		return nil
	}
	return NewValue(v.Type, v.Value)
}

var unitValue = NewValue(UnitType, nil)

func UnitValue() *Value {
	return unitValue
}
