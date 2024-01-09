package parser

import (
	lexer "github.com/whirl-lang/whirl/pkg/lexer"
)

type Instruction interface {}


func CreateAst(tokens []lexer.Token) []Instruction {
	return nil
}