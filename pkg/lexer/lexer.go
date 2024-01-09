package lexer

import (
	"bytes"
	"errors"
	"fmt"
)

type TokenIterator struct {
	bytes []byte
}

func (iter *TokenIterator) Next() (Token, error) {
	fmt.Println(string(iter.bytes))
	if len(iter.bytes) == 0 {
		
		return Token{UNKNOWN, ""}, nil
	}

	//skip whitespace
	for iter.bytes[0] == ' ' || iter.bytes[0] == '\n' || iter.bytes[0] == '\t' {
		iter.bytes = iter.bytes[1:]
		
		if len(iter.bytes) == 0 {
			
			return Token{UNKNOWN, ""}, nil
		}
	}

	// check for keywords
	//keywords must have a space, tab or newline after them
	for i := IF; i <= STRUCT; i++ {
		word := TokensWithSpace[i]
		
		if bytes.HasPrefix(iter.bytes, word) {
			length := len(word)

			if length >= len(iter.bytes) || iter.bytes[length] == ' ' || iter.bytes[length] == '\n' || iter.bytes[length] == '\t' {
				iter.bytes = iter.bytes[length:]
				
				return Token{i, ""}, nil
			}
		}
	}

	//check for the rest
	for i := EQ; i <= STRING; i++ {
		word := TokensWithoutSpace[i]

		if bytes.HasPrefix(iter.bytes, word) {
			iter.bytes = iter.bytes[len(word):]
			
			
			return Token{i, ""}, nil
		}
	}

	//check for identifiers and literals
	

	
	
	return Token{}, errors.New("unknown token")
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

	PARENOPEN:  []byte("("),
	PARENCLOSE: []byte(")"),
	CURLYOPEN:  []byte("{"),
	CURLYCLOSE: []byte("}"),

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

	UNKNOWN
	
)

func Iterator(input []byte) TokenIterator {
	return TokenIterator{input}
}