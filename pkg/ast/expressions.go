package ast

import (
	"bytes"

	"github.com/LingByte/memole/pkg/lexer"
	"strings"
)

// Identifier 标识符结构
type Identifier struct {
	Token lexer.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// IntegerLiteral 整数字面量结构
type IntegerLiteral struct {
	Token lexer.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

// StringLiteral 字符串字面量结构
type StringLiteral struct {
	Token lexer.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

// Boolean 布尔字面量
type Boolean struct {
	Token lexer.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

// PrefixExpression 前缀表达式（用于!操作符）
type PrefixExpression struct {
	Token    lexer.Token // ! 或 -
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	return "(" + pe.Operator + pe.Right.String() + ")"
}

// InfixExpression 中缀表达式结构
type InfixExpression struct {
	Token    lexer.Token // 运算符token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	if ie.Left != nil {
		out.WriteString(ie.Left.String())
	}
	out.WriteString(" " + ie.Operator + " ")
	if ie.Right != nil {
		out.WriteString(ie.Right.String())
	}
	out.WriteString(")")
	return out.String()
}

// AssignmentExpression 赋值表达式结构
type AssignmentExpression struct {
	Token lexer.Token
	Name  *Identifier
	Value Expression
}

func (ae *AssignmentExpression) expressionNode()      {}
func (ae *AssignmentExpression) TokenLiteral() string { return ae.Token.Literal }
func (ae *AssignmentExpression) String() string {
	var out bytes.Buffer
	out.WriteString(ae.Name.String())
	out.WriteString(" = ")
	if ae.Value != nil {
		out.WriteString(ae.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// MemberAssignmentExpression 成员赋值表达式 obj.member = value
type MemberAssignmentExpression struct {
	Token  lexer.Token
	Object Expression
	Member *Identifier
	Value  Expression
}

func (mae *MemberAssignmentExpression) expressionNode()      {}
func (mae *MemberAssignmentExpression) TokenLiteral() string { return mae.Token.Literal }
func (mae *MemberAssignmentExpression) String() string {
	var out bytes.Buffer
	out.WriteString(mae.Object.String())
	out.WriteString(".")
	out.WriteString(mae.Member.String())
	out.WriteString(" = ")
	if mae.Value != nil {
		out.WriteString(mae.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// MemberAccessExpression 成员访问表达式结构
type MemberAccessExpression struct {
	Token  lexer.Token
	Object Expression
	Member *Identifier
}

func (mae *MemberAccessExpression) expressionNode()      {}
func (mae *MemberAccessExpression) TokenLiteral() string { return mae.Token.Literal }
func (mae *MemberAccessExpression) String() string {
	var out bytes.Buffer
	out.WriteString(mae.Object.String())
	out.WriteString(".")
	out.WriteString(mae.Member.String())
	return out.String()
}

// IfExpression if表达式
type IfExpression struct {
	Token       lexer.Token // if 关键字
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())
	if ie.Alternative != nil {
		out.WriteString(" else ")
		out.WriteString(ie.Alternative.String())
	}
	return out.String()
}

// FunctionLiteral 函数字面量
type FunctionLiteral struct {
	Token      lexer.Token
	Name       string
	ReturnType string
	Parameters []*TypedParameter
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := make([]string, len(fl.Parameters))
	for i, p := range fl.Parameters {
		params[i] = p.String()
	}
	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())
	return out.String()
}

// TypedParameter 带类型的参数
type TypedParameter struct {
	Type string
	Name *Identifier
}

func (tp *TypedParameter) String() string {
	return tp.Type + " " + tp.Name.String()
}

// CallExpression 函数调用表达式
type CallExpression struct {
	Token     lexer.Token // '('符号
	Function  Expression  // 标识符或函数字面量
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	args := make([]string, len(ce.Arguments))
	for i, a := range ce.Arguments {
		args[i] = a.String()
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}
