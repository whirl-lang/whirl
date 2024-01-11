package codegen

import (
	"bufio"
	"io"
)

type InstructionIterator interface {
	Next() (Instruction, error)
}

type Context struct {
	Namespace string
	Transpile func(content []byte, namespace string) string
}

func WriteC(ctx Context, nodes InstructionIterator, out io.Writer) {
	writer := bufio.NewWriter(out)

	defer writer.Flush()

	for {
		node, err := nodes.Next()

		if err != nil {
			panic(err)
		}

		if node == nil {
			break
		}

		writer.WriteString(node.CInstruction(ctx))
	}
}
