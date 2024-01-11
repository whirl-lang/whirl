package parser

import (
	"github.com/whirl-lang/whirl/pkg/codegen"
	"github.com/whirl-lang/whirl/pkg/lexer"
)

func ParseStruct(tokens *lexer.TokenIterator) (codegen.Struct, error) {
	// get "struct"
	_, err := ExpectToken(tokens, lexer.STRUCT)

	if err != nil {
		return codegen.Struct{}, err
	}

	// get ident
	path, err := ParsePath(tokens)

	if err != nil {
		return codegen.Struct{}, err
	}

	// get open brace
	_, err = ExpectToken(tokens, lexer.CURLYOPEN)

	if err != nil {
		return codegen.Struct{}, err
	}

	structure := codegen.Struct{
		Ident: path,
	}

	next, err := tokens.Peek()

	if err != nil {
		return codegen.Struct{}, err
	}

	// fields...
	for next.Kind != lexer.CURLYCLOSE {
		field, err := ParseField(tokens)

		if err != nil {
			return codegen.Struct{}, err
		}

		// get comma
		_, err = ExpectToken(tokens, lexer.COMMA)

		if err != nil {
			return codegen.Struct{}, err
		}

		// add field to struct
		structure.Fields = append(structure.Fields, field)

		next, err = tokens.Peek()

		if err != nil {
			return codegen.Struct{}, err
		}
	}

	// get close brace
	_, err = ExpectToken(tokens, lexer.CURLYCLOSE)

	if err != nil {
		return codegen.Struct{}, err
	}

	return structure, nil
}

func ParseField(tokens *lexer.TokenIterator) (codegen.Field, error) {
	// get ident
	ident, err := ParseIdent(tokens)

	if err != nil {
		return codegen.Field{}, err
	}

	// get colon
	_, err = ExpectToken(tokens, lexer.COLON)

	if err != nil {
		return codegen.Field{}, err
	}

	// get type
	typ, err := ParseType(tokens)

	if err != nil {
		return codegen.Field{}, err
	}

	return codegen.Field{Ident: ident, Type: typ}, nil
}
