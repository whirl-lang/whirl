package codegen

import (
	"bytes"
	"fmt"
	"strconv"
)

type CType interface {
	CType() string
}

type CValue interface {
	CValue() string
}

type CInstruction interface {
	CInstruction() string
}

func (i Int) CType() string {
	return "int"
}

func (s String) CType() string {
	return "char*"
}

func (s String) CValue() string {
	return "\"" + s.Value + "\""
}

func (i Int) CValue() string {
	return strconv.FormatInt(i.Value, 10)
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

func (c Char) CType() string {
	return "char"
}

func (c Char) CValue() string {
	return "'" + string(c.Value) + "'"
}

func (v Void) CType() string {
	return "void"
}

func (i Ident) CType() string {
	return i.Name
}

func (a Array) CType() string {
	return a.Type.CType()
}

func (a Array) Brackets() string {
	switch a.Type.(type) {
	case Array:
		//incase we somehow ever need 2d array (bad we need Vec)
		return "[]" + a.Type.(Array).Brackets()
	}

	return "[]"
}

func (a Array) CValue() string {
	var buffer bytes.Buffer

	buffer.WriteString("{")

	for i, value := range a.Value {
		buffer.WriteString(value.(Value).CValue())

		if i != len(a.Value)-1 {
			buffer.WriteString(", ")
		}
	}

	buffer.WriteString("}")

	return buffer.String()
}

func (p Procedure) CInstruction() string {
	var buffer bytes.Buffer

	buffer.WriteString(p.ReturnType.CType())
	buffer.WriteString(" ")
	buffer.WriteString(p.Ident)
	buffer.WriteString("(")

	for i, arg := range p.Args {
		buffer.WriteString(arg.Type.CType())
		buffer.WriteString(" ")
		buffer.WriteString(arg.Ident)

		if i != len(p.Args)-1 {
			buffer.WriteString(", ")
		}
	}

	buffer.WriteString(") { ")

	for _, instruction := range p.Instructions {
		buffer.WriteString(instruction.CInstruction())
		buffer.WriteString(" ")
	}

	buffer.WriteString("}")

	return buffer.String()
}

func (s Struct) CType() string {
	var buffer bytes.Buffer

	buffer.WriteString("struct ")
	buffer.WriteString(s.Ident)
	buffer.WriteString(" { ")

	for _, field := range s.Fields {
		buffer.WriteString(field.Type.CType())
		buffer.WriteString(" ")
		buffer.WriteString(field.Ident)

		buffer.WriteString("; ")
	}

	buffer.WriteString(" }")

	return buffer.String()
}

func (a Assignment) CInstruction() string {

	switch a.Type.(type) {
	case Ident:
		return fmt.Sprintf("struct %s %s = %s;", a.Type.CType(), a.Ident, a.Expr.(Value).CValue())
	case Array:
		return fmt.Sprintf("%s %s%s = %s;", a.Type.CType(), a.Ident, a.Type.(Array).Brackets(), a.Expr.(Value).CValue())
	}

	return fmt.Sprintf("%s %s = %s;", a.Type.CType(), a.Ident, a.Expr.(Value).CValue())
}

func (s StructInit) CValue() string {
	var buffer bytes.Buffer

	buffer.WriteString(" { ")

	for i, field := range s.Fields {
		buffer.WriteString(".")
		buffer.WriteString(field.Ident)
		buffer.WriteString(" = ")
		buffer.WriteString(field.Expr.(Value).CValue())

		if i != len(s.Fields)-1 {
			buffer.WriteString(", ")
		}
	}

	buffer.WriteString(" }")

	return buffer.String()
}

func (s Struct) CInstruction() string {
	return s.CType() + ";"
}

func (i If) CInstruction() string {
	var buffer bytes.Buffer

	buffer.WriteString("if (")
	buffer.WriteString(i.Condition.(Value).CValue())
	buffer.WriteString(") { ")

	for _, instruction := range i.Body {
		buffer.WriteString(instruction.CInstruction())
		buffer.WriteString(" ")
	}

	if len(i.Else) == 0 {
		buffer.WriteString("}")

		return buffer.String()
	}

	buffer.WriteString("} else { ")

	for _, instruction := range i.Else {
		buffer.WriteString(instruction.CInstruction())
		buffer.WriteString(" ")
	}

	buffer.WriteString("}")

	return buffer.String()
}

func (esc Escape) CInstruction() string {
	return fmt.Sprintf("return %s;", esc.Expr.(Value).CValue())
}

func (e ExprMath) CValue() string {
	var buffer bytes.Buffer

	for _, token := range e.Tokens {
		buffer.WriteString(token.CValue())
	}

	return buffer.String()
}

func (e ExprToken) CValue() string {
	return e.Token.Value
}

func (p ProcedureCall) CValue() string {
	var buffer bytes.Buffer

	buffer.WriteString(p.Ident)
	buffer.WriteString("(")

	for i, arg := range p.Args {
		buffer.WriteString(arg.(Value).CValue())

		if i != len(p.Args)-1 {
			buffer.WriteString(", ")
		}
	}

	buffer.WriteString(")")

	return buffer.String()
}

func (p ProcedureCall) CInstruction() string {
	return fmt.Sprintf("%s;", p.CValue())
}
