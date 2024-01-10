package parser

import (
	"fmt"
	"strconv"

	lexer "github.com/whirl-lang/whirl/pkg/lexer"
)

type Instruction interface{}

type Procedure struct {
	Indentifier  string
	Arguments    []Argument
	Instructions []Instruction
	ReturnType   Type
}

type Argument struct {
	Identifier string
	Type       Type
}

type ProcedureCall struct {
	Identifier string
	Arguments  []Expression
}

func ParseProcedureCall(tokens *lexer.TokenIterator) (ProcedureCall, error) {
	// get ident
	ident, err := ParseIdent(tokens)

	if err != nil {
		return ProcedureCall{}, err
	}

	// get open parens
	_, err = ExpectToken(tokens, lexer.PARENOPEN)

	if err != nil {
		return ProcedureCall{}, err
	}

	// get args
	next, err := tokens.Peek()

	if err != nil {
		return ProcedureCall{}, err
	}

	var args []Expression

	// FIXME: make this work
	for next.Kind != lexer.PARENCLOSE {
		arg, err := ParseExpression(tokens)

		if err != nil {
			return ProcedureCall{}, err
		}

		args = append(args, arg)

		next, err = tokens.Peek()

		if err != nil {
			return ProcedureCall{}, err
		}

		if next.Kind == lexer.COMMA {
		}
	}

	// get close parens
	_, err = ExpectToken(tokens, lexer.PARENCLOSE)

	if err != nil {
		return ProcedureCall{}, err
	}

	return ProcedureCall{ident.Name, args}, nil
}

func ParseProcedure(tokens *lexer.TokenIterator) (Procedure, error) {
	// get "proc"
	_, err := ExpectToken(tokens, lexer.PROC)

	if err != nil {
		return Procedure{}, err
	}

	// get ident
	ident, err := ParseIdent(tokens)

	if err != nil {
		return Procedure{}, err
	}

	// get open parens
	_, err = ExpectToken(tokens, lexer.PARENOPEN)

	if err != nil {
		return Procedure{}, err
	}

	// get args
	args, err := ParseArgs(tokens)

	if err != nil {
		return Procedure{}, err
	}

	// get close parens
	_, err = ExpectToken(tokens, lexer.PARENCLOSE)

	if err != nil {
		return Procedure{}, err
	}

	// two colons
	_, err = ExpectToken(tokens, lexer.COLON)

	if err != nil {
		return Procedure{}, err
	}

	_, err = ExpectToken(tokens, lexer.COLON)

	if err != nil {
		return Procedure{}, err
	}

	// get return type
	returnType, err := ParseType(tokens)

	if err != nil {
		return Procedure{}, err
	}

	// get body
	body, err := ParseBody(tokens)

	if err != nil {
		return Procedure{}, err
	}

	return Procedure{
		Indentifier:  ident.Name,
		Arguments:    args,
		Instructions: body,
		ReturnType:   returnType,
	}, nil
}

func ParseArgs(tokens *lexer.TokenIterator) ([]Argument, error) {
	var args []Argument
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
	}

	return args, nil
}

func ParseArg(tokens *lexer.TokenIterator) (Argument, error) {
	// get ident
	ident, err := ParseIdent(tokens)

	if err != nil {
		return Argument{}, err
	}

	// get colon
	_, err = ExpectToken(tokens, lexer.COLON)

	if err != nil {
		return Argument{}, err
	}

	// get type
	typ, err := ParseType(tokens)

	if err != nil {
		return Argument{}, err
	}

	return Argument{ident.Name, typ}, nil
}

func ParseBody(tokens *lexer.TokenIterator) ([]Instruction, error) {
	var instructions []Instruction
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

func ParseInstruction(tokens *lexer.TokenIterator) (Instruction, error) {
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

func ParseEscape(tokens *lexer.TokenIterator) (Escape, error) {
	// get "escape"
	_, err := ExpectToken(tokens, lexer.ESCAPE)

	if err != nil {
		return Escape{}, err
	}

	// get expression
	expr, err := ParseExpression(tokens)

	if err != nil {
		return Escape{}, err
	}

	// get semi
	_, err = ExpectToken(tokens, lexer.SEMICOLON)

	if err != nil {
		return Escape{}, err
	}

	return Escape{expr}, nil
}

func ParseIf(tokens *lexer.TokenIterator) (If, error) {
	// get "if"
	_, err := ExpectToken(tokens, lexer.IF)

	if err != nil {
		return If{}, err
	}

	// get condition
	condition, err := ParseExpression(tokens)

	if err != nil {
		return If{}, err
	}

	// get body
	body, err := ParseBody(tokens)

	if err != nil {
		return If{}, err
	}

	// get "else"
	_, err = ExpectToken(tokens, lexer.ELSE)

	if err != nil {
		return If{
			condition,
			body,
			nil,
		}, err
	}

	// get else body
	elseBody, err := ParseBody(tokens)

	if err != nil {
		return If{}, err
	}

	return If{condition, body, elseBody}, nil
}

func ParseAssignment(tokens *lexer.TokenIterator) (Assignment, error) {
	// get "let"
	_, err := ExpectToken(tokens, lexer.LET)

	if err != nil {
		return Assignment{}, err
	}

	// get ident
	ident, err := ParseIdent(tokens)

	if err != nil {
		return Assignment{}, err
	}

	// parse colon
	_, err = ExpectToken(tokens, lexer.COLON)

	if err != nil {
		return Assignment{}, err
	}

	// parse type
	typ, err := ParseType(tokens)

	if err != nil {
		return Assignment{}, err
	}

	// get equals
	_, err = ExpectToken(tokens, lexer.ASSIGN)

	if err != nil {
		return Assignment{}, err
	}

	switch typ.(type) {
	case Ident:
		structure, err := ParseStructInit(tokens)

		if err != nil {
			return Assignment{}, err
		}

		// get semi
		_, err = ExpectToken(tokens, lexer.SEMICOLON)

		if err != nil {
			return Assignment{}, err
		}

		if ident.Name != structure.Identifier {
			return Assignment{}, fmt.Errorf("expected identifier %s, got %s", ident.Name, structure.Identifier)
		}

		return Assignment{ident.Name, structure, typ}, nil
	}

	// get expression
	expr, err := ParseExpression(tokens)

	if err != nil {
		return Assignment{}, err
	}

	// get semi
	_, err = ExpectToken(tokens, lexer.SEMICOLON)

	if err != nil {
		return Assignment{}, err
	}

	return Assignment{ident.Name, expr, typ}, nil
}

type Struct struct {
	Identifier string
	Fields     []Field
}

func ParseStruct(tokens *lexer.TokenIterator) (Struct, error) {
	// get "struct"
	_, err := ExpectToken(tokens, lexer.STRUCT)

	if err != nil {
		return Struct{}, err
	}

	// get ident
	ident, err := ParseIdent(tokens)

	if err != nil {
		return Struct{}, err
	}

	// get open brace
	_, err = ExpectToken(tokens, lexer.CURLYOPEN)

	if err != nil {
		return Struct{}, err
	}

	structure := Struct{
		Identifier: ident.Name,
	}
	next, err := tokens.Peek()

	if err != nil {
		return Struct{}, err
	}

	// fields...
	for next.Kind != lexer.CURLYCLOSE {
		field, err := ParseField(tokens)

		if err != nil {
			return Struct{}, err
		}

		// get comma
		_, err = ExpectToken(tokens, lexer.COMMA)

		if err != nil {
			return Struct{}, err
		}

		// add field to struct
		structure.Fields = append(structure.Fields, field)

		next, err = tokens.Peek()

		if err != nil {
			return Struct{}, err
		}
	}

	// get close brace
	_, err = ExpectToken(tokens, lexer.CURLYCLOSE)

	if err != nil {
		return Struct{}, err
	}

	return structure, nil
}

type StructInit struct {
	Identifier string
	Fields     []FieldInit
}

func ParseStructInit(tokens *lexer.TokenIterator) (StructInit, error) {
	// get "struct"
	_, err := ExpectToken(tokens, lexer.STRUCT)

	if err != nil {
		return StructInit{}, err
	}

	// get ident
	ident, err := ParseIdent(tokens)

	if err != nil {
		return StructInit{}, err
	}

	// get open brace
	_, err = ExpectToken(tokens, lexer.CURLYOPEN)

	if err != nil {
		return StructInit{}, err
	}

	structure := StructInit{
		Identifier: ident.Name,
	}
	next, err := tokens.Peek()

	if err != nil {
		return StructInit{}, err
	}

	// fields...
	for next.Kind != lexer.CURLYCLOSE {
		field, err := ParseInitField(tokens)

		if err != nil {
			return StructInit{}, err
		}

		// get comma
		_, err = ExpectToken(tokens, lexer.COMMA)

		if err != nil {
			return StructInit{}, err
		}

		// add field to struct
		structure.Fields = append(structure.Fields, field)

		next, err = tokens.Peek()

		if err != nil {
			return StructInit{}, err
		}
	}

	// get close brace
	_, err = ExpectToken(tokens, lexer.CURLYCLOSE)

	if err != nil {
		return StructInit{}, err
	}

	return structure, nil
}

type Field struct {
	Identifier string
	Type       Type
}

func ParseField(tokens *lexer.TokenIterator) (Field, error) {
	// get ident
	ident, err := ParseIdent(tokens)

	if err != nil {
		return Field{}, err
	}

	// get colon
	_, err = ExpectToken(tokens, lexer.COLON)

	if err != nil {
		return Field{}, err
	}

	// get type
	typ, err := ParseType(tokens)

	if err != nil {
		return Field{}, err
	}

	return Field{ident.Name, typ}, nil
}

type FieldInit struct {
	Identifier string
	Expression Expression
}

func ParseInitField(tokens *lexer.TokenIterator) (FieldInit, error) {
	// get ident
	ident, err := ParseIdent(tokens)

	if err != nil {
		return FieldInit{}, err
	}

	// get colon
	_, err = ExpectToken(tokens, lexer.COLON)

	if err != nil {
		return FieldInit{}, err
	}

	// get expression
	expr, err := ParseExpression(tokens)

	if err != nil {
		return FieldInit{}, err
	}

	return FieldInit{ident.Name, expr}, nil
}

type If struct {
	Condition Expression
	Body      []Instruction
	Else      []Instruction
}

type Expression interface{}

type ExpressionMath struct {
	Tokens []lexer.Token
}

func ParseExpression(tokens *lexer.TokenIterator) (Expression, error) {
	next, err := tokens.Peek()

	if err != nil {
		return ExpressionMath{}, err
	}

	counter := 0
	expr := ExpressionMath{}

	// capture until open parens equal closed parens
	// and the next token isn't a closeable one
	for !next.IsSeparator() {
		if next.Kind == lexer.CURLYOPEN {
			counter++
		}

		if next.Kind == lexer.CURLYCLOSE {
			counter--
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
	}

	return expr, nil
}

type Assignment struct {
	Identifier string
	Expression Expression
	Type       Type
}

type FunctionCall struct {
	Identifier string
	Arguments  []Expression
}

type Escape struct {
	Expression Expression
}

type Type interface {
	CType() string
}

type Value interface {
	CValue() string
}

func ParseType(tokens *lexer.TokenIterator) (Type, error) {
	tok, err := tokens.Next()

	if err != nil {
		return nil, err
	}

	switch tok.Kind {
	case lexer.INT:
		return Int{}, nil
	case lexer.STRING:
		return String{}, nil
	case lexer.BOOLEAN:
		return Bool{}, nil
	case lexer.CHAR:
		return Char{}, nil
	case lexer.VOID:
		return Void{}, nil
	case lexer.IDENT:
		return Ident{Name: tok.Value}, nil
	}

	return nil, fmt.Errorf("unexpected token %s", lexer.TokensPretty[tok.Kind])
}

type Int struct {
	Value int64
}

func (i Int) CType() string {
	return "int64_t"
}

func (i Int) CValue() string {
	return strconv.FormatInt(i.Value, 10)
}

func ParseInt(tokens *lexer.TokenIterator) (Int, error) {
	token, err := ExpectToken(tokens, lexer.INT)

	if err != nil {
		return Int{}, err
	}

	value, err := strconv.ParseInt(token.Value, 10, 64)

	if err != nil {
		return Int{}, err
	}

	return Int{value}, nil
}

type String struct {
	Value string
}

func (s String) CType() string {
	return "char*"
}

func (s String) CValue() string {
	return "\"" + s.Value + "\""
}

func ParseString(tokens *lexer.TokenIterator) (String, error) {
	token, err := ExpectToken(tokens, lexer.STRING)

	if err != nil {
		return String{}, err
	}

	return String{token.Value}, nil
}

type Bool struct {
	Value bool
}

func (b Bool) CType() string {
	return "bool"
}

func (b Bool) CValue() string {
	if b.Value {
		return "true"
	}

	return "false"
}

func ParseBool(tokens *lexer.TokenIterator) (Bool, error) {
	token, err := ExpectToken(tokens, lexer.BOOLEAN)

	if err != nil {
		return Bool{}, err
	}

	return Bool{token.Value == "true"}, nil
}

type Void struct{}

func (v Void) CType() string {
	return "void"
}

func ParseVoid(tokens *lexer.TokenIterator) (Void, error) {
	_, err := ExpectToken(tokens, lexer.VOID)

	if err != nil {
		return Void{}, err
	}

	return Void{}, nil
}

type Char struct {
	Value rune
}

func (c Char) CType() string {
	return "char"
}

func (c Char) CValue() string {
	return "'" + string(c.Value) + "'"
}

func ParseChar(tokens *lexer.TokenIterator) (Char, error) {
	token, err := ExpectToken(tokens, lexer.CHAR)

	if err != nil {
		return Char{}, err
	}

	return Char{rune(token.Value[0])}, nil
}

type Ident struct {
	Name string
}

func (i Ident) CType() string {
	return i.Name
}

func ParseIdent(tokens *lexer.TokenIterator) (Ident, error) {
	token, err := ExpectToken(tokens, lexer.IDENT)

	if err != nil {
		return Ident{}, err
	}

	return Ident{Name: token.Value}, nil
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

func (iter *InstructionIterator) Next() (Instruction, error) {
	return ParseInstruction(&iter.Tokens)
}

func Iterator(tokens lexer.TokenIterator) InstructionIterator {
	return InstructionIterator{tokens}
}
