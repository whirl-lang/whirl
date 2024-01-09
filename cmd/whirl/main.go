package main

import (
	lexer "github.com/whirl-lang/whirl/pkg/lexer"
)

func main() {
	lexer.Tokenize("if (x == 1) { print(x); }")
}