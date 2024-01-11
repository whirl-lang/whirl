package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/whirl-lang/whirl/pkg/codegen"
	"github.com/whirl-lang/whirl/pkg/lexer"
	"github.com/whirl-lang/whirl/pkg/parser"
)

type Args struct {
	args []string
}

func main() {

	args := os.Args

	if len(args) < 2 {
		fmt.Println("Usage: whirl <filename> [args]")
		return
	}

	filename := args[1]
	nodes := ParseFile(filename)
	file := GenerateCode(nodes)
	out := ExecuteFile(filename)
	ParseArgs(Args{args: args}, file)

	fmt.Println(string(out))
}

func ParseArgs(args Args, file *os.File) {
	//TODO: Add proper arg parsing and flags
	deletingFile := true
	if len(args.args) > 2 {
		for _, arg := range args.args[2:] {
			switch arg {
			case "-c":
				deletingFile = false
			}
		}
	}
	file.Close()
	if deletingFile {

		err := os.Remove("out.c")
		if err != nil {
			panic(err)
		}
	}
}

func ExecuteFile(filename string) []byte {
	out, err := exec.Command("tcc", "-run", "out.c").Output()

	if err != nil {
		panic(err)
	}

	return out
}

func ParseFile(filename string) parser.InstructionIterator {

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

	return nodes
}

func GenerateCode(nodes parser.InstructionIterator) *os.File {

	file, err := os.Create("out.c")

	if err != nil {
		panic(err)
	}

	codegen.Generate(&nodes, file)

	return file
}
