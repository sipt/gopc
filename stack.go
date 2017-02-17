package gopc

//StackEmpty : stack is empty
type StackEmpty error

//Stack : string stack
type Stack struct {
	stack []string
}

//NewStringStack : new string stack
func NewStringStack() *Stack {
	return &Stack{
		stack: make([]string, 0),
	}
}

//Pop : stack Pop
func (s *Stack) Pop() (string, *StackEmpty) {
	result, err := s.Peep()
	if err != nil {
		return "", err
	}
	s.stack = s.stack[:len(s.stack)-1]
	return result, nil
}

//Peep : stack Peep not pop
func (s *Stack) Peep() (string, *StackEmpty) {
	len := len(s.stack)
	if len > 0 {
		result := s.stack[len-1]
		return result, nil
	}
	return "", new(StackEmpty)
}

//Push : stack push
func (s *Stack) Push(v string) error {
	s.stack = append(s.stack, v)
	return nil
}

//Clear : Clear stack
func (s *Stack) Clear() {
	s.stack = s.stack[:0]
}
