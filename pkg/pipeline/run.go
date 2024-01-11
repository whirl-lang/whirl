package pipeline

import (
	"io"

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
func transpileC(content []byte, path string, out io.Writer) {
	nodes := transpile(content)

	codegen.WriteC(codegen.Context{
		Namespace: codegen.PathToNamespace(path),
		Path:      path,
		Transpile: TranspileC,
	}, &nodes, out)
}

func TranspileC(content []byte, path string, out io.Writer) {
	nodes := transpile(content)

	out.Write([]byte("#include <stdio.h>\n\n"))

	codegen.WriteC(codegen.Context{
		Namespace: "",
		Path:      path,
		Transpile: transpileC,
	}, &nodes, out)
}
