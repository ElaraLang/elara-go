package interpreter

func IntValue(int int64) *Value {
	return &Value{
		Type:  IntType,
		Value: int,
	}
}
