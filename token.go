package gocalc

type TokenType int

const (
	TT_UNKNOWN TokenType = iota
	TT_IDENTIFIER
	TT_INTEGER
	TT_FLOAT
	TT_DIVIDE
	TT_MODDIV
	TT_MULTI
	TT_MINUS
	TT_EXPON
	TT_PLUS
	TT_LPAREN
	TT_RPAREN
	TT_IGNORE
)

type Token struct {
	ttype TokenType
	value string
	pos   int
}

func (tok Token) Type() TokenType {
	return tok.ttype
}

func (tok *Token) SetType(val TokenType) {
	tok.ttype = val
}

func (tok Token) Value() string {
	return tok.value
}

func (tok *Token) SetValue(val string) {
	tok.value = val
}

func (tok Token) Pos() int {
	return tok.pos
}
