package gocalc

import "testing"

func TestToken(t *testing.T) {
	etype := TT_DIVIDE
	evalue := "/"
	tok := new(Token)
	tok.SetType(etype)
	tok.SetValue(evalue)

	if tok.Value() != evalue {
		t.Error("Expected /, got ", tok.Value())
	} else if tok.Type() != etype {
		t.Error("Expected 4, got", etype)
	}
}
