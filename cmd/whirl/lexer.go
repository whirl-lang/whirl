package main

import "bufio"

type Token int

const (
	EOF = iota
	ILLEGAL
	IDENT
	INT
	SEMICOLON
	COLON

	// Operators
	ADD
	SUB
	MUL
	DIV

	ASSIGN
)

var tokens = []string{
	EOF:       "EOF",
	ILLEGAL:   "ILLEGAL",
	IDENT:     "IDENT",
	INT:       "INT",
	COLON:     ":",
	SEMICOLON: ";",

	// Operators
	ADD: "+",
	SUB: "-",
	MUL: "*",
	DIV: "/",

	ASSIGN: "=",
}

func (t Token) String() string {
	return tokens[t]
}

type Position struct {
	line   int
	column int
}

type Lexer struct {
	pos    Position
	reader *bufio.Reader
}

