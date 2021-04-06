package gocalc

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Priority int

var (
	RIGHT Priority = 0
	LEFT  Priority = 1
)

type Calculator struct {
	variables map[string](*Token)
	operators map[TokenType][]Priority
	lexer     *Lexer
	input     string
}

func NewCalculator(input string) *Calculator {
	calc := new(Calculator)
	calc.operators = map[TokenType][]Priority{
		TT_EXPON:  {RIGHT, 3},
		TT_DIVIDE: {LEFT, 2},
		TT_MULTI:  {LEFT, 2},
		TT_MODDIV: {LEFT, 1},
		TT_MINUS:  {LEFT, 1},
		TT_PLUS:   {LEFT, 1},
		TT_LPAREN: {0, 0},
		TT_RPAREN: {0, 0},
	}

	calc.variables = map[string](*Token){
		"E":   MakeFloat(math.E),
		"PI":  MakeFloat(math.Pi),
		"PHI": MakeFloat(math.Phi),
	}
	calc.lexer = NewLexer(strings.ReplaceAll(input, "\n", ""))
	calc.input = input
	return calc
}

func (calc *Calculator) Process() (*Token, error) {
	tokens := calc.lexer.Parse()
	opStack := [](*Token){}
	numStack := [](*Token){}

	for _, val := range tokens {
		if val.Type() == TT_IDENTIFIER {
			if v, isset := calc.variables[val.Value()]; isset {
				numStack = append(numStack, v)
			} else {
				return nil, errors.New("Identifition error: Undefined '" + val.Value() + "' in " + strconv.Itoa(val.Pos()) + "-" + strconv.Itoa(len(val.Value())+val.Pos()-1) + " pos!")
			}
		} else if val.Type() == TT_INTEGER || val.Type() == TT_FLOAT {
			numStack = append(numStack, val)
		} else if val.Type() == TT_RPAREN {
			closed := false
			for len(opStack) != 0 {
				op := Back(opStack)
				opStack = PopBack(opStack)
				if op.Type() == TT_LPAREN {
					closed = true
					break
				}

				if len(numStack) < 2 {
					return nil, errors.New("Operation error: Invalid Syntax!")
				}
				lnum := Back(numStack)
				numStack = PopBack(numStack)
				rnum := Back(numStack)
				numStack = PopBack(numStack)

				if lnum.Type() != rnum.Type() {
					return nil, errors.New("Operation error: Type Error!")
				}

				res := &Token{}
				if lnum.Type() == TT_INTEGER {
					res = ExecuteInt(op, rnum, lnum)
				} else {
					res = ExecuteFloat(op, rnum, lnum)
				}
				numStack = append(numStack, res)
			}

			if !closed {
				return nil, errors.New("Parse error: Unexpeceted '" + val.Value() + "' in " + strconv.Itoa(val.Pos()) + " pos!")
			}
		} else if pty, isset := calc.operators[val.Type()]; isset {
			for {
				bpty := []Priority{0, 0}
				if len(opStack) > 0 {
					bpty, _ = calc.operators[Back(opStack).Type()]
				}

				if len(opStack) == 0 || val.Type() == TT_LPAREN || pty[1] > bpty[1] || (pty[1] == bpty[1] && pty[0] == RIGHT) {
					opStack = append(opStack, val)
					break
				}

				if len(numStack) < 2 {
					return nil, errors.New("Operation error: Invalid Syntax!")
				}

				op := Back(opStack)
				opStack = PopBack(opStack)
				lnum := Back(numStack)
				numStack = PopBack(numStack)
				rnum := Back(numStack)
				numStack = PopBack(numStack)

				if lnum.Type() != rnum.Type() {
					return nil, errors.New("Operation error: Type Error!")
				}

				res := &Token{}
				if lnum.Type() == TT_INTEGER {
					res = ExecuteInt(op, rnum, lnum)
				} else {
					res = ExecuteFloat(op, rnum, lnum)
				}
				numStack = append(numStack, res)
			}
		} else {
			return nil, errors.New("Parse error: Unexpeceted '" + val.Value() + "' in " + strconv.Itoa(val.Pos()) + " pos!")
		}
	}

	for len(opStack) != 0 {
		op := Back(opStack)
		opStack = PopBack(opStack)

		if op.Type() == TT_LPAREN {
			return nil, errors.New("Parse error: unclosed parentheses found!")
		}

		if len(numStack) < 2 {
			return nil, errors.New("Operation error: Invalid Syntax!")
		}

		lnum := Back(numStack)
		numStack = PopBack(numStack)
		rnum := Back(numStack)
		numStack = PopBack(numStack)

		if lnum.Type() != rnum.Type() {
			return nil, errors.New("Operation error: Type Error!")
		}

		res := &Token{}
		if lnum.Type() == TT_INTEGER {
			res = ExecuteInt(op, rnum, lnum)
		} else {
			res = ExecuteFloat(op, rnum, lnum)
		}
		numStack = append(numStack, res)
	}
	return Back(numStack), nil
}

func ExecuteInt(op *Token, num1 *Token, num2 *Token) *Token {
	if num1.Type() != TT_INTEGER || num2.Type() != TT_INTEGER {
		return nil
	}

	n1, _ := ToInt(num1)
	n2, _ := ToInt(num2)
	res := 0

	switch op.Type() {
	case TT_MINUS:
		res = int(n1 - n2)
	case TT_PLUS:
		res = int(n1 + n2)
	case TT_MULTI:
		res = int(n1 * n2)
	case TT_DIVIDE:
		res = int(n1 / n2)
	case TT_EXPON:
		res = int(math.Pow(float64(n1), float64(n2)))
	case TT_MODDIV:
		res = int(n1 % n2)
	}
	return MakeInt(res)
}

func ExecuteFloat(op *Token, num1 *Token, num2 *Token) *Token {
	if num1.Type() != TT_FLOAT || num2.Type() != TT_FLOAT {
		return nil
	}

	n1, _ := ToFloat(num1, 64)
	n2, _ := ToFloat(num2, 64)
	res := 0.0

	switch op.Type() {
	case TT_MINUS:
		res = float64(n1 - n2)
	case TT_PLUS:
		res = float64(n1 + n2)
	case TT_MULTI:
		res = float64(n1 * n2)
	case TT_DIVIDE:
		res = float64(n1 / n2)
	case TT_EXPON:
		res = float64(math.Pow(n1, n2))
	case TT_MODDIV:
		res = float64(int(n1) % int(n2))
	}
	return MakeFloat(res)
}

func MakeFloat(val float64) *Token {
	return &Token{TT_FLOAT, fmt.Sprintf("%f", val), 0}
}

func MakeInt(val int) *Token {
	return &Token{TT_INTEGER, strconv.Itoa(val), 0}
}

func ToFloat(tok *Token, bsize int) (float64, error) {
	return strconv.ParseFloat(tok.Value(), bsize)
}

func ToInt(tok *Token) (int, error) {
	return strconv.Atoi(tok.Value())
}

func Back(arr []*Token) *Token {
	return arr[len(arr)-1]
}

func PopBack(arr []*Token) []*Token {
	return arr[:len(arr)-1]
}
