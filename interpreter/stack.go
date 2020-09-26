package interpreter

type Stack struct {
	values []Value
}

func (s *Stack) Append(value Value) {
	s.values = append(s.values, value)
}

func (s *Stack) Peek() Value {
	return s.values[len(s.values)-1]
}

func (s *Stack) Pop() Value {
	n := len(s.values) - 1
	value := s.values[n]
	s.values = s.values[:n]

	return value
}

func NewStack() *Stack {
	return &Stack{
		values: []Value{},
	}
}
