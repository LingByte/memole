package parser

import (
	"memmole/pkg/ast"
	"memmole/pkg/lexer"
	"fmt"
	"strconv"
)

// parseIntegerLiteral 解析整数字面量
func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		p.errors = append(p.errors, fmt.Sprintf("无法解析 %q 为整数", p.curToken.Literal))
		return nil
	}

	lit.Value = value
	return lit
}

// parseStringLiteral 解析字符串字面量
func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

// parseIdentifier 解析标识符
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

// parseBoolean 解析布尔字面量
func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{
		Token: p.curToken,
		Value: p.curToken.Literal == "true",
	}
}

// parsePrefixExpression 解析前缀表达式
func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

// parseInfixExpression 解析中缀表达式
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

// parseAssignmentExpression 解析赋值表达式
func (p *Parser) parseAssignmentExpression(left ast.Expression) ast.Expression {
	// 标识符赋值
	if ident, ok := left.(*ast.Identifier); ok {
		expression := &ast.AssignmentExpression{
			Token: p.curToken,
			Name:  ident,
		}

		precedence := p.curPrecedence()
		p.nextToken()
		expression.Value = p.parseExpression(precedence)

		return expression
	}

	// 成员赋值：obj.member = expr
	if mem, ok := left.(*ast.MemberAccessExpression); ok {
		expression := &ast.MemberAssignmentExpression{
			Token:  p.curToken,
			Object: mem.Object,
			Member: mem.Member,
		}
		precedence := p.curPrecedence()
		p.nextToken()
		expression.Value = p.parseExpression(precedence)
		return expression
	}

	p.errors = append(p.errors, fmt.Sprintf("赋值左侧必须是标识符，得到 %T", left))
	return nil
}

// parseMemberAccessExpression 解析成员访问表达式
func (p *Parser) parseMemberAccessExpression(left ast.Expression) ast.Expression {
	expression := &ast.MemberAccessExpression{
		Token:  p.curToken,
		Object: left,
	}

	// 读取成员名
	if !p.expectPeek(lexer.Identifier) {
		return nil
	}

	expression.Member = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	return expression
}

// parseIfExpression 解析if表达式
func (p *Parser) parseIfExpression() ast.Expression {
	expr := &ast.IfExpression{Token: p.curToken}

	if !p.expectPeek(lexer.LParen) {
		return nil
	}

	p.nextToken()
	expr.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.RParen) {
		return nil
	}

	if !p.expectPeek(lexer.Lbrace) {
		return nil
	}

	expr.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(lexer.Else) {
		p.nextToken()
		if !p.expectPeek(lexer.Lbrace) {
			return nil
		}
		expr.Alternative = p.parseBlockStatement()
	}

	return expr
}

// parseFunctionLiteral 解析函数字面量
func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	// 支持三种形式：
	// 1) fn(name-only): fn main(...) { }
	// 2) fn ReturnType Name(Type x, ... ) { }
	// 3) fn(...): 匿名无类型参数
	if p.peekTokenIs(lexer.LParen) {
		// 匿名无类型参数 fn(...)
		if !p.expectPeek(lexer.LParen) {
			return nil
		}
		lit.Parameters = p.parseUntypedFunctionParameters()
	} else if p.peekTokenIs(lexer.Identifier) {
		// 读取第一个标识符，可能是 函数名 或 返回类型
		if !p.expectPeek(lexer.Identifier) {
			return nil
		}
		first := p.curToken.Literal

		if p.peekTokenIs(lexer.LParen) {
			// 形如 fn Name(...)
			lit.Name = first
			if !p.expectPeek(lexer.LParen) {
				return nil
			}
			lit.Parameters = p.parseUntypedFunctionParameters()
		} else {
			// 形如 fn ReturnType Name(...)
			lit.ReturnType = first
			if !p.expectPeek(lexer.Identifier) {
				return nil
			}
			lit.Name = p.curToken.Literal
			if !p.expectPeek(lexer.LParen) {
				return nil
			}
			lit.Parameters = p.parseTypedFunctionParameters()
		}
	} else {
		return nil
	}

	// 读取函数体
	if !p.expectPeek(lexer.Lbrace) {
		return nil
	}
	lit.Body = p.parseBlockStatement()

	return lit
}

// parseTypedFunctionParameters 解析带类型的函数参数
func (p *Parser) parseTypedFunctionParameters() []*ast.TypedParameter {
	params := []*ast.TypedParameter{}

	if p.peekTokenIs(lexer.RParen) {
		p.nextToken()
		return params
	}

	p.nextToken() // 读取类型
	paramType := p.curToken.Literal

	if !p.expectPeek(lexer.Identifier) {
		return nil
	}
	paramName := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	params = append(params, &ast.TypedParameter{Type: paramType, Name: paramName})

	for p.peekTokenIs(lexer.Comma) {
		p.nextToken() // 跳过逗号
		p.nextToken() // 读取下一个类型
		ptype := p.curToken.Literal

		if !p.expectPeek(lexer.Identifier) {
			return nil
		}
		pname := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		params = append(params, &ast.TypedParameter{Type: ptype, Name: pname})
	}

	if !p.expectPeek(lexer.RParen) {
		return nil
	}

	return params
}

// parseUntypedFunctionParameters 解析无类型参数 (x, y)
func (p *Parser) parseUntypedFunctionParameters() []*ast.TypedParameter {
	params := []*ast.TypedParameter{}

	if p.peekTokenIs(lexer.RParen) {
		p.nextToken()
		return params
	}

	p.nextToken()
	if p.curToken.Type != lexer.Identifier {
		return nil
	}
	first := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	params = append(params, &ast.TypedParameter{Type: "", Name: first})

	for p.peekTokenIs(lexer.Comma) {
		p.nextToken()
		p.nextToken()
		if p.curToken.Type != lexer.Identifier {
			return nil
		}
		name := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		params = append(params, &ast.TypedParameter{Type: "", Name: name})
	}

	if !p.expectPeek(lexer.RParen) {
		return nil
	}

	return params
}

// parseCallExpression 解析函数调用表达式
func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{
		Token:    p.curToken,
		Function: function,
	}
	exp.Arguments = p.parseExpressionList(lexer.RParen)
	return exp
}

// parseExpressionList 解析表达式列表
func (p *Parser) parseExpressionList(end lexer.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(lexer.Comma) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}
