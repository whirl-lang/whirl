package lexer

import (
	"bytes"
)

type TokenIterator struct {
	Bytes     []byte
	NextToken *Token
}

func (iter *TokenIterator) Peek() (Token, error) {
	if iter.NextToken == nil {
		token, err := iter.Next()

		if err != nil {
			return Token{}, err
		}

		iter.NextToken = &token
	}

	return *iter.NextToken, nil
}

func (iter *TokenIterator) Next() (Token, error) {
	if iter.NextToken != nil {
		token := *iter.NextToken
		iter.NextToken = nil

		return token, nil
	}

	if len(iter.Bytes) == 0 {
		return Token{EOF, ""}, nil
	}

	for iter.Bytes[0] == ' ' || iter.Bytes[0] == '\n' || iter.Bytes[0] == '\t' || iter.Bytes[0] == '\r' {
		iter.Bytes = iter.Bytes[1:]

		if len(iter.Bytes) == 0 {
			return Token{EOF, ""}, nil
		}
	}

	// check for keywords
	//keywords must have a space, tab or newline after them
	for i := IF; i <= STRUCT; i++ {
		word := TokensWithSpace[i]

		if iter.FoundToken(word, true) {

			return Token{i, string(word)}, nil
		}

	}

	//check for the rest
	for i := EQ; i <= STRING; i++ {
		word := TokensWithoutSpace[i]

		if iter.FoundToken(word, false) {
			return Token{i, string(word)}, nil
		}
	}

	//Check for booleans
	if iter.FoundToken([]byte("false"), true) {
		return Token{BOOLEAN_LIT, "false"}, nil
	}

	if iter.FoundToken([]byte("true"), true) {
		return Token{BOOLEAN_LIT, "true"}, nil
	}

	//Check for strings
	if iter.Bytes[0] == '"' {
		iter.Bytes = iter.Bytes[1:]
		var str []byte

		for iter.Bytes[0] != '"' {
			str = append(str, iter.Bytes[0])
			iter.Bytes = iter.Bytes[1:]
		}

		iter.Bytes = iter.Bytes[1:]

		return Token{STRING_LIT, "\"" + string(str) + "\""}, nil
	}

	//check for char
	if iter.Bytes[0] == '\'' {
		if len(iter.Bytes) >= 3 && iter.Bytes[2] == '\'' {
			char := iter.Bytes[1]
			iter.Bytes = iter.Bytes[3:]

			return Token{CHAR_LIT, string(char)}, nil
		}
	}

	//check for int
	if iter.Bytes[0] >= '0' && iter.Bytes[0] <= '9' {
		length := 0

		for iter.Bytes[length] >= '0' && iter.Bytes[length] <= '9' {
			length++
		}

		num := iter.Bytes[:length]
		iter.Bytes = iter.Bytes[length:]
		return Token{INT_LIT, string(num)}, nil
	}

	//check for identifier
	index := 0

	for (iter.Bytes[index] >= 'a' && iter.Bytes[index] <= 'z') || (iter.Bytes[index] >= 'A' && iter.Bytes[index] <= 'Z') || iter.Bytes[index] == '_' {
		index++
	}
	str := iter.Bytes[:index]
	iter.Bytes = iter.Bytes[index:]
	return Token{IDENT, string(str)}, nil

}

func (iter *TokenIterator) FoundToken(token []byte, seperation bool) bool {
	length := len(token)
	hasPrefix := bytes.HasPrefix(iter.Bytes, token)

	if seperation && len(iter.Bytes) > length && !IsSeperationByte(iter.Bytes[length]) {
		return false
	}

	if hasPrefix && len(iter.Bytes) >= length {
		iter.Bytes = iter.Bytes[length:]
	}

	return hasPrefix

}

func IsSeperationByte(b byte) bool {
	return b == ' ' || b == '\n' || b == '\t' || b == '\r' || b == '*' || b == '/' || b == '%' || b == '+' || b == '-' || b == '(' || b == ')' || b == '{' || b == '}' || b == ',' || b == ';' || b == ':' || b == '=' || b == '!' || b == '<' || b == '>'
}

var TokensWithSpace = [][]byte{
	IF:       []byte("if"),
	LET:      []byte("let"),
	UNTIL:    []byte("until"),
	ITER:     []byte("iter"),
	ELSE:     []byte("else"),
	PROC:     []byte("proc"),
	CONTINUE: []byte("continue"),
	ESCAPE:   []byte("escape"),
	BREAK:    []byte("break"),
	STRUCT:   []byte("struct"),
}

var TokensWithoutSpace = [][]byte{
	EQ:  []byte("=="),
	NE:  []byte("!="),
	LT:  []byte("<"),
	GT:  []byte(">"),
	LE:  []byte("<="),
	GE:  []byte(">="),
	AND: []byte("&&"),
	OR:  []byte("||"),
	NOT: []byte("!"),

	COLON:     []byte(":"),
	COMMA:     []byte(","),
	SEMICOLON: []byte(";"),
	ASSIGN:    []byte("="),

	PARENOPEN:    []byte("("),
	PARENCLOSE:   []byte(")"),
	CURLYOPEN:    []byte("{"),
	CURLYCLOSE:   []byte("}"),
	BRACKETOPEN:  []byte("["),
	BRACKETCLOSE: []byte("]"),

	PLUS:  []byte("+"),
	MINUS: []byte("-"),
	MUL:   []byte("*"),
	DIV:   []byte("/"),
	MOD:   []byte("%"),

	BOOLEAN: []byte("bool"),
	CHAR:    []byte("char"),
	VOID:    []byte("void"),
	INT:     []byte("int"),
	STRING:  []byte("string"),
}

type Token struct {
	Kind  int
	Value string
}

func (t Token) IsSeparator() bool {
	return (t.Kind >= COLON && t.Kind <= BRACKETCLOSE) || (t.Kind >= IF && t.Kind <= STRUCT)
}

var TokensPretty = []string{
	IF:       "if",
	LET:      "let",
	UNTIL:    "until",
	ITER:     "iter",
	ELSE:     "else",
	PROC:     "proc",
	CONTINUE: "continue",
	ESCAPE:   "escape",
	BREAK:    "break",
	STRUCT:   "struct",

	EQ:  "==",
	NE:  "!=",
	LT:  "<",
	GT:  ">",
	LE:  "<=",
	GE:  ">=",
	AND: "&&",
	OR:  "||",
	NOT: "!",

	COLON:     ":",
	COMMA:     ",",
	SEMICOLON: ";",
	ASSIGN:    "=",

	PARENOPEN:    "(",
	PARENCLOSE:   ")",
	CURLYOPEN:    "{",
	CURLYCLOSE:   "}",
	BRACKETOPEN:  "[",
	BRACKETCLOSE: "]",

	PLUS:  "+",
	MINUS: "-",
	MUL:   "*",
	DIV:   "/",
	MOD:   "%",

	BOOLEAN: "bool",
	CHAR:    "char",
	VOID:    "void",
	INT:     "int",
	STRING:  "string",

	IDENT:       "<identifier>",
	INT_LIT:     "<integer>",
	STRING_LIT:  "<string>",
	BOOLEAN_LIT: "<boolean>",
	CHAR_LIT:    "<character>",

	EOF: "EOF",
}

const (
	//the order matters here

	//Keywords
	IF = iota
	LET
	UNTIL
	ITER
	ELSE
	PROC
	CONTINUE
	ESCAPE
	BREAK
	STRUCT

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

	COLON
	COMMA
	SEMICOLON
	ASSIGN // this is = while EQ is ==

	//Delimeter
	PARENOPEN
	PARENCLOSE
	CURLYOPEN
	CURLYCLOSE
	BRACKETOPEN
	BRACKETCLOSE

	//Arithmetic
	PLUS
	MINUS
	MUL
	DIV
	MOD

	//Types
	BOOLEAN
	CHAR
	VOID
	INT
	STRING

	//anything after this line will not be checked for keywords

	//Literals
	IDENT
	INT_LIT
	STRING_LIT
	BOOLEAN_LIT
	CHAR_LIT

	EOF
)

func Iterator(input []byte) TokenIterator {
	return TokenIterator{input, nil}
}
