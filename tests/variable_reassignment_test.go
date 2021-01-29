package tests

import (
	"testing"
)

func TestSimpleVariableReassignment(t *testing.T) {
	/*code := `let mut a = 3
	a = 4
	a`
	results, _, _, _ := base.Execute(nil, code, false)
	expectedResults := []*interpreter.Value{
		nil,
		nil,
		interpreter.IntValue(4),
	}

	if !reflect.DeepEqual(results, expectedResults) {
		t.Errorf("Incorrect parsing output, got %v but expected %v", formatValues(results), formatValues(expectedResults))
	}*/
}

func TestSimpleVariableReassignmentWithType(t *testing.T) {
	/*code := `let mut a: Int = 3
	a = 4
	a`
	results, _, _, _ := base.Execute(nil, code, false)
	expectedResults := []*interpreter.Value{
		nil,
		nil,
		interpreter.IntValue(4),
	}

	if !reflect.DeepEqual(results, expectedResults) {
		t.Errorf("Incorrect parsing output, got %v but expected %v", formatValues(results), formatValues(expectedResults))
	}*/
}

func TestSimpleVariableReassignmentWithTypeAndInvalidValue(t *testing.T) {
	/*defer func() {
		if r := recover(); r == nil {
			t.Errorf("Intepreter allows reassignment of incorrect types / values")
		}
	}()

	code := `let mut a: Int = 3
	a = 3.5`
	base.Execute(nil, code, false)*/
}

func TestSimpleVariableReassignmentWithImmutableVariable(t *testing.T) {
	/*defer func() {
		if r := recover(); r == nil {
			t.Errorf("Intepreter allows reassignment of immutable variable")
		}
	}()

	code := `let a = 3
	a = 4`
	results, _, _, _ := base.Execute(nil, code, false)

	expectedResults := []*interpreter.Value{
		nil,
		interpreter.IntValue(4),
	}
	if !reflect.DeepEqual(results, expectedResults) {
		t.Errorf("Incorrect parsing output, got %v but expected %v", formatValues(results), formatValues(expectedResults))
	}*/
}
