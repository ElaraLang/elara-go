package interpreter

import "fmt"

type Variable struct {
	Name    string
	Mutable bool
	Type    Type
	Value   *Value
}

func (v Variable) string() string {
	return fmt.Sprintf("Variable { Name: %s, mutable: %T, type: %s, Value: %s", v.Name, v.Mutable, v.Type, v.Value)
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
