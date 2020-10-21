package tests

import (
	"elara/base"
	"elara/interpreter"
	"reflect"
	"testing"
)

func TestSimpleVariableDeclaration(t *testing.T) {
	code := `let a = 3
	a`
	results := base.Execute(nil, code, false)
	expectedResults := []*interpreter.Value{
		nil,
		interpreter.IntValue(3),
	}

	if !reflect.DeepEqual(results, expectedResults) {
		t.Fail()
		t.Errorf("Incorrect lexing output, got %v but expected %v", formatValues(results), formatValues(expectedResults))
	}
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
