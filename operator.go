package ps

import "fmt"

type Operator interface {
	Execute(vm *VM) error
}

type OperatorFunc func(vm *VM) error

type PrimitiveOperator struct {
	f OperatorFunc
}

func NewOperator(f OperatorFunc) *PrimitiveOperator {
	return &PrimitiveOperator{f}
}

func (o *PrimitiveOperator) Execute(vm *VM) error {
	return o.f(vm)
}

// stack

func stack(vm *VM) error {
	fmt.Println(vm)
	return nil
}

// math

//num1  num2 add sum -> Return num1 plus num2
func add(vm *VM) error {
	num2, err := vm.popFloat()
	if err != nil {
		return err
	}
	num1, err := vm.popFloat()
	if err != nil {
		return err
	}
	vm.push(num1 + num2)
	return nil
}
