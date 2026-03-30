package ast

import (
	"bytes"
)

// Node AST节点的基础接口
type Node interface {
	TokenLiteral() string
	String() string
}

// Statement 语句接口
type Statement interface {
	Node
	statementNode()
}

// Expression 表达式接口
type Expression interface {
	Node
	expressionNode()
}

// Program 程序根节点
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}
