package parser

import (
	"fmt"
	"strconv"

	"github.com/whirl-lang/whirl/pkg/codegen"
	"github.com/whirl-lang/whirl/pkg/lexer"
)

func ParseStructInit(tokens *lexer.TokenIterator) (codegen.StructInit, error) {
	// get "struct"
	_, err := ExpectToken(tokens, lexer.STRUCT)

	if err != nil {
		return codegen.StructInit{}, err
	}

	// get ident
	ident, err := ParseIdent(tokens)

	if err != nil {
		return codegen.StructInit{}, err
	}

	// get open brace
	_, err = ExpectToken(tokens, lexer.CURLYOPEN)

	if err != nil {
		return codegen.StructInit{}, err
	}

	structure := codegen.StructInit{
		Ident: ident.Name,
	}
	next, err := tokens.Peek()

	if err != nil {
		return codegen.StructInit{}, err
	}

	// fields...
	for next.Kind != lexer.CURLYCLOSE {
		field, err := ParseInitField(tokens)

		if err != nil {
			return codegen.StructInit{}, err
		}

		// get comma
		_, err = ExpectToken(tokens, lexer.COMMA)

		if err != nil {
			return codegen.StructInit{}, err
		}

		// add field to struct
		structure.Fields = append(structure.Fields, field)

		next, err = tokens.Peek()

		if err != nil {
			return codegen.StructInit{}, err
		}
	}

	// get close brace
	_, err = ExpectToken(tokens, lexer.CURLYCLOSE)

	if err != nil {
		return codegen.StructInit{}, err
	}

	return structure, nil
}

func ParseInitField(tokens *lexer.TokenIterator) (codegen.FieldInit, error) {
	// get ident
	ident, err := ParseIdent(tokens)

	if err != nil {
		return codegen.FieldInit{}, err
	}

	// get colon
	_, err = ExpectToken(tokens, lexer.COLON)

	if err != nil {
		return codegen.FieldInit{}, err
	}

	// get expression
	expr, err := ParseExpr(tokens)

	if err != nil {
		return codegen.FieldInit{}, err
	}

	return codegen.FieldInit{Ident: ident.Name, Expr: expr}, nil
}

func ParseExpr(tokens *lexer.TokenIterator) (codegen.Expr, error) {
	next, err := tokens.Peek()

	if err != nil {
		return codegen.ExprMath{}, err
	}

	counter := 0
	expr := codegen.ExprMath{}

	for {
		if next.Kind == lexer.PARENOPEN {
			counter++
		}

		if next.Kind == lexer.PARENCLOSE {
			counter--
		}

		if next.IsSeparator() && counter == 0 && next.Kind != lexer.PARENCLOSE {
			break
		}

		if counter < 0 {
			return expr, fmt.Errorf("unexpected token %s", lexer.TokensPretty[next.Kind])
		}

		token, err := tokens.Next()

		if err != nil {
			return expr, err
		}

		expr.Tokens = append(expr.Tokens, token)

		next, err = tokens.Peek()

		if err != nil {
			return expr, err
		}

		if next.Kind == lexer.PARENCLOSE && counter == 0 {
			break
		}
	}

	return expr, nil
}

func ParseInt(tokens *lexer.TokenIterator) (codegen.Int, error) {
	token, err := ExpectToken(tokens, lexer.INT)

	if err != nil {
		return codegen.Int{}, err
	}

	value, err := strconv.ParseInt(token.Value, 10, 64)

	if err != nil {
		return codegen.Int{}, err
	}

	return codegen.Int{Value: value}, nil
}

func ParseString(tokens *lexer.TokenIterator) (codegen.String, error) {
	token, err := ExpectToken(tokens, lexer.STRING)

	if err != nil {
		return codegen.String{}, err
	}

	return codegen.String{Value: token.Value}, nil
}

func ParseBool(tokens *lexer.TokenIterator) (codegen.Bool, error) {
	token, err := ExpectToken(tokens, lexer.BOOLEAN)

	if err != nil {
		return codegen.Bool{}, err
	}

	return codegen.Bool{Value: token.Value == "true"}, nil
}

func ParseVoid(tokens *lexer.TokenIterator) (codegen.Void, error) {
	_, err := ExpectToken(tokens, lexer.VOID)

	if err != nil {
		return codegen.Void{}, err
	}

	return codegen.Void{}, nil
}

func ParseChar(tokens *lexer.TokenIterator) (codegen.Char, error) {
	token, err := ExpectToken(tokens, lexer.CHAR)

	if err != nil {
		return codegen.Char{}, err
	}

	return codegen.Char{Value: rune(token.Value[0])}, nil
}
