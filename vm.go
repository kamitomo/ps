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

	// バインドされている値がオペーレータかどうか判定
	operator, isOperator := value.(Operator)
	if isOperator {
		// オペレータ実行
		return operator.Execute(vm)
	}

	// キーにバインドされた値をプッシュ
	vm.push(value)

	return nil
}

// misc pop

func (vm *VM) popFloat() (float64, error) {
	operand, err := vm.pop()
	if err != nil {
		return 0, err
	}
	f, ok := operand.(float64)
	if !ok {
		err = fmt.Errorf("%v は float64 型と互換性がない: %w", operand, ErrType)
		return 0, err
	}
	return f, nil
}

func (vm *VM) popInt() (int, error) {
	f, err := vm.popFloat()
	if err != nil {
		return 0, err
	}
	return int(f), nil
}

func (vm *VM) popOperator() (Operator, error) {
	operator, err := vm.pop()
	if err != nil {
		return nil, err
	}
	return operator.(Operator), nil
}

func (vm *VM) popString() (string, error) {
	operand, err := vm.pop()
	if err != nil {
		return "", err
	}
	str, ok := operand.(string)
	if !ok {
		err = fmt.Errorf("%v は string 型と互換性がない: %w", operand, ErrType)
		return "", err
	}
	return str, nil
}

func (vm *VM) popBoolean() (bool, error) {
	operand, err := vm.pop()
	if err != nil {
		return false, err
	}
	b, ok := operand.(bool)
	if !ok {
		err = fmt.Errorf("%v は bool 型と互換性がない: %w", operand, ErrType)
		return false, err
	}
	return b, nil
}
