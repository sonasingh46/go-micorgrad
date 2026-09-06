package value

// Value can be imagined a node in graph
// which is created due to a mathematical operation of
// one or more than one Values
// Example : Value1 = Value2 + Value3
type Value struct {
	Data      int
	Grad      int
	Operation string
	Prev      []*Value
	Label     string
}

func New(data int, label string) *Value {
	return &Value{Data: data, Label: label, Prev: []*Value{}}
}

func (V *Value) Add(a *Value, label string) *Value {
	data := V.Data + a.Data
	// Add the operation string for visualization
	Operation := "+"
	// Add the operands of the addition as their child nodes
	operandA := a
	operandB := V
	children := []*Value{operandA, operandB}
	return &Value{Data: data, Operation: Operation, Label: label, Prev: children}
}

func (V *Value) Mult(a *Value, label string) *Value {
	data := V.Data * a.Data
	// Add the operation string for visualization
	Operation := "*"
	// Add the operands of the addition as their child nodes
	operandA := a
	operandB := V
	children := []*Value{operandA, operandB}
	return &Value{Data: data, Operation: Operation, Label: label, Prev: children}
}

func (V *Value) Sub(a *Value, label string) *Value {
	data := V.Data - a.Data
	Operation := "-"
	children := []*Value{V, a}
	return &Value{Data: data, Operation: Operation, Label: label, Prev: children}
}

func (V *Value) Div(a *Value, label string) *Value {
	data := V.Data / a.Data
	Operation := "/"
	children := []*Value{V, a}
	return &Value{Data: data, Operation: Operation, Label: label, Prev: children}
}
