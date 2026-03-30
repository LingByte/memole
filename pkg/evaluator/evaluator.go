package evaluator

import (
	"fmt"
	"memmole/pkg/ast"
	"memmole/pkg/object"
	"memmole/pkg/parser"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

// FunctionObject 简化的函数对象
type FunctionObject struct {
	Parameters []*ast.TypedParameter
	Body       *ast.BlockStatement
	Env        *parser.Environment
}

func (f *FunctionObject) Type() object.ObjectType { return "FUNCTION" }
func (f *FunctionObject) Inspect() string {
	return fmt.Sprintf("function with %d parameters", len(f.Parameters))
}

// Eval 主求值函数
func Eval(node ast.Node, env *parser.Environment) object.Object {
	switch node := node.(type) {

	case *ast.Program:
		return evalProgram(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.PackageStatement:
		return evalPackageStatement(node, env)

	case *ast.ImportStatement:
		return evalImportStatement(node, env)

	case *ast.WhileStatement:
		return evalWhileStatement(node, env)

	case *ast.TypeStatement:
		return evalTypeStatement(node, env)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		right := Eval(node.Right, env)
		return evalInfixExpression(node.Operator, left, right)

	case *ast.AssignmentExpression:
		return evalAssignmentExpression(node, env)

	case *ast.MemberAccessExpression:
		return evalMemberAccessExpression(node, env)

	case *ast.MemberAssignmentExpression:
		return evalMemberAssignmentExpression(node, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		env.Set(node.Name.Value, val)
		return val

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		return &object.ReturnValue{Value: val}

	case *ast.FunctionLiteral:
		// 暂时创建一个简化的函数对象
		fn := &FunctionObject{
			Parameters: node.Parameters,
			Body:       node.Body,
			Env:        env,
		}
		if node.Name != "" {
			env.Set(node.Name, fn)
		}
		return fn

	case *ast.CallExpression:
		return evalCallExpression(node, env)
	}

	return NULL
}

// evalProgram 求值程序
func evalProgram(program *ast.Program, env *parser.Environment) object.Object {
	for _, stmt := range program.Statements {
		Eval(stmt, env)
	}

	mainFunc, ok := env.Get("main")
	if !ok {
		fmt.Println("未找到 main 函数")
		return NULL
	}

	fnObj, ok := mainFunc.(*FunctionObject)
	if !ok {
		fmt.Println("main 不是一个函数")
		return NULL
	}

	return applyFunction(fnObj, []object.Object{})
}

// evalBlockStatement 求值代码块
func evalBlockStatement(block *ast.BlockStatement, env *parser.Environment) object.Object {
	var result object.Object

	for _, stmt := range block.Statements {
		result = Eval(stmt, env)

		// 注释掉详细的执行语句日志
		// fmt.Printf("执行语句: %s => %s\n", stmt.String(), result.Inspect())

		if returnObj, ok := result.(*object.ReturnValue); ok {
			return returnObj
		}
	}
	return result
}

// applyFunction 应用函数
func applyFunction(fn *FunctionObject, args []object.Object) object.Object {
	// fmt.Printf("执行函数: %s\n", fn.Inspect())
	extendedEnv := extendFunctionEnv(fn, args)
	evaluated := Eval(fn.Body, extendedEnv)

	if returnValue, ok := evaluated.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return evaluated
}

// extendFunctionEnv 扩展函数环境
func extendFunctionEnv(fn *FunctionObject, args []object.Object) *parser.Environment {
	env := parser.NewEnvironment(fn.Env)

	for i, param := range fn.Parameters {
		env.Set(param.Name.Value, args[i])
	}

	return env
}

// nativeBoolToBooleanObject 转换原生布尔值到布尔对象
func nativeBoolToBooleanObject(input bool) object.Object {
	if input {
		return TRUE
	}
	return FALSE
}

// isTruthy 判断对象是否为真
func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL, FALSE:
		return false
	default:
		return true
	}
}
