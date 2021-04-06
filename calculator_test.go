package gocalc

import (
	"fmt"
	"testing"
)

func TestCalculator(t *testing.T) {
	calc := NewCalculator("2 ^ 2 ^ (2 * 2 - 1)")
	res, err := calc.Process()

	if err == nil {
		fmt.Println(res.Value())
	} else {
		fmt.Println(err)
	}
}
