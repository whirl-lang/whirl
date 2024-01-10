package codegen

import (
	"os"
)

type InstructionIterator interface {
	Next() (Instruction, error)
}

func Generate(nodes InstructionIterator) {

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

func generateCode(node Instruction) []byte {
	out := []byte{}
	switch node := node.(type) {
	case Procedure:
		//return type
		out = append(out, []byte(node.ReturnType.CType())...)
		out = append(out, []byte(" ")...)

		//name
		out = append(out, []byte(node.Ident)...)
		out = append(out, []byte("(")...)

		//args
		for i, arg := range node.Args {
			out = append(out, []byte(arg.Type.CType())...)
			out = append(out, []byte(" ")...)
			out = append(out, []byte(arg.Ident)...)
			if i < len(node.Args)-1 {
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
	case Escape:
		out = append(out, []byte("\treturn")...)
		out = append(out, []byte(" ")...)

		out = append(out, GenerateExpr(node.Expr)...)

		out = append(out, []byte(";\n")...)
	}

	return out

}

func GenerateExpr(node Expr) []byte {

	out := []byte{}

	switch node := node.(type) {
	case ExprMath:
		for _, token := range node.Tokens {
			out = append(out, []byte(token.Value)...)
		}
	}

	return out

}
