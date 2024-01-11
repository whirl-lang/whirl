package codegen

import (
	"fmt"
	"regexp"
)

func TransformIdent(ctx Context, ident string) string {
	if len(ctx.Namespace) == 0 {
		return ident
	}

	return fmt.Sprintf("__whirl_%s_%s", ctx.Namespace, ident)
}

func PathToNamespace(path string) string {
	return regexp.
		MustCompile("[^a-zA-Z0-9]").
		ReplaceAllString(path, "_")
}
