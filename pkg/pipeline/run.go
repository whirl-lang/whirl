package pipeline

import (
	"bytes"

	"github.com/whirl-lang/whirl/pkg/codegen"
	"github.com/whirl-lang/whirl/pkg/lexer"
	"github.com/whirl-lang/whirl/pkg/parser"
)

func transpile(content []byte) parser.InstructionIterator {
	tokens := lexer.Iterator(content)
	nodes := parser.Iterator(tokens)

	return nodes
}

// Transpiles the given Whirl source code into C source code.
func TranspileC(content []byte, namespace string) string {
	nodes := transpile(content)
	buf := bytes.NewBufferString("")

	codegen.WriteC(codegen.Context{
		Namespace: namespace,
		Transpile: TranspileC,
	}, &nodes, buf)

	return buf.String()
}
