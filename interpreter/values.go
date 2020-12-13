package interpreter

func IntValue(int int64) *Value {
	return &Value{
		Type:  IntType,
		Value: int,
	}
}

func FloatValue(num float64) *Value {
	return &Value{
		Type:  FloatType,
		Value: num,
	}
}

func BooleanValue(value bool) *Value {
	return &Value{
		Type:  BooleanType,
		Value: value,
	}
}

func StringValue(value string) *Value {
	return &Value{
		Type:  StringType,
		Value: value,
	}
}
