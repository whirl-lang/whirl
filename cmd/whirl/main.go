package main

import (
	"fmt"
	"os"

	lexer "github.com/whirl-lang/whirl/pkg/lexer"
	"github.com/whirl-lang/whirl/pkg/parser"
)

func main() {

	args := os.Args

	if len(args) != 2 {
		panic("Invalid number of arguments")
	}

	filename := args[1]

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}

	bytes := make([]byte, stat.Size())
	file.Read(bytes)

	tokens := lexer.Iterator(bytes)
	nodes := parser.Iterator(tokens)

	for {
		node, err := nodes.Next()

		if err != nil {
			panic(err)
		}

		if node == nil {
			break
		}

		fmt.Println(node)
	}
}
