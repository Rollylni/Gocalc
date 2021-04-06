package gocalc

import "unicode"

type Lexer struct {
	input  string
	length int
	tokens map[string]TokenType
}

func NewLexer(input string) *Lexer {
	lex := new(Lexer)
	lex.input = input
	lex.length = 0
	lex.tokens = map[string]TokenType{
		"/":  TT_DIVIDE,
		"%":  TT_MODDIV,
		"*":  TT_MULTI,
		"-":  TT_MINUS,
		"^":  TT_EXPON,
		"+":  TT_PLUS,
		"(":  TT_LPAREN,
		")":  TT_RPAREN,
		" ":  TT_IGNORE,
		"\r": TT_IGNORE,
		"\t": TT_IGNORE,
	}
	return lex
}

func ToRune(chr string) rune {
	return []rune(chr)[0]
}

func IsIdentifier(chr rune) bool {
	return unicode.IsLetter(chr) || unicode.IsDigit(chr) || chr == '_'
}

func (lex Lexer) Parse() [](*Token) {
	tokens := [](*Token){}
	offset := 0

	for lex.Has() {
		chr := lex.Consume()
		offset++

		if val, isset := lex.tokens[chr]; isset {
			if val != TT_IGNORE {
				tokens = append(tokens, &Token{val, chr, offset})
			}
		} else if unicode.IsDigit(ToRune(chr)) {
			isfloat := false
			value := chr

			for lex.Has() {
				chr := lex.Look()
				if unicode.IsDigit(ToRune(chr)) {
					lex.Apply()
					value += chr
				} else if chr == "." && !isfloat {
					lex.Apply()
					value += chr
					isfloat = true
				} else {
					break
				}
			}

			tt := TT_INTEGER
			if isfloat {
				tt = TT_FLOAT
			}

			tokens = append(tokens, &Token{tt, value, offset})
		} else if IsIdentifier(ToRune(chr)) {
			value := chr

			for lex.Has() {
				chr := lex.Look()
				if IsIdentifier(ToRune(chr)) {
					lex.Apply()
					value += chr
				} else {
					break
				}
			}
			tokens = append(tokens, &Token{TT_IDENTIFIER, value, offset})
		} else {
			tokens = append(tokens, &Token{TT_UNKNOWN, chr, offset})
		}
	}
	return tokens
}

func (lex Lexer) Has() bool {
	return len(lex.input) != 0
}

func (lex *Lexer) Look(length ...int) string {
	if len(length) <= 0 {
		length = append(length, 1)
	}

	if max := len(lex.input); length[0] > max {
		length[0] = max
	}

	lex.length = length[0]
	return lex.input[:lex.length]
}

func (lex *Lexer) Apply(length ...int) {
	if len(length) <= 0 {
		length = append(length, 0)
	}

	lex.length += length[0]
	lex.input = lex.input[lex.length:]
	lex.length = 0
}

func (lex *Lexer) Consume(length ...int) string {
	if len(length) <= 0 {
		length = append(length, 1)
	}

	res := lex.Look(length[0])
	lex.Apply()
	return res
}
