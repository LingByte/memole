package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/LingByte/memole/pkg/evaluator"
	"github.com/LingByte/memole/pkg/lexer"
	"github.com/LingByte/memole/pkg/object"
	"github.com/LingByte/memole/pkg/parser"
	"github.com/LingByte/memole/repl"
)

func main() {
	// 创建基础环境
	env := parser.NewEnvironment(nil)

	args := os.Args[1:]
	if len(args) == 0 {
		printUsage()
		return
	}

	switch args[0] {
	case "run":
		if len(args) < 2 {
			fmt.Println("错误: 缺少 .mml 文件路径")
			printUsage()
			os.Exit(1)
		}
		if !strings.HasSuffix(args[1], ".mml") {
			fmt.Println("错误: run 仅支持 .mml 文件")
			printUsage()
			os.Exit(1)
		}
		runFile(args[1], env)
		return
	case "repl":
		repl.StartREPL(env)
		return
	default:
		// 兼容直接传文件路径: ./memole app.mml
		if strings.HasSuffix(args[0], ".mml") {
			runFile(args[0], env)
			return
		}
		fmt.Printf("错误: 未知命令或参数 %q\n", args[0])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Memole 解释器")
	fmt.Println("")
	fmt.Println("用法:")
	fmt.Println("  memole run <file.mml>   执行 mml 文件")
	fmt.Println("  memole <file.mml>       执行 mml 文件（简写）")
	fmt.Println("  memole repl             进入交互模式")
}

func runFile(filename string, env *parser.Environment) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("读取文件出错:", err)
		os.Exit(1)
	}

	l := lexer.New(string(data))
	p, err := parser.New(l)
	if err != nil {
		fmt.Println("解析错误:", err)
		os.Exit(1)
	}

	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		fmt.Println("语法错误:")
		for _, msg := range p.Errors() {
			fmt.Println("\t", msg)
		}
		os.Exit(1)
	}

	// 执行整个程序
	env.Set("__exec_mode__", &object.String{Value: "file"})
	env.Set("__module_dir__", &object.String{Value: filepath.Dir(filename)})
	result := evaluator.Eval(program, env)

	// 只输出值本身，不添加前缀
	if result != nil && result.Type() != object.NULL_OBJ {
		fmt.Println(result.Inspect())
	}
}
