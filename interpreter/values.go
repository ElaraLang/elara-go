package interpreter

func NewValue(valueType Type, value interface{}) *Value {
	return &Value{
		Type:  valueType,
		Value: value,
	}
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
