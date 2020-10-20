package interpreter

import (
	"elara/util"
	"fmt"
)

type Value struct {
	Type  *Type
	Value interface{}
}

func (v *Value) String() *string {
	if v == nil {
		return nil
	}
	formatted := fmt.Sprintf("%s (%s)", util.Stringify(v.Value), v.Type.Name)
	return &formatted
}

var unitValue = &Value{
	Type:  UnitType,
	Value: "Unit",
}

func UnitValue() *Value {
	return unitValue
}

type Variable struct {
	Name    string
	Mutable bool
	Type    Type
	Value   Value
}

func (v Variable) string() string {
	return fmt.Sprintf("Variable { name: %s, mutable: %T, type: %s, Value: %s", v.Name, v.Mutable, v.Type, v.Value)
}
