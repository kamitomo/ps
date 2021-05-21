package ps

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
)

type Token interface{}

type Dictionary map[string]Token

type VM struct {
	data []Token
	dic  Dictionary
}

func NewVM() *VM {
	vm := new(VM)
	vm.data = make([]Token, 0, 100)
	vm.dic = make(Dictionary, 100)
	vm.dic["stack"] = NewOperator(stack)
	vm.dic["add"] = NewOperator(add)
	return vm
}

func (vm *VM) push(operand Token) {
	vm.data = append(vm.data, operand)
}

func (vm *VM) pop() Token {
	lastIndex := len(vm.data) - 1
	operand := vm.data[lastIndex]
	vm.data = vm.data[0:lastIndex]
	return operand
}

func (vm *VM) Peek() (Token, error) {
	if len(vm.data) == 0 {
		return 0, ErrStackUnderFlow
	}

	lastIndex := len(vm.data) - 1
	return vm.data[lastIndex], nil
}

func (vm *VM) String() string {
	var buf string
	for i := len(vm.data) - 1; i >= 0; i-- {
		buf += fmt.Sprintf("%v\n", vm.data[i])
	}
	return buf
}

func (vm *VM) Size() int {
	return len(vm.data)
}

func (vm *VM) Clear() {
	vm.data = vm.data[0:0]
}

func (vm *VM) Execute(f io.Reader) error {
	reader := NewReader(f)

	for {
		err := vm.parse(reader)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}

func (vm *VM) parse(reader *Reader) error {
	reader.SkipWhite()
	r, err := reader.readRune()
	if err != nil {
		return err
	}

	if isLetterDigitSymbol(r) {
		reader.unreadRune()
		return vm.parsePrimitive(reader)
	}

	if r == '(' {
		return vm.parseString(reader)
	}

	return ErrInvalidToken
}

func (vm *VM) parsePrimitive(reader *Reader) error {
	var buf bytes.Buffer
	for {
		r, err := reader.readRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if !isLetterDigitSymbol(r) {
			reader.unreadRune()
			break
		}
		buf.WriteRune(r)
	}

	s := buf.String()

	if s == "null" {
		vm.push(nil)
		return nil
	}
	if s == "true" {
		vm.push(true)
		return nil
	}
	if s == "false" {
		vm.push(false)
		return nil
	}
	// if i, err := strconv.ParseInt(s, 10, 64); err == nil {
	// 	vm.push(i)
	// 	return nil
	// }
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		vm.push(f)
		return nil
	}

	// トークンをキーとして辞書を参照
	return vm.computeReference(s)
}

func (vm *VM) parseString(reader *Reader) error {
	var buf bytes.Buffer
	for {
		r, err := reader.readRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if r == ')' {
			break
		}
		buf.WriteRune(r)
	}
	vm.push(buf.String())
	return nil
}

func (vm *VM) computeReference(key string) error {
	value := vm.dic[key]
	if value == nil {
		return fmt.Errorf("未定義のキーによる辞書参照: %s", key)
	}
	operator, isOperator := value.(Operator)
	if isOperator {
		operator.Execute(vm)
	} else {
		vm.push(value)
	}
	return nil
}

// misc pop

func (vm *VM) popFloat() float64 {
	operand := vm.pop()
	return operand.(float64)
}

func (vm *VM) popInt() int {
	f := vm.popFloat()
	return int(f)
}

func (vm *VM) popOperator() Operator {
	operator := vm.pop()
	return operator.(Operator)
}

func (vm *VM) popName() string {
	name := vm.pop().(string)
	return name[1:]
}

func (vm *VM) popString() string {
	s := vm.pop().(string)
	return s[1 : len(s)-1]
}

func (vm *VM) popBoolean() bool {
	s := vm.pop()
	return s.(bool)
}
