package tests

import (
	"elara/base"
	"elara/interpreter"
	"reflect"
	"testing"
)

func TestSimpleVariableAssignment(t *testing.T) {
	code := `let a = 3
	a`
	results, _, _, _ := base.Execute(nil, code, false)
	expectedResults := []*interpreter.Value{
		nil,
		interpreter.IntValue(3),
	}

	if !reflect.DeepEqual(results, expectedResults) {
		t.Errorf("Incorrect parsing output, got %v but expected %v", formatValues(results), formatValues(expectedResults))
	}
}

func TestSimpleVariableAssignmentWithType(t *testing.T) {
	code := `let a: Int = 3
	a`
	results, _, _, _ := base.Execute(nil, code, false)
	expectedResults := []*interpreter.Value{
		nil,
		interpreter.IntValue(3),
	}

	if !reflect.DeepEqual(results, expectedResults) {
		t.Errorf("Incorrect parsing output, got %v but expected %v", formatValues(results), formatValues(expectedResults))
	}
}

func TestSimpleVariableAssignmentWithTypeAndInvalidValue(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Intepreter allows assignment of incorrect types / values")
		}
	}()

	code := `let a: Int = 3.5
	a`
	base.Execute(nil, code, false)
}

func formatValues(values []*interpreter.Value) string {
	str := "["
	for i, value := range values {
		formatted := value.String()
		if formatted != nil {
			str += *formatted
		} else {
			str += "<nil>"
		}
		if i != len(values)-1 {
			str += ", "
		}
	}
	str += "]"

	return str
}
