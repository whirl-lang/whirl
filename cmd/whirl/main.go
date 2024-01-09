package main

import (
	lexer "github.com/whirl-lang/whirl/pkg/lexer"
	"os"
)

func main() {

	args := os.Args

	if len(args) != 2 {
		panic("Invalid number of arguments")
	}

	filename := args[1]

	file, err := os.Open(filename)
	check(err)
	
	stat, err := file.Stat()
	check(err)

	bytes = make([]byte, stat.Size())
	count, err := file.Read(bytes)
	check(err)

	lexer.Tokenize(bytes)
}