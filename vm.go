package ps

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
)

type Token interface{}

type VM struct {
	data []Token
}

func NewVM() *VM {
	vm := new(VM)
	vm.data = make([]Token, 0, 100)
	return vm
}

func (vm *VM) push(operand Token) {
	fmt.Println(operand)
	vm.data = append(vm.data, operand)
}

func (vm *VM) pop() (Token, error) {
	if len(vm.data) == 0 {
		return 0, ErrStackUnderFlow
	}

	lastIndex := len(vm.data) - 1
	operand := vm.data[lastIndex]
	vm.data = vm.data[0:lastIndex]
	return operand, nil
}

func (vm *VM) Peek() (Token, error) {
	if len(vm.data) == 0 {
		return 0, ErrStackUnderFlow
	}

	lastIndex := len(vm.data) - 1
	return vm.data[lastIndex], nil
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
		token, err := vm.parsePrimitive(reader)
		if err != nil {
			return err
		}
		vm.push(token)
		return nil
	}

	if r == '(' {
		token, err := vm.parseString(reader)
		if err != nil {
			return err
		}
		vm.push(token)
		return nil
	}

	return ErrInvalidToken
}

func (vm *VM) parsePrimitive(reader *Reader) (Token, error) {
	var buf bytes.Buffer
	for {
		r, err := reader.readRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if !isLetterDigitSymbol(r) {
			reader.unreadRune()
			break
		}
		buf.WriteRune(r)
	}

	s := buf.String()

	if s == "null" {
		return nil, nil
	}
	if s == "true" {
		return true, nil
	}
	if s == "false" {
		return false, nil
	}
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return i, nil
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f, nil
	}

	// TODO: 評価
	return s, nil
}

func (vm *VM) parseString(reader *Reader) (Token, error) {
	var buf bytes.Buffer
	for {
		r, err := reader.readRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if r == ')' {
			break
		}
		buf.WriteRune(r)
	}
	return buf.String(), nil
}
