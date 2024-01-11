package codegen

import (
	"bytes"
	"fmt"
	"regexp"
)

var reserved = map[string]bool{
	"main":   true,
	"printf": true,
}

func TransformIdent(ctx Context, ident string) string {
	if len(ctx.Namespace) == 0 || reserved[ident] {
		return ident
	}

	return fmt.Sprintf("__whirl_%s_%s", PathToNamespace(ctx.Path), ident)
}

func TransformPath(ctx Context, path Path) string {
	var buffer bytes.Buffer

	// if we're in the main namespace and there's just one token, use it
	if len(ctx.Namespace) == 0 && len(path.Tokens) == 1 {
		return path.Tokens[0].Name
	}

	if len(path.Tokens) == 1 && reserved[path.Tokens[0].Name] {
		return path.Tokens[0].Name
	}

	buffer.WriteString("__whirl_")
	buffer.WriteString(PathToNamespace(ctx.Path))
	buffer.WriteString("_")

	for i, token := range path.Tokens {
		buffer.WriteString(token.Name)

		if i != len(path.Tokens)-1 {
			buffer.WriteString("_")
		}

		// on second-last, append whirl_
		if i == len(path.Tokens)-2 {
			buffer.WriteString("whirl_")
		}
	}

	return buffer.String()
}

func PathToNamespace(path string) string {
	return regexp.
		MustCompile("[^a-zA-Z0-9]").
		ReplaceAllString(path, "_")
}
