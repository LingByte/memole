package repl

import (
	"bufio"
	"fmt"
	"os"

	"github.com/LingByte/memole/pkg/evaluator"
	"github.com/LingByte/memole/pkg/lexer"
	"github.com/LingByte/memole/pkg/object"
	"github.com/LingByte/memole/pkg/parser"
)

// StartREPL 启动REPL
func StartREPL(env *parser.Environment) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("MemMole语言解释器，输入表达式：")

	for {
		fmt.Print(">> ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		runLine(line, env)
	}
}

// runLine 运行单行代码
func runLine(input string, env *parser.Environment) {
	env.Set("__exec_mode__", &object.String{Value: "repl"})
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
