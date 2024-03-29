package parser

import (
	"fmt"

	"github.com/whirl-lang/whirl/pkg/codegen"
	"github.com/whirl-lang/whirl/pkg/lexer"
)

func ParseImport(tokens *lexer.TokenIterator) (codegen.Import, error) {
	// get "import"
	_, err := ExpectToken(tokens, lexer.IMPORT)

	if err != nil {
		return codegen.Import{}, err
	}

	// get path
	path, err := ParseString(tokens)

	if err != nil {
		return codegen.Import{}, err
	}

	// get semi
	_, err = ExpectToken(tokens, lexer.SEMICOLON)

	if err != nil {
		return codegen.Import{}, err
	}

	return codegen.Import{Path: path.Value[1 : len(path.Value)-1]}, nil
}

func ParseBody(tokens *lexer.TokenIterator) ([]codegen.Instruction, error) {
	// parse curly open
	_, err := ExpectToken(tokens, lexer.CURLYOPEN)

	if err != nil {
		// parse one instruction
		instruction, err := ParseInstruction(tokens)

		if err != nil {
			return nil, err
		}

		return []codegen.Instruction{instruction}, nil
	}

	var instructions []codegen.Instruction

	next, err := tokens.Peek()

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
	case lexer.BREAK:
		return ParseBreak(tokens)
	case lexer.CONTINUE:
		return ParseContinue(tokens)
	case lexer.STRUCT:
		return ParseStruct(tokens)
	case lexer.IF:
		return ParseIf(tokens)
	case lexer.ESCAPE:
		return ParseEscape(tokens)
	case lexer.PROC:
		return ParseProcedure(tokens)
	case lexer.UNTIL:
		return ParseUntil(tokens)
	case lexer.ITER:
		return ParseIter(tokens)
	case lexer.IMPORT:
		return ParseImport(tokens)
	case lexer.IDENT:
		path, err := ParsePath(tokens)

		if err != nil {
			return codegen.ProcedureCall{}, err
		}

		var instruction codegen.Instruction
		instruction, err = ParseProcedureCall(tokens, path)

		if err != nil {
			instruction, err = ParseReassign(tokens, path)

			if err != nil {
				return nil, err
			}
		}

		_, err = ExpectToken(tokens, lexer.SEMICOLON)

		if err != nil {
			return nil, err
		}

		return instruction, nil
	}

	_, err = ExpectToken(tokens, lexer.EOF)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func ParseType(tokens *lexer.TokenIterator) (codegen.Type, error) {
	tok, err := tokens.Next()
	//fmt.Println(tok)
	if err != nil {
		return nil, err
	}

	var typ codegen.Type

	switch tok.Kind {
	case lexer.INT:
		typ = codegen.Int{}
	case lexer.STRING:
		typ = codegen.String{}
	case lexer.BOOLEAN:
		typ = codegen.Bool{}
	case lexer.CHAR:
		typ = codegen.Char{}
	case lexer.VOID:
		typ = codegen.Void{}
	case lexer.IDENT:
		typ = codegen.Ident{Name: tok.Value}
	default:
		return nil, fmt.Errorf("unexpected token %s", lexer.TokensPretty[tok.Kind])
	}

	tok, err = tokens.Peek()

	if err != nil {
		return nil, err
	}

	if tok.Kind == lexer.BRACKETOPEN {
		_, err = ExpectToken(tokens, lexer.BRACKETOPEN)

		if err != nil {

			return nil, err
		}

		_, err = ExpectToken(tokens, lexer.BRACKETCLOSE)

		if err != nil {

			return nil, err
		}

		typ = codegen.Array{Type: typ}
	}

	return typ, nil
}

func ParseIdent(tokens *lexer.TokenIterator) (codegen.Ident, error) {
	token, err := ExpectToken(tokens, lexer.IDENT)

	if err != nil {
		return codegen.Ident{}, err
	}

	return codegen.Ident{Name: token.Value}, nil
}

func ParsePath(tokens *lexer.TokenIterator) (codegen.Path, error) {
	token, err := ExpectToken(tokens, lexer.IDENT)

	if err != nil {
		return codegen.Path{}, err
	}

	next, err := tokens.Peek()

	if err != nil {
		return codegen.Path{}, err
	}

	path := codegen.Path{
		Tokens: []codegen.Ident{{Name: token.Value}},
	}

	for next.Kind == lexer.COLONCOLON {
		_, err = ExpectToken(tokens, lexer.COLONCOLON)

		if err != nil {
			return codegen.Path{}, err
		}

		token, err = ExpectToken(tokens, lexer.IDENT)

		if err != nil {
			return codegen.Path{}, err
		}

		path.Tokens = append(path.Tokens, codegen.Ident{Name: token.Value})

		next, err = tokens.Peek()

		if err != nil {
			return codegen.Path{}, err
		}
	}

	return path, nil
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
