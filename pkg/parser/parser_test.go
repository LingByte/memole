package parser

import (
	"memmole/pkg/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
        let x = 5;
        let y = 10;
        let foo = 20;
    `

	l := lexer.New(input)
	p, err := New(l)
	if err != nil {
		t.Fatalf("创建解析器失败: %v", err)
	}

	program := p.ParseProgram()
	if len(program.Statements) != 3 {
		t.Fatalf("应解析出3条语句，实际得到%d条", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foo"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("Token字面量应为'let'，得到=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*LetStatement)
	if !ok {
		t.Errorf("s 类型应为*LetStatement，实际=%T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("变量名应为%q，实际=%q", name, letStmt.Name.Value)
		return false
	}

	return true
}

func TestExpressionStatements(t *testing.T) {
	input := "3 + 5;"

	l := lexer.New(input)
	p, _ := New(l)
	program := p.ParseProgram()

	if len(program.Statements) != 1 {
		t.Fatalf("应解析出1条语句，实际得到%d条", len(program.Statements))
	}

	_, ok := program.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("语句类型应为ExpressionStatement，实际=%T", program.Statements[0])
	}
}

func TestOperatorPrecedence(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1 + 2 * 3", "(1 + (2 * 3))"},
		{"(1 + 2) * 3", "((1 + 2) * 3)"},
		{"-1 * 2", "((-1) * 2)"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p, _ := New(l)
		program := p.ParseProgram()

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("预期=%q, 实际=%q", tt.expected, actual)
		}
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"true", "true"},
		{"false", "false"},
		{"!true", "(!true)"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p, _ := New(l)
		program := p.ParseProgram()

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("预期=%q, 实际=%q", tt.expected, actual)
		}
	}
}

func TestComparisonExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1 > 2", "(1 > 2)"},
		{"1 < 2", "(1 < 2)"},
		{"1 == 2", "(1 == 2)"},
		{"1 != 2", "(1 != 2)"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p, _ := New(l)
		program := p.ParseProgram()

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("预期=%q, 实际=%q", tt.expected, actual)
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x > y) { x } else { y }`

	l := lexer.New(input)
	p, _ := New(l)
	program := p.ParseProgram()

	if len(program.Statements) != 1 {
		t.Fatalf("应解析出1条语句，实际得到%d条", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("语句类型应为ExpressionStatement，实际=%T", program.Statements[0])
	}

	ifExp, ok := stmt.Expression.(*IfExpression)
	if !ok {
		t.Fatalf("表达式类型应为IfExpression，实际=%T", stmt.Expression)
	}

	// 验证条件表达式
	if ifExp.Condition.String() != "(x > y)" {
		t.Errorf("条件表达式错误，期望=(x > y)，实际=%s", ifExp.Condition.String())
	}

	// 验证结果块
	if len(ifExp.Consequence.Statements) != 1 {
		t.Errorf("结果块应包含1条语句，实际=%d", len(ifExp.Consequence.Statements))
	}

	// 验证else块
	if ifExp.Alternative == nil {
		t.Error("应存在else分支")
	}

	// 验证整个结构
	expected := "if(x > y) x else y"
	if ifExp.String() != expected {
		t.Errorf("期望=%q, 实际=%q", expected, ifExp.String())
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	l := lexer.New(input)
	p, _ := New(l)
	program := p.ParseProgram()

	if len(program.Statements) != 1 {
		t.Fatalf("应解析出1条语句，实际得到%d条", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("语句类型应为ExpressionStatement，实际=%T", program.Statements[0])
	}

	function, ok := stmt.Expression.(*FunctionLiteral)
	if !ok {
		t.Fatalf("表达式类型应为FunctionLiteral，实际=%T", stmt.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("应解析出2个参数，实际得到%d个", len(function.Parameters))
	}

	expectedBody := "(x + y)"
	if function.Body.Statements[0].(*ExpressionStatement).Expression.String() != expectedBody {
		t.Errorf("函数体错误，期望=%q，实际=%q", expectedBody, function.Body.String())
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := `add(1, 2 * 3)`

	l := lexer.New(input)
	p, _ := New(l)
	program := p.ParseProgram()

	stmt := program.Statements[0].(*ExpressionStatement)
	exp, ok := stmt.Expression.(*CallExpression)
	if !ok {
		t.Fatalf("表达式类型应为CallExpression，实际=%T", stmt.Expression)
	}

	if exp.Function.String() != "add" {
		t.Errorf("函数名错误，期望=add，实际=%q", exp.Function.String())
	}

	if len(exp.Arguments) != 2 {
		t.Fatalf("应解析出2个参数，实际得到%d个", len(exp.Arguments))
	}

	expectedArgs := []string{"1", "(2 * 3)"}
	for i, arg := range exp.Arguments {
		if arg.String() != expectedArgs[i] {
			t.Errorf("参数%d错误，期望=%q，实际=%q", i, expectedArgs[i], arg.String())
		}
	}
}

func TestReturnStatements(t *testing.T) {
	input := `
        return 5;
        return add(10);
    `

	l := lexer.New(input)
	p, _ := New(l)
	program := p.ParseProgram()

	if len(program.Statements) != 2 {
		t.Fatalf("应解析出2条语句，实际得到%d条", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ReturnStatement)
		if !ok {
			t.Errorf("语句类型应为ReturnStatement，实际=%T", stmt)
			continue
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("Token字面量应为'return'，实际=%q", returnStmt.TokenLiteral())
		}
	}
}

func TestEnvironment(t *testing.T) {
	env := NewEnvironment(nil)

	// 测试变量存储
	env.Set("x", &Integer{Value: 10})
	val, exists := env.Get("x")
	if !exists {
		t.Fatal("变量x应存在")
	}
	if val.(*Integer).Value != 10 {
		t.Errorf("期望值=10，实际=%d", val.(*Integer).Value)
	}

	// 测试作用域链
	innerEnv := NewEnvironment(env)
	innerEnv.Set("y", &Integer{Value: 20})

	if _, exists := innerEnv.Get("x"); !exists {
		t.Error("应能访问外部作用域变量x")
	}
}
