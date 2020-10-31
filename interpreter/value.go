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
	Value   *Value
}

func (v Variable) string() string {
	return fmt.Sprintf("Variable { name: %s, mutable: %T, type: %s, Value: %s", v.Name, v.Mutable, v.Type, v.Value)
}

func (v *Variable) Equals(other Variable) bool {
	if v.Name != other.Name {
		return false
	}
	if v.Mutable != other.Mutable {
		return false
	}
	if !v.Type.Accepts(other.Type) {
		//TODO exact equality necessary?
		return false
	}
	if v.Value != other.Value {
		return false
	}
	return true
}
