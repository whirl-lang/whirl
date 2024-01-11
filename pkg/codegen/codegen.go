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
	Path      string
	Transpile func(content []byte, path string, out io.Writer)
}

func WriteC(ctx Context, nodes InstructionIterator, out io.Writer) error {
	writer := bufio.NewWriter(out)

	defer writer.Flush()

	for {
		node, err := nodes.Next()

		if err != nil {
			return err
		}

		if node == nil {
			break
		}

		writer.WriteString(node.CInstruction(ctx))
	}

	return nil
}
