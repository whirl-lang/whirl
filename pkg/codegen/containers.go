package codegen

import "github.com/whirl-lang/whirl/pkg/lexer"

type Instruction interface {
	CInstruction(ctx Context) string
}

type Type interface {
	CType(ctx Context) string
}

type Value interface {
	CValue(ctx Context) string
}

type Procedure struct {
	Ident        Ident
	Args         []Argument
	Instructions []Instruction
	ReturnType   Type
}

type Argument struct {
	Ident Ident
	Type  Type
}

type ProcedureCall struct {
	Ident Path
	Args  []Expr
}

type Struct struct {
	Ident  Path
	Fields []Field
}

type StructInit struct {
	Ident  Ident
	Fields []FieldInit
}

type Field struct {
	Ident Ident
	Type  Type
}

type FieldInit struct {
	Ident Ident
	Expr  Expr
}

type If struct {
	Condition Expr
	Body      []Instruction
	Else      []Instruction
}

type Expr interface {
	CValue(ctx Context) string
}

type ExprMath struct {
	Tokens []Expr
}

type ExprToken struct {
	Token lexer.Token
}

type Assignment struct {
	Ident Ident
	Expr  Expr
	Type  Type
}

type FunctionCall struct {
	Ident Path
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
	Ident Ident
	Lower Expr
	Upper Expr
	Body  []Instruction
}

type Until struct {
	Condition Expr
	Body      []Instruction
}

type Reassign struct {
	Ident Path
	Expr  Expr
}

type Break struct{}

type Continue struct{}

type Import struct {
	Path string
}

type Path struct {
	Tokens []Ident
}
