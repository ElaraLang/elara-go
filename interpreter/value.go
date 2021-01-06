package interpreter

import (
	"github.com/ElaraLang/elara/util"
	"sync"
)

type Value struct {
	Type  Type
	Value interface{}
}

func (v *Value) String() string {
	if v == nil {
		return ""
	}
	return util.Stringify(v.Value)
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

var returnedValues = sync.Pool{
	New: func() interface{} {
		return &ReturnedValue{
			Value:       nil,
			IsReturning: false,
		}
	},
}

type ReturnedValue struct {
	Value       *Value
	IsReturning bool
}

func NewReturningValue(value *Value, returning bool) *ReturnedValue {
	r := returnedValues.Get().(*ReturnedValue)
	r.Value = value
	r.IsReturning = returning
	return r
}
func NonReturningValue(value *Value) *ReturnedValue {
	return NewReturningValue(value, false)
}

func ReturningValue(value *Value) *ReturnedValue {
	return NewReturningValue(value, true)
}

func NilValue() *ReturnedValue {
	return NewReturningValue(nil, false)
}

func (r *ReturnedValue) clean() {
	r.Value = nil
	r.IsReturning = false
	returnedValues.Put(r)
}

func (r ReturnedValue) Unwrap() *Value {
	if r.IsReturning {
		panic("Value should return")
	}
	val := r.Value
	r.clean()
	return val
}
func (r ReturnedValue) UnwrapNotNil() *Value {
	if r.IsReturning {
		panic("Value should return")
	}
	if r.Value == nil {
		panic("Value must not be nil")
	}
	val := r.Value
	r.clean()
	return val
}
