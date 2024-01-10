package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/whirl-lang/whirl/pkg/codegen"
	"github.com/whirl-lang/whirl/pkg/lexer"
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
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}

	bytes := make([]byte, stat.Size())
	file.Read(bytes)

	tokens := lexer.Iterator(bytes)
	nodes := parser.Iterator(tokens)

	codegen.Generate(nodes)

	out, err := exec.Command("tcc", "-run", "out.c").Output()

	if err != nil {
		panic(err)
	}

	fmt.Println(string(out))
}
