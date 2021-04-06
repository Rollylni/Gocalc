package gocalc

import (
	"fmt"
	"testing"
)

func TestLexer(t *testing.T) {
	lex := NewLexer("/*+-()%^,12 12.00")
	tokens := lex.Parse()

	for _, tok := range tokens {
		fmt.Println(tok.Type(), ": ", tok.Value())
	}
}
