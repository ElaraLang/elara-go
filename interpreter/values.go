package interpreter

import (
	"sync"
)

var valuePool = sync.Pool{
	New: func() interface{} {
		return &Value{
			Type:  nil,
			Value: nil,
		}
	},
}

func NewValue(valueType Type, value interface{}) *Value {
	v := valuePool.Get().(*Value)
	v.Type = valueType
	v.Value = value

	return v
}

func (v *Value) Cleanup() {
	v.Value = nil
	v.Type = nil
	valuePool.Put(v)
}

func IntValue(int int64) *Value {
	return NewValue(IntType, int)
}

func FloatValue(num float64) *Value {
	return NewValue(FloatType, num)
}

func BooleanValue(value bool) *Value {
	return NewValue(BooleanType, value)
}

func StringValue(value string) *Value {
	return NewValue(StringType, value)
}
