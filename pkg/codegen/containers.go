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

type Expr interface{}

type ExprMath struct {
	Tokens []lexer.Token
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
