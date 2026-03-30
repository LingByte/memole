package repl

import (
	"bufio"
	"memmole/pkg/evaluator"
	"memmole/pkg/lexer"
	"memmole/pkg/parser"
	"fmt"
	"os"
)

// StartREPL 启动REPL
func StartREPL(env *parser.Environment, packageManager *parser.PackageManager) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("CVM语言解释器（支持包管理），输入表达式：")

	for {
		fmt.Print(">> ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		runLine(line, env, packageManager)
	}
}

// runLine 运行单行代码
func runLine(input string, env *parser.Environment, packageManager *parser.PackageManager) {
	l := lexer.New(input)
	p, err := parser.New(l)
	if err != nil {
		fmt.Println("解析错误:", err)
		return
	}
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		fmt.Println("语法错误:")
		for _, msg := range p.Errors() {
			fmt.Println("\t", msg)
		}
		return
	}
	result := evaluator.Eval(program, env)
	if result != nil {
		fmt.Println(result.Inspect())
	}
}
