package parser

import (
	"github.com/whirl-lang/whirl/pkg/codegen"
	"github.com/whirl-lang/whirl/pkg/lexer"
)

func ParseAssignment(tokens *lexer.TokenIterator) (codegen.Assignment, error) {
	// get "let"
	_, err := ExpectToken(tokens, lexer.LET)

	if err != nil {
		return codegen.Assignment{}, err
	}

	// get ident
	ident, err := ParseIdent(tokens)

	if err != nil {
		return codegen.Assignment{}, err
	}

	// parse colon
	_, err = ExpectToken(tokens, lexer.COLON)

	if err != nil {
		return codegen.Assignment{}, err
	}

	// parse type
	typ, err := ParseType(tokens)

	if err != nil {
		return codegen.Assignment{}, err
	}

	// get equals
	_, err = ExpectToken(tokens, lexer.ASSIGN)

	if err != nil {
		return codegen.Assignment{}, err
	}

	switch typ.(type) {
	case codegen.Ident:
		structure, err := ParseStructInit(tokens)

		if err != nil {
			return codegen.Assignment{}, err
		}

		// get semi
		_, err = ExpectToken(tokens, lexer.SEMICOLON)

		if err != nil {
			return codegen.Assignment{}, err
		}

		return codegen.Assignment{Ident: ident, Expr: structure, Type: typ}, nil
	}

	// get expression
	expr, err := ParseExpr(tokens)

	if err != nil {
		return codegen.Assignment{}, err
	}

	// get semi
	_, err = ExpectToken(tokens, lexer.SEMICOLON)

	if err != nil {
		return codegen.Assignment{}, err
	}

	return codegen.Assignment{Ident: ident, Expr: expr, Type: typ}, nil
}

func ParseIf(tokens *lexer.TokenIterator) (codegen.If, error) {
	// get "if"
	_, err := ExpectToken(tokens, lexer.IF)

	if err != nil {
		return codegen.If{}, err
	}

	// get condition
	condition, err := ParseExpr(tokens)

	if err != nil {
		return codegen.If{}, err
	}

	// get body
	body, err := ParseBody(tokens)

	if err != nil {
		return codegen.If{}, err
	}

	// get "else"
	_, err = ExpectToken(tokens, lexer.ELSE)

	if err != nil {
		return codegen.If{
			Condition: condition,
			Body:      body,
			Else:      nil,
		}, nil
	}

	// get else body
	elseBody, err := ParseBody(tokens)

	if err != nil {
		return codegen.If{}, err
	}

	return codegen.If{Condition: condition, Body: body, Else: elseBody}, nil
}

func ParseEscape(tokens *lexer.TokenIterator) (codegen.Escape, error) {
	// get "escape"
	_, err := ExpectToken(tokens, lexer.ESCAPE)

	if err != nil {
		return codegen.Escape{}, err
	}

	// get expression
	expr, err := ParseExpr(tokens)

	if err != nil {
		return codegen.Escape{}, err
	}

	// get semi
	_, err = ExpectToken(tokens, lexer.SEMICOLON)

	if err != nil {
		return codegen.Escape{}, err
	}

	return codegen.Escape{Expr: expr}, nil
}

func ParseProcedure(tokens *lexer.TokenIterator) (codegen.Procedure, error) {
	// get "proc"
	_, err := ExpectToken(tokens, lexer.PROC)

	if err != nil {
		return codegen.Procedure{}, err
	}

	// get ident
	ident, err := ParseIdent(tokens)

	if err != nil {
		return codegen.Procedure{}, err
	}

	// get open parens
	_, err = ExpectToken(tokens, lexer.PARENOPEN)

	if err != nil {
		return codegen.Procedure{}, err
	}

	// get args
	args, err := ParseArgs(tokens)

	if err != nil {
		return codegen.Procedure{}, err
	}

	// get close parens
	_, err = ExpectToken(tokens, lexer.PARENCLOSE)

	if err != nil {
		return codegen.Procedure{}, err
	}

	// two colons
	_, err = ExpectToken(tokens, lexer.COLON)

	if err != nil {
		return codegen.Procedure{}, err
	}

	_, err = ExpectToken(tokens, lexer.COLON)

	if err != nil {
		return codegen.Procedure{}, err
	}

	// get return type
	returnType, err := ParseType(tokens)

	if err != nil {

		return codegen.Procedure{}, err
	}

	// get body
	body, err := ParseBody(tokens)

	if err != nil {
		return codegen.Procedure{}, err
	}

	return codegen.Procedure{
		Ident:        ident,
		Args:         args,
		Instructions: body,
		ReturnType:   returnType,
	}, nil
}

func ParseUntil(tokens *lexer.TokenIterator) (codegen.Until, error) {
	// get "until"
	_, err := ExpectToken(tokens, lexer.UNTIL)

	if err != nil {
		return codegen.Until{}, err
	}

	// get condition
	condition, err := ParseExpr(tokens)

	if err != nil {
		return codegen.Until{}, err
	}

	// get body
	body, err := ParseBody(tokens)

	if err != nil {
		return codegen.Until{}, err
	}

	return codegen.Until{Condition: condition, Body: body}, nil
}

func ParseArg(tokens *lexer.TokenIterator) (codegen.Argument, error) {
	// get ident
	ident, err := ParseIdent(tokens)

	if err != nil {
		return codegen.Argument{}, err
	}

	// get colon
	_, err = ExpectToken(tokens, lexer.COLON)

	if err != nil {
		return codegen.Argument{}, err
	}

	// get type
	typ, err := ParseType(tokens)

	if err != nil {
		return codegen.Argument{}, err
	}

	return codegen.Argument{Ident: ident, Type: typ}, nil
}

func ParseReassign(tokens *lexer.TokenIterator, path codegen.Path) (codegen.Reassign, error) {
	// get equals
	_, err := ExpectToken(tokens, lexer.ASSIGN)

	if err != nil {
		return codegen.Reassign{}, err
	}

	// get expression
	expr, err := ParseExpr(tokens)

	if err != nil {
		return codegen.Reassign{}, err
	}

	return codegen.Reassign{Ident: path, Expr: expr}, nil
}

func ParseIter(tokens *lexer.TokenIterator) (codegen.Iter, error) {
	// get "iter"
	_, err := ExpectToken(tokens, lexer.ITER)

	if err != nil {
		return codegen.Iter{}, err
	}

	// get ident
	ident, err := ParseIdent(tokens)

	if err != nil {
		return codegen.Iter{}, err
	}

	_, err = ExpectToken(tokens, lexer.IN)

	if err != nil {
		return codegen.Iter{}, err
	}

	// get lower
	lower, err := ParseExpr(tokens)

	if err != nil {
		return codegen.Iter{}, err
	}

	_, err = ExpectToken(tokens, lexer.COLON)

	if err != nil {
		return codegen.Iter{}, err
	}

	// get upper
	upper, err := ParseExpr(tokens)

	if err != nil {
		return codegen.Iter{}, err
	}

	// get body
	body, err := ParseBody(tokens)

	if err != nil {
		return codegen.Iter{}, err
	}

	return codegen.Iter{Ident: ident, Lower: lower, Upper: upper, Body: body}, nil
}

func ParseBreak(tokens *lexer.TokenIterator) (codegen.Break, error) {

	// get "break"
	_, err := ExpectToken(tokens, lexer.BREAK)

	if err != nil {
		return codegen.Break{}, err
	}

	// get semi
	_, err = ExpectToken(tokens, lexer.SEMICOLON)

	if err != nil {
		return codegen.Break{}, err
	}

	return codegen.Break{}, nil

}

func ParseContinue(tokens *lexer.TokenIterator) (codegen.Continue, error) {

	// get "continue"
	_, err := ExpectToken(tokens, lexer.CONTINUE)

	if err != nil {
		return codegen.Continue{}, err
	}

	// get semi
	_, err = ExpectToken(tokens, lexer.SEMICOLON)

	if err != nil {
		return codegen.Continue{}, err
	}

	return codegen.Continue{}, nil

}

func ParseArgs(tokens *lexer.TokenIterator) ([]codegen.Argument, error) {
	var args []codegen.Argument
	next, err := tokens.Peek()

	if err != nil {
		return nil, err
	}

	for next.Kind != lexer.PARENCLOSE {
		arg, err := ParseArg(tokens)

		if err != nil {
			return nil, err
		}

		args = append(args, arg)

		next, err = tokens.Peek()

		if err != nil {
			return nil, err
		}

		_, err = ExpectToken(tokens, lexer.COMMA)

		if err != nil {
			break
		}
	}

	return args, nil
}

func ParseProcedureCall(tokens *lexer.TokenIterator, path codegen.Path) (codegen.ProcedureCall, error) {

	// get open parens
	_, err := ExpectToken(tokens, lexer.PARENOPEN)

	if err != nil {
		return codegen.ProcedureCall{}, err
	}

	// get args
	next, err := tokens.Peek()

	if err != nil {
		return codegen.ProcedureCall{}, err
	}

	var args []codegen.Expr

	for next.Kind != lexer.PARENCLOSE {
		arg, err := ParseExpr(tokens)

		if err != nil {
			return codegen.ProcedureCall{}, err
		}

		args = append(args, arg)

		next, err = tokens.Peek()

		if err != nil {
			return codegen.ProcedureCall{}, err
		}

		_, err = ExpectToken(tokens, lexer.COMMA)

		if err != nil {
			break
		}
	}

	// get close parens
	_, err = ExpectToken(tokens, lexer.PARENCLOSE)

	if err != nil {
		return codegen.ProcedureCall{}, err
	}

	return codegen.ProcedureCall{Ident: path, Args: args}, nil
}
