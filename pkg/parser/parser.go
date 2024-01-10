package parser

import (
	lexer "github.com/whirl-lang/whirl/pkg/lexer"
)

type Instruction interface{}

type Procedure struct {
	Indentifier  string
	Instructions []Instruction
	ReturnType   Type
}

type Struct struct {
	Fields []Field
}

type Field struct {
	Identifier string
	Type       Type
}

type If struct {
	Condition Expression
	Body      []Instruction
}

type Expression interface {
	Value() Literal
}

type Assignment struct {
	Identifier string
	Expression Expression
}

type FunctionCall struct {
	Identifier string
	Arguments  []Expression
}

type Escape struct {
	Expression Expression
}

type Literal interface{}

// Types
type Type interface{}
type Int struct{}
type String struct{}
type Bool struct{}
type Char struct{}

func Iterator(tokens lexer.TokenIterator) []Instruction {
	return nil
}
