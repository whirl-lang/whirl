package codegen

import (
	"os"

	"github.com/whirl-lang/whirl/pkg/parser"
)

func Generate(nodes parser.InstructionIterator) {

	output := []byte{}

	for {
		node, err := nodes.Next()

		if err != nil {
			panic(err)
		}

		if node == nil {
			break
		}

		output = append(output, generateCode(node)...)

	}

	file, err := os.Create("out.c")

	if err != nil {
		panic(err)
	}

	defer file.Close()

	file.Write(output)

}

func generateCode(node parser.Instruction) []byte {
	out := []byte{}
	switch node := node.(type) {
	case parser.Procedure:
		out = append(out, []byte(node.ReturnType.CType())...)
		out = append(out, []byte(" ")...)
		out = append(out, []byte(node.Indentifier)...)
		out = append(out, []byte("(")...)
		//args

		for i, arg := range node.Arguments {
			out = append(out, []byte(arg.Type.CType())...)
			out = append(out, []byte(" ")...)
			out = append(out, []byte(arg.Identifier)...)
			if i < len(node.Arguments)-1 {
				out = append(out, []byte(", ")...)
			}
		}

		out = append(out, []byte(")")...)
		out = append(out, []byte(" {\n")...)
		//body

		for _, instruction := range node.Instructions {
			out = append(out, generateCode(instruction)...)
		}

		out = append(out, []byte("}\n\n")...)
	case parser.Escape:
		out = append(out, []byte("\treturn")...)
		out = append(out, []byte(" ")...)
		out = append(out, []byte("0")...)
		out = append(out, []byte(";\n")...)
	}

	return out

	/*
		returntype name (args) {
			body

			return value (if not void)
		}

	*/
}
