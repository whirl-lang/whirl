package lexer

type Token struct {
	kind  int
	value string
}

const (
	EOF = iota

	COLON
	COMMA
	SEMICOLON
	ASSIGN // this is = while EQ is ==
	IDENT

	//Keywords
	IF
	LET
	UNTIL
	ITER
	ELSE
	PROC
	CONTINUE
	ESCAPE
	BREAK

	//Delimeter
	PARENOPEN
	PARENCLOSE
	CURLYOPEN
	CURLYCLOSE

	//Arithmetic
	PLUS
	MINUS
	MUL
	DIV
	MOD

	//Operators
	EQ
	NE
	LT
	GT
	LE
	GE
	AND
	OR
	NOT

	//Literals
	NUMBER
	STRING
	BOOLEAN
	CHAR
)

func Tokenize(input string) []Token {
	return []Token{ Token{ kind: EOF, value: "" } }
}