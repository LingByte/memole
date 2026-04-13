package parser

import (
	"fmt"

	"github.com/LingByte/memole/pkg/ast"
	"github.com/LingByte/memole/pkg/lexer"
)

// 定义优先级常量
const (
	_ int = iota
	LOWEST
	ASSIGNMENT // = 赋值运算符
	COMPARE    // > < == !=
	SUM        // + -
	PRODUCT    // * /
	PREFIX     // -X
	CALL       // 函数调用
	MEMBER     // . 成员访问
)

// 定义优先级表
var precedences = map[lexer.TokenType]int{
	lexer.Assign:   ASSIGNMENT,
	lexer.Plus:     SUM,
	lexer.Minus:    SUM,
	lexer.Multiply: PRODUCT,
	lexer.Slash:    PRODUCT,
	lexer.Dot:      MEMBER,
	lexer.GT:       COMPARE,
	lexer.LT:       COMPARE,
	lexer.EQ:       COMPARE,
	lexer.NotEQ:    COMPARE,
	lexer.LParen:   CALL,
}

// 解析函数类型
type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// Parser 解析器结构体
type Parser struct {
	l              *lexer.Lexer
	curToken       lexer.Token
	peekToken      lexer.Token
	prefixParseFns map[lexer.TokenType]prefixParseFn
	infixParseFns  map[lexer.TokenType]infixParseFn
	errors         []string
}

// New 创建Parser实例
func New(l *lexer.Lexer) (*Parser, error) {
	p := &Parser{
		l:              l,
		prefixParseFns: make(map[lexer.TokenType]prefixParseFn),
		infixParseFns:  make(map[lexer.TokenType]infixParseFn),
	}

	// 注册前缀解析函数
	p.registerPrefix(lexer.Int, p.parseIntegerLiteral)
	p.registerPrefix(lexer.String, p.parseStringLiteral)
	p.registerPrefix(lexer.Identifier, p.parseIdentifier)
	p.registerPrefix(lexer.True, p.parseBoolean)
	p.registerPrefix(lexer.False, p.parseBoolean)
	p.registerPrefix(lexer.Bang, p.parsePrefixExpression)
	p.registerPrefix(lexer.If, p.parseIfExpression)
	p.registerPrefix(lexer.Function, p.parseFunctionLiteral)

	// 注册中缀解析函数
	p.registerInfix(lexer.Assign, p.parseAssignmentExpression)
	p.registerInfix(lexer.Dot, p.parseMemberAccessExpression)
	p.registerInfix(lexer.Plus, p.parseInfixExpression)
	p.registerInfix(lexer.Minus, p.parseInfixExpression)
	p.registerInfix(lexer.Multiply, p.parseInfixExpression)
	p.registerInfix(lexer.Slash, p.parseInfixExpression)
	p.registerInfix(lexer.GT, p.parseInfixExpression)
	p.registerInfix(lexer.LT, p.parseInfixExpression)
	p.registerInfix(lexer.EQ, p.parseInfixExpression)
	p.registerInfix(lexer.NotEQ, p.parseInfixExpression)
	p.registerInfix(lexer.LParen, p.parseCallExpression)

	// 读取两个Token初始化状态
	p.nextToken()
	p.nextToken()

	if p.curToken.Type == lexer.Illegal {
		return nil, fmt.Errorf("非法Token: %s", p.curToken.Literal)
	}
	return p, nil
}

// nextToken 推进Token指针
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// ParseProgram 解析整个程序
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}

	for p.curToken.Type != lexer.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

// parseStatement 解析语句
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Literal {
	case "let":
		return p.parseLetStatement()
	case "return":
		return p.parseReturnStatement()
	case "package":
		return p.parsePackageStatement()
	case "import":
		return p.parseImportStatement()
	case "while":
		return p.parseWhileStatement()
	case "ty":
		return p.parseTypeStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// parseExpression 解析表达式
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(lexer.Semicolon) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

// 注册函数方法
func (p *Parser) registerPrefix(tokenType lexer.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType lexer.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// 辅助方法
func (p *Parser) curTokenIs(t lexer.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t lexer.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t lexer.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekError(t lexer.TokenType) {
	msg := fmt.Sprintf("预期下一个Token是 %s, 但得到 %s",
		t.String(), p.peekToken.Type.String())
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

// Errors 获取错误列表
func (p *Parser) Errors() []string {
	return p.errors
}
