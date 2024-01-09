package main

import (
	"fmt"
	"os"

	lexer "github.com/whirl-lang/whirl/pkg/lexer"
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

	t, err := tokens.Next()
	
	if err != nil {
		panic(err)
	}
	
	for t.Kind != lexer.UNKNOWN {
		fmt.Println(t)
		t, err = tokens.Next()

		if err != nil {
			panic(err)
		}
		
	}
	fmt.Println(t)


}