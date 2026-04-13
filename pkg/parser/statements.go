package parser

import (
	"github.com/LingByte/memole/pkg/ast"
	"github.com/LingByte/memole/pkg/lexer"
)

// parseLetStatement 解析变量声明
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(lexer.Identifier) {
		return nil
	}

	stmt.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	if !p.expectPeek(lexer.Assign) {
		return nil
	}

	p.nextToken() // 跳到表达式起始位置
	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(lexer.Semicolon) {
		p.nextToken()
	}

	return stmt
}

// parseReturnStatement 解析返回语句
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()
	stmt.ReturnValue = p.parseExpression(LOWEST)

	// 分号可选：遇到 ';' 消费掉，遇到 EOF 直接结束，避免死循环。
	if p.peekTokenIs(lexer.Semicolon) {
		p.nextToken()
	}

	return stmt
}

// parsePackageStatement 解析包声明
func (p *Parser) parsePackageStatement() *ast.PackageStatement {
	stmt := &ast.PackageStatement{Token: p.curToken}

	if !p.expectPeek(lexer.Identifier) {
		return nil
	}

	stmt.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	if p.peekTokenIs(lexer.Semicolon) {
		p.nextToken()
	}

	return stmt
}

// parseImportStatement 解析导入语句
func (p *Parser) parseImportStatement() *ast.ImportStatement {
	stmt := &ast.ImportStatement{Token: p.curToken}

	if !p.expectPeek(lexer.Identifier) {
		return nil
	}

	stmt.Path = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	// 检查是否有别名 (import math as m)
	if p.peekTokenIs(lexer.Identifier) && p.peekToken.Literal == "as" {
		p.nextToken() // 跳过 "as"
		if !p.expectPeek(lexer.Identifier) {
			return nil
		}
		stmt.Alias = &ast.Identifier{
			Token: p.curToken,
			Value: p.curToken.Literal,
		}
	}

	if p.peekTokenIs(lexer.Semicolon) {
		p.nextToken()
	}

	return stmt
}

// parseWhileStatement 解析while循环语句
func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	stmt := &ast.WhileStatement{Token: p.curToken}

	if !p.expectPeek(lexer.LParen) {
		return nil
	}

	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.RParen) {
		return nil
	}

	if !p.expectPeek(lexer.Lbrace) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

// parseExpressionStatement 解析表达式语句
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(lexer.Semicolon) {
		p.nextToken()
	}
	return stmt
}

// parseBlockStatement 解析代码块语句
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	p.nextToken()

	for !p.curTokenIs(lexer.Rbrace) && !p.curTokenIs(lexer.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}

// parseTypeStatement 解析类型声明： ty Name stru { name type, ... }
func (p *Parser) parseTypeStatement() *ast.TypeStatement {
	stmt := &ast.TypeStatement{Token: p.curToken}

	if !p.expectPeek(lexer.Identifier) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// 期望关键字 stru
	if !p.expectPeek(lexer.Stru) {
		return nil
	}
	stmt.Kind = "stru"

	if !p.expectPeek(lexer.Lbrace) {
		return nil
	}

	fields := []*ast.StructField{}

	for !p.peekTokenIs(lexer.Rbrace) {
		// 字段名
		p.nextToken()
		if p.curToken.Type != lexer.Identifier {
			return nil
		}
		fname := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		// 字段类型
		if !p.expectPeek(lexer.Identifier) {
			return nil
		}
		ftype := p.curToken.Literal
		fields = append(fields, &ast.StructField{Name: fname, Type: ftype})
		if p.peekTokenIs(lexer.Comma) || p.peekTokenIs(lexer.Semicolon) {
			p.nextToken()
		}
	}

	if !p.expectPeek(lexer.Rbrace) {
		return nil
	}

	stmt.Fields = fields
	return stmt
}
