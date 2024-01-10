package codegen

import (
	"bufio"
	"io"
)

type InstructionIterator interface {
	Next() (Instruction, error)
}

func Generate(nodes InstructionIterator, out io.Writer) {
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

		writer.WriteString(node.CInstruction())
	}
}
