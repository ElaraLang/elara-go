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

func CharValue(value rune) *Value {
	return NewValue(CharType, value)
}

func StringValue(value string) *Value {
	chars := make([]*Value, len(value))
	for i, c := range value {
		chars[i] = CharValue(c)
	}
	val := &Collection{Elements: chars, ElementType: CharType}
	return NewValue(NewCollectionTypeOf(CharType), val)
}
