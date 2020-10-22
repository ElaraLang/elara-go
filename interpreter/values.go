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
