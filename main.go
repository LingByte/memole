package main

import (
	"fmt"
	"memmole/pkg/ast"
	"memmole/database"
	"memmole/pkg/evaluator"
	"memmole/pkg/lexer"
	"memmole/pkg/logger"
	"memmole/network"
	"memmole/pkg/parser"
	"memmole/repl"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// 初始化日志系统
	logger.SetLevel(logger.INFO)

	// 创建包管理器，设置基础路径为当前目录
	packageManager := parser.NewPackageManager(".")
	env := parser.NewEnvironment(nil)

	// 注册内置网络包
	networkPkg := network.NewNetworkPackage()
	env.Set("network", &network.NetworkObject{Package: networkPkg})

	// 注册内置数据库包
	dbPkg := database.NewDBPackage()
	env.Set("db", &database.DBObject{Package: dbPkg})

	// 支持从文件读取
	if len(os.Args) > 1 && strings.HasSuffix(os.Args[1], ".mml") {
		runFile(os.Args[1], env, packageManager)
		return
	}

	// 默认进入REPL
	repl.StartREPL(env, packageManager)
}

func runFile(filename string, env *parser.Environment, packageManager *parser.PackageManager) {
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

	// 创建包对象
	pkg := &parser.Package{
		Name:        filepath.Base(filename),
		Path:        filename,
		Statements:  program.Statements,
		Environment: env,
		Imports:     []*ast.ImportStatement{},
	}

	// 收集导入语句
	for _, stmt := range program.Statements {
		if importStmt, ok := stmt.(*ast.ImportStatement); ok {
			pkg.Imports = append(pkg.Imports, importStmt)
		}
	}

	// 解析导入
	if err := packageManager.ResolveImports(pkg); err != nil {
		fmt.Println("导入解析错误:", err)
		os.Exit(1)
	}

	// 执行整个程序
	logger.Info("开始执行程序...")
	result := evaluator.Eval(program, env)
	logger.Info("程序执行完成")

	// 输出最终返回值
	if result != nil {
		fmt.Println("返回值:", result.Inspect())
	} else {
		fmt.Println("返回值为null")
	}
}
