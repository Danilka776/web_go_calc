package calc

import (
	"fmt"
	"math"
	"strconv"
)

const length = 100
const epsilon = 1e-6

type stack struct {
	top   int
	items [length]interface{}
}

type stackfloat struct {
	top   int
	items [length]float64
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9' || r == '.'
}

func isOperator(r rune) bool {
	return r == '+' || r == '-' || r == '*' || r == '/'
}

func priority(r rune) int {
	switch r {
	case '+', '-':
		return 1
	case '*', '/':
		return 2
	default:
		return 0
	}
}

func converToPostfix(infix string) stack {
	var postfix stack
	st := make([]rune, 0)
	for i := 0; i < len(infix); i++ {
		if rune(infix[i]) == ' ' {
			continue
		}
		if isDigit(rune(infix[i])) {
			var curNumber string = ""
			for i < len(infix) && isDigit(rune(infix[i])) {
				curNumber += string(rune(infix[i]))
				i++
			}
			i--
			number, err := strconv.ParseFloat(curNumber, 64)
			if err != nil {
				fmt.Println("Wrong cast", err)
			}
			postfix.push(number)
		} else if rune(infix[i]) == '(' {
			st = append(st, rune(infix[i]))
		} else if rune(infix[i]) == ')' {
			for len(st) > 0 && st[len(st)-1] != '(' {
				postfix.push(st[len(st)-1])
				st = st[:len(st)-1]
			}
			st = st[:len(st)-1]
		} else if isOperator(rune(infix[i])) {
			for len(st) > 0 && priority(st[len(st)-1]) >= priority(rune(infix[i])) {
				postfix.push(st[len(st)-1])
				st = st[:len(st)-1]
			}
			st = append(st, rune(infix[i]))
		} else {
			fmt.Printf("Undefined character %c\n", rune(infix[i]))
		}
	}
	for len(st) > 0 {
		postfix.push(st[len(st)-1])
		st = st[:len(st)-1]
	}
	return postfix
}

func oper(symb rune, op1 float64, op2 float64) (float64, error) {
	switch symb {
	case '+':
		return op1 + op2, nil
	case '-':
		return op1 - op2, nil
	case '*':
		return op1 * op2, nil
	case '/':
		if math.Abs(op2) < epsilon {
			return 0.0, fmt.Errorf("Division by zero or too small value")
		}
		return op1 / op2, nil
	default:
		return 0.0, nil
	}
}

func (s *stack) push(x interface{}) {
	s.items[s.top+1] = x
	s.top++
}

func (s *stack) pop() interface{} {
	if s.top == -1 {
		return 0
	}
	s.top--
	return s.items[s.top+1]
}

func (s *stackfloat) push(x float64) {
	s.items[s.top+1] = x
	s.top++
}

func (s *stackfloat) pop() float64 {
	if s.top == -1 {
		return 0
	}
	s.top--
	return s.items[s.top+1]
}

func eva(postfix stack) (float64, error) {
	var opndstk stackfloat
	for i := 1; i <= postfix.top; i++ {
		_, ok := postfix.items[i].(float64)
		if ok {
			opndstk.push(postfix.items[i].(float64))
		} else {
			opnd2 := opndstk.pop()
			opnd1 := opndstk.pop()
			val, error := oper(postfix.items[i].(rune), opnd1, opnd2)
			if error != nil {
				return 0.0, error
			}
			opndstk.push(val)
		}
	}
	return opndstk.pop(), nil
}

func isValidExpression(infix string) error {
	// проверка пустого ввода
	if len(infix) == 0 {
		return fmt.Errorf("Empty expression")
	}

	// Проверка на наличие некорректных символов
	for _, r := range infix {
		if !isDigit(r) && !isOperator(r) && r != '(' && r != ')' && r != ' ' && r != '.' {
			return fmt.Errorf("Wrong character")
		}
	}

	// Проверка скобок
	var brackets int = 0
	for _, r := range infix {
		if r == '(' {
			brackets++
		} else if r == ')' {
			brackets--
			if brackets < 0 {
				return fmt.Errorf("Wrong number of brackets")
			}
		}
	}
	if brackets != 0 {
		return fmt.Errorf("Wrong number of brackets")
	}

	// Провека двойной операции
	var flagOp bool = false
	for _, r := range infix {
		if isOperator(r) {
			if !flagOp {
				flagOp = true
			} else {
				return fmt.Errorf("Double operation")
			}
		} else {
			flagOp = false
		}
	}

	// Проверка корректности выражения (переделать!!!)
	if isOperator(rune(infix[0])) || isOperator(rune(infix[len(infix)-1])) {
		return fmt.Errorf("Wrong expression")
	}

	return nil
}

func Calc(expression string) (float64, error) {
	err := isValidExpression(expression)
	if err != nil {
		return 0, err
	}
	postfix := converToPostfix(expression)
	value, err := eva(postfix)
	if err != nil {
		return 0, err
	}
	return value, nil
}
