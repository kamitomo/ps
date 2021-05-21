package ps

import "fmt"

type Operator interface {
	Execute(vm *VM)
}

type OperatorFunc func(vm *VM)

type PrimitiveOperator struct {
	f OperatorFunc
}

func NewOperator(f OperatorFunc) *PrimitiveOperator {
	return &PrimitiveOperator{f}
}

func (o *PrimitiveOperator) Execute(vm *VM) {
	o.f(vm)
}

// stack

func stack(vm *VM) {
	fmt.Println(vm)
}

// math

//num1  num2 add sum -> Return num1 plus num2
func add(vm *VM) {
	num2 := vm.popFloat()
	num1 := vm.popFloat()
	vm.push(num1 + num2)
}
