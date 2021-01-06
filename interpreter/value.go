package interpreter

import (
	"fmt"
	"github.com/ElaraLang/elara/util"
)

type Value struct {
	Type  Type
	Value interface{}
}

func (v *Value) String() *string {
	if v == nil {
		return nil
	}
	formatted := fmt.Sprintf("%s (%s)", util.Stringify(v.Value), v.Type.Name())
	return &formatted
}

var unitValue = &Value{
	Type:  UnitType,
	Value: "Unit",
}

func UnitValue() *Value {
	return unitValue
}
