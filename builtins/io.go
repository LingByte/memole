package builtins

import (
	"memmole/pkg/ast"
	"memmole/pkg/parser"
	"memmole/pkg/object"
	"fmt"
)

// EvalIOCall 处理IO包调用
func EvalIOCall(exp *ast.MemberAccessExpression, args []object.Object, env *parser.Environment) object.Object {
	switch exp.Member.Value {
	case "print":
		return evalPrint(args)
	case "println":
		return evalPrintln(args)
	case "printf":
		return evalPrintf(args)
	default:
		return &object.String{Value: fmt.Sprintf("未知的IO方法: %s", exp.Member.Value)}
	}
}

// evalPrint 打印函数（不换行）
func evalPrint(args []object.Object) object.Object {
	if len(args) == 0 {
		return &object.String{Value: "print需要至少1个参数"}
	}

	// 将所有参数转换为字符串并打印
	for i, arg := range args {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(arg.Inspect())
	}

	return &object.String{Value: "print完成"}
}

// evalPrintln 打印函数（换行）
func evalPrintln(args []object.Object) object.Object {
	if len(args) == 0 {
		fmt.Println()
		return &object.String{Value: "println完成"}
	}

	// 将所有参数转换为字符串并打印
	for i, arg := range args {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(arg.Inspect())
	}
	fmt.Println()

	return &object.String{Value: "println完成"}
}

// evalPrintf 格式化打印函数
func evalPrintf(args []object.Object) object.Object {
	if len(args) < 2 {
		return &object.String{Value: "printf需要至少2个参数: format, values..."}
	}

	// 第一个参数是格式字符串
	format, ok := args[0].(*object.String)
	if !ok {
		return &object.String{Value: "printf第一个参数必须是字符串"}
	}

	// 将其他参数转换为interface{}切片
	values := make([]interface{}, len(args)-1)
	for i, arg := range args[1:] {
		switch v := arg.(type) {
		case *object.Integer:
			values[i] = v.Value
		case *object.String:
			values[i] = v.Value
		case *object.Boolean:
			values[i] = v.Value
		default:
			values[i] = arg.Inspect()
		}
	}

	fmt.Printf(format.Value, values...)
	return &object.String{Value: "printf完成"}
}
