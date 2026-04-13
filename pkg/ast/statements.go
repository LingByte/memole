package ast

import (
	"bytes"

	"github.com/LingByte/memole/pkg/lexer"
)

// LetStatement let语句结构
type LetStatement struct {
	Token lexer.Token // let Token
	Name  *Identifier // 变量名
	Value Expression  // 表达式
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// PackageStatement 包声明语句
type PackageStatement struct {
	Token lexer.Token
	Name  *Identifier
}

func (ps *PackageStatement) statementNode()       {}
func (ps *PackageStatement) TokenLiteral() string { return ps.Token.Literal }
func (ps *PackageStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ps.TokenLiteral() + " ")
	out.WriteString(ps.Name.String())
	out.WriteString(";")
	return out.String()
}

// ImportStatement 导入语句
type ImportStatement struct {
	Token lexer.Token
	Path  *Identifier
	Alias *Identifier // 可选的别名
}

func (is *ImportStatement) statementNode()       {}
func (is *ImportStatement) TokenLiteral() string { return is.Token.Literal }
func (is *ImportStatement) String() string {
	var out bytes.Buffer
	out.WriteString(is.TokenLiteral() + " ")
	out.WriteString(is.Path.String())
	if is.Alias != nil {
		out.WriteString(" as ")
		out.WriteString(is.Alias.String())
	}
	out.WriteString(";")
	return out.String()
}

// WhileStatement while循环语句
type WhileStatement struct {
	Token     lexer.Token
	Condition Expression
	Body      *BlockStatement
}

func (ws *WhileStatement) statementNode()       {}
func (ws *WhileStatement) TokenLiteral() string { return ws.Token.Literal }
func (ws *WhileStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ws.TokenLiteral())
	out.WriteString("(")
	out.WriteString(ws.Condition.String())
	out.WriteString(") ")
	out.WriteString(ws.Body.String())
	return out.String()
}

// ExpressionStatement 表达式语句
type ExpressionStatement struct {
	Token      lexer.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// ReturnStatement 返回语句
type ReturnStatement struct {
	Token       lexer.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

// BlockStatement 代码块语句
type BlockStatement struct {
	Token      lexer.Token // { 符号
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// TypeStatement 类型声明（结构体）
type TypeStatement struct {
	Token  lexer.Token // ty
	Name   *Identifier
	Kind   string // "stru"
	Fields []*StructField
}

func (ts *TypeStatement) statementNode()       {}
func (ts *TypeStatement) TokenLiteral() string { return ts.Token.Literal }
func (ts *TypeStatement) String() string {
	var out bytes.Buffer
	out.WriteString("ty ")
	out.WriteString(ts.Name.String())
	out.WriteString(" stru {")
	for i, f := range ts.Fields {
		if i > 0 {
			out.WriteString(", ")
		}
		out.WriteString(f.String())
	}
	out.WriteString("}")
	return out.String()
}

// StructField 结构体字段
type StructField struct {
	Name *Identifier
	Type string
}

func (sf *StructField) String() string {
	return sf.Name.String() + " " + sf.Type
}
