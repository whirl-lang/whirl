package parser

import (
	"fmt"

	"github.com/whirl-lang/whirl/pkg/codegen"
	"github.com/whirl-lang/whirl/pkg/lexer"
)

func ParseBody(tokens *lexer.TokenIterator) ([]codegen.Instruction, error) {
	var instructions []codegen.Instruction
	next, err := tokens.Peek()

	if err != nil {
		return nil, err
	}

	// parse curly open
	_, err = ExpectToken(tokens, lexer.CURLYOPEN)

	if err != nil {
		return nil, err
	}

	for next.Kind != lexer.CURLYCLOSE {
		instruction, err := ParseInstruction(tokens)

		if err != nil {
			return nil, err
		}

		instructions = append(instructions, instruction)

		next, err = tokens.Peek()

		if err != nil {
			return nil, err
		}
	}

	// parse curly close
	_, err = ExpectToken(tokens, lexer.CURLYCLOSE)

	if err != nil {
		return nil, err
	}

	return instructions, nil
}

func ParseInstruction(tokens *lexer.TokenIterator) (codegen.Instruction, error) {
	next, err := tokens.Peek()

	if err != nil {
		return nil, err
	}

	switch next.Kind {
	case lexer.LET:
		return ParseAssignment(tokens)
	case lexer.STRUCT:
		return ParseStruct(tokens)
	case lexer.IF:
		return ParseIf(tokens)
	case lexer.ESCAPE:
		return ParseEscape(tokens)
	case lexer.PROC:
		return ParseProcedure(tokens)
	}

	_, err = ExpectToken(tokens, lexer.EOF)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func ParseType(tokens *lexer.TokenIterator) (codegen.Type, error) {
	tok, err := tokens.Next()

	if err != nil {
		return nil, err
	}

	switch tok.Kind {
	case lexer.INT:
		return codegen.Int{}, nil
	case lexer.STRING:
		return codegen.String{}, nil
	case lexer.BOOLEAN:
		return codegen.Bool{}, nil
	case lexer.CHAR:
		return codegen.Char{}, nil
	case lexer.VOID:
		return codegen.Void{}, nil
	case lexer.IDENT:
		return codegen.Ident{Name: tok.Value}, nil
	}

	return nil, fmt.Errorf("unexpected token %s", lexer.TokensPretty[tok.Kind])
}

func ParseIdent(tokens *lexer.TokenIterator) (codegen.Ident, error) {
	token, err := ExpectToken(tokens, lexer.IDENT)

	if err != nil {
		return codegen.Ident{}, err
	}

	return codegen.Ident{Name: token.Value}, nil
}

func ExpectToken(tokens *lexer.TokenIterator, token int) (lexer.Token, error) {
	tok, err := tokens.Peek()

	if err != nil {
		return tok, err
	}

	if tok.Kind != token {
		return tok, fmt.Errorf("expected %s, got %s", lexer.TokensPretty[token], lexer.TokensPretty[tok.Kind])
	}

	return tokens.Next()
}

type InstructionIterator struct {
	Tokens lexer.TokenIterator
}

func (iter *InstructionIterator) Next() (codegen.Instruction, error) {
	return ParseInstruction(&iter.Tokens)
}

func Iterator(tokens lexer.TokenIterator) InstructionIterator {
	return InstructionIterator{tokens}
}
