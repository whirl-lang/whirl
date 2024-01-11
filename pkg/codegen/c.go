package codegen

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"strconv"
)

type CType interface {
	CType(ctx Context) string
}

type CValue interface {
	CValue(ctx Context) string
}

type CInstruction interface {
	CInstruction(ctx Context) string
}

func (i Int) CType(ctx Context) string {
	return "int"
}

func (s String) CType(ctx Context) string {
	return "char*"
}

func (s String) CValue(ctx Context) string {
	return "\"" + s.Value + "\""
}

func (i Int) CValue(ctx Context) string {
	return strconv.FormatInt(i.Value, 10)
}

func (b Bool) CType(ctx Context) string {
	return "bool"
}

func (b Bool) CValue(ctx Context) string {
	if b.Value {
		return "true"
	}

	return "false"
}

func (c Char) CType(ctx Context) string {
	return "char"
}

func (c Char) CValue(ctx Context) string {
	return "'" + string(c.Value) + "'"
}

func (v Void) CType(ctx Context) string {
	return "void"
}

func (i Ident) CType(ctx Context) string {
	return TransformIdent(ctx, i.Name)
}

func (p Path) CType(ctx Context) string {
	return p.CValue(ctx)
}

func (p Path) CValue(ctx Context) string {
	var buffer bytes.Buffer

	buffer.WriteString("__whirl_")

	for i, token := range p.Tokens {
		buffer.WriteString(token.CType(ctx))

		if i != len(p.Tokens)-1 {
			buffer.WriteString("_")
		}
	}

	return buffer.String()
}

func (a Array) CType(ctx Context) string {
	return a.Type.CType(ctx)
}

func (a Array) Brackets() string {
	switch a.Type.(type) {
	case Array:
		//incase we somehow ever need 2d array (bad we need Vec)
		return "[]" + a.Type.(Array).Brackets()
	}

	return "[]"
}

func (a Array) CValue(ctx Context) string {
	var buffer bytes.Buffer

	buffer.WriteString("{")

	for i, value := range a.Value {
		buffer.WriteString(value.(Value).CValue(ctx))

		if i != len(a.Value)-1 {
			buffer.WriteString(", ")
		}
	}

	buffer.WriteString("}")

	return buffer.String()
}

func (p Procedure) CInstruction(ctx Context) string {
	var buffer bytes.Buffer

	buffer.WriteString(p.ReturnType.CType(ctx))
	buffer.WriteString(" ")
	buffer.WriteString(p.Ident.CType(ctx))
	buffer.WriteString("(")

	for i, arg := range p.Args {
		buffer.WriteString(arg.Type.CType(ctx))
		buffer.WriteString(" ")
		buffer.WriteString(arg.Ident.CType(ctx))

		if i != len(p.Args)-1 {
			buffer.WriteString(", ")
		}
	}

	buffer.WriteString(") { ")

	for _, instruction := range p.Instructions {
		buffer.WriteString(instruction.CInstruction(ctx))
		buffer.WriteString(" ")
	}

	buffer.WriteString("}")

	return buffer.String()
}

func (s Struct) CType(ctx Context) string {
	var buffer bytes.Buffer

	buffer.WriteString("struct ")
	buffer.WriteString(s.Ident.CType(ctx))
	buffer.WriteString(" { ")

	for _, field := range s.Fields {
		buffer.WriteString(field.Type.CType(ctx))
		buffer.WriteString(" ")
		buffer.WriteString(field.Ident.CType(ctx))

		buffer.WriteString("; ")
	}

	buffer.WriteString(" }")

	return buffer.String()
}

func (a Assignment) CInstruction(ctx Context) string {

	switch a.Type.(type) {
	case Ident:
		return fmt.Sprintf("struct %s %s = %s;", a.Type.CType(ctx), a.Ident, a.Expr.(Value).CValue(ctx))
	case Array:
		return fmt.Sprintf("%s %s%s = %s;", a.Type.CType(ctx), a.Ident, a.Type.(Array).Brackets(), a.Expr.(Value).CValue(ctx))
	}

	return fmt.Sprintf("%s %s = %s;", a.Type.CType(ctx), a.Ident, a.Expr.(Value).CValue(ctx))
}

func (s StructInit) CValue(ctx Context) string {
	var buffer bytes.Buffer

	buffer.WriteString(" { ")

	for i, field := range s.Fields {
		buffer.WriteString(".")
		buffer.WriteString(field.Ident.CType(ctx))
		buffer.WriteString(" = ")
		buffer.WriteString(field.Expr.(Value).CValue(ctx))

		if i != len(s.Fields)-1 {
			buffer.WriteString(", ")
		}
	}

	buffer.WriteString(" }")

	return buffer.String()
}

func (s Struct) CInstruction(ctx Context) string {
	return s.CType(ctx) + ";"
}

func (u Until) CInstruction(ctx Context) string {
	var buffer bytes.Buffer

	buffer.WriteString("while (!(")
	buffer.WriteString(u.Condition.(Value).CValue(ctx))
	buffer.WriteString(")) { ")

	for _, instruction := range u.Body {
		buffer.WriteString(instruction.CInstruction(ctx))
		buffer.WriteString(" ")
	}

	buffer.WriteString("}")

	return buffer.String()
}

func (i If) CInstruction(ctx Context) string {
	var buffer bytes.Buffer

	buffer.WriteString("if (")
	buffer.WriteString(i.Condition.(Value).CValue(ctx))
	buffer.WriteString(") { ")

	for _, instruction := range i.Body {
		buffer.WriteString(instruction.CInstruction(ctx))
		buffer.WriteString(" ")
	}

	if len(i.Else) == 0 {
		buffer.WriteString("}")

		return buffer.String()
	}

	buffer.WriteString("} else { ")

	for _, instruction := range i.Else {
		buffer.WriteString(instruction.CInstruction(ctx))
		buffer.WriteString(" ")
	}

	buffer.WriteString("}")

	return buffer.String()
}

func (esc Escape) CInstruction(ctx Context) string {
	return fmt.Sprintf("return %s;", esc.Expr.(Value).CValue(ctx))
}

func (r Reassign) CInstruction(ctx Context) string {
	return fmt.Sprintf("%s = %s;", r.Ident, r.Expr.(Value).CValue(ctx))
}

func (b Break) CInstruction(ctx Context) string {
	return "break;"
}

func (c Continue) CInstruction(ctx Context) string {
	return "continue;"
}

func (i Iter) CInstruction(ctx Context) string {
	var buffer bytes.Buffer

	buffer.WriteString("for (int ")
	buffer.WriteString(i.Ident.CType(ctx))
	buffer.WriteString(" = ")
	buffer.WriteString(i.Lower.CValue(ctx))
	buffer.WriteString("; ")
	buffer.WriteString(i.Ident.CType(ctx))
	buffer.WriteString(" < ")
	buffer.WriteString(i.Upper.CValue(ctx))
	buffer.WriteString("; ")
	buffer.WriteString(i.Ident.CType(ctx))
	buffer.WriteString("++) { ")

	for _, instruction := range i.Body {
		buffer.WriteString(instruction.CInstruction(ctx))
		buffer.WriteString(" ")
	}

	buffer.WriteString("}")

	return buffer.String()
}

func (e ExprMath) CValue(ctx Context) string {
	var buffer bytes.Buffer

	for _, token := range e.Tokens {
		buffer.WriteString(token.CValue(ctx))
	}

	return buffer.String()
}

func (e ExprToken) CValue(ctx Context) string {
	return e.Token.Value
}

func (p ProcedureCall) CValue(ctx Context) string {
	var buffer bytes.Buffer

	buffer.WriteString(p.Ident.CType(ctx))
	buffer.WriteString("(")

	for i, arg := range p.Args {
		buffer.WriteString(arg.(Value).CValue(ctx))

		if i != len(p.Args)-1 {
			buffer.WriteString(", ")
		}
	}

	buffer.WriteString(")")

	return buffer.String()
}

func (p ProcedureCall) CInstruction(ctx Context) string {
	return fmt.Sprintf("%s;", p.CValue(ctx))
}

func (i Import) CInstruction(ctx Context) string {
	c, err := fs.ReadFile(os.DirFS(i.Root), i.Path)

	if err != nil {
		panic(err)
	}

	// FIXME: relative paths pointing to the same file from different
	// locations should not generate different namespaces
	return ctx.Transpile(c, PathToNamespace(i.Path))
}
