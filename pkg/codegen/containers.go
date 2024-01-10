package codegen

import "github.com/whirl-lang/whirl/pkg/lexer"

type Instruction interface {
	CInstruction() string
}

type Type interface {
	CType() string
}

type Value interface {
	CValue() string
}

type Procedure struct {
	Ident        string
	Args         []Argument
	Instructions []Instruction
	ReturnType   Type
}

type Argument struct {
	Ident string
	Type  Type
}

type ProcedureCall struct {
	Ident string
	Args  []Expr
}

type Struct struct {
	Ident  string
	Fields []Field
}

type StructInit struct {
	Ident  string
	Fields []FieldInit
}

type Field struct {
	Ident string
	Type  Type
}

type FieldInit struct {
	Ident string
	Expr  Expr
}

type If struct {
	Condition Expr
	Body      []Instruction
	Else      []Instruction
}

type Expr interface {
	CValue() string
}

type ExprMath struct {
	Tokens []Expr
}

type ExprToken struct {
	Token lexer.Token
}

type Assignment struct {
	Ident string
	Expr  Expr
	Type  Type
}

type FunctionCall struct {
	Ident string
	Args  []Expr
}

type Escape struct {
	Expr Expr
}

type Array struct {
	Type  Type
	Value []Expr
}

type Int struct {
	Value int64
}

type String struct {
	Value string
}

type Bool struct {
	Value bool
}

type Void struct{}

type Char struct {
	Value rune
}

type Ident struct {
	Name string
}

type Iter struct {
	Ident string
	Lower Expr
	Upper Expr
	Body  []Instruction
}

type Until struct {
	Condition Expr
	Body      []Instruction
}

type Reassign struct {
	Ident string
	Expr  Expr
}

type Break struct{}

type Continue struct{}
