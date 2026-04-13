package evaluator

import (
	"github.com/LingByte/memole/pkg/ast"
	"github.com/LingByte/memole/pkg/object"
	"github.com/LingByte/memole/pkg/parser"
)

// evalPrefixExpression 求值前缀表达式
func evalPrefixExpression(op string, right object.Object) object.Object {
	switch op {
	case "!":
		return evalBangOperator(right)
	case "-":
		if val, ok := right.(*object.Integer); ok {
			return &object.Integer{Value: -val.Value}
		}
	}
	return NULL
}

// evalBangOperator 求值!操作符
func evalBangOperator(obj object.Object) object.Object {
	switch obj {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

// evalInfixExpression 求值中缀表达式
func evalInfixExpression(op string, left, right object.Object) object.Object {
	lval, lok := left.(*object.Integer)
	rval, rok := right.(*object.Integer)

	if lok && rok {
		switch op {
		case "+":
			return &object.Integer{Value: lval.Value + rval.Value}
		case "-":
			return &object.Integer{Value: lval.Value - rval.Value}
		case "*":
			return &object.Integer{Value: lval.Value * rval.Value}
		case "/":
			return &object.Integer{Value: lval.Value / rval.Value}
		case "==":
			return nativeBoolToBooleanObject(lval.Value == rval.Value)
		case "!=":
			return nativeBoolToBooleanObject(lval.Value != rval.Value)
		case ">":
			return nativeBoolToBooleanObject(lval.Value > rval.Value)
		case "<":
			return nativeBoolToBooleanObject(lval.Value < rval.Value)
		}
	}

	// 比较布尔对象
	if lb, ok1 := left.(*object.Boolean); ok1 {
		if rb, ok2 := right.(*object.Boolean); ok2 {
			switch op {
			case "==":
				return nativeBoolToBooleanObject(lb.Value == rb.Value)
			case "!=":
				return nativeBoolToBooleanObject(lb.Value != rb.Value)
			}
		}
	}

	// 比较字符串对象
	if ls, ok1 := left.(*object.String); ok1 {
		if rs, ok2 := right.(*object.String); ok2 {
			switch op {
			case "==":
				return nativeBoolToBooleanObject(ls.Value == rs.Value)
			case "!=":
				return nativeBoolToBooleanObject(ls.Value != rs.Value)
			}
		}
	}

	return NULL
}

// evalIfExpression 求值if表达式
func evalIfExpression(expr *ast.IfExpression, env *parser.Environment) object.Object {
	condition := Eval(expr.Condition, env)
	if isTruthy(condition) {
		return Eval(expr.Consequence, env)
	} else if expr.Alternative != nil {
		return Eval(expr.Alternative, env)
	}
	return NULL
}

// evalIdentifier 求值标识符
func evalIdentifier(node *ast.Identifier, env *parser.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	return NULL
}

// evalAssignmentExpression 求值赋值表达式
func evalAssignmentExpression(exp *ast.AssignmentExpression, env *parser.Environment) object.Object {
	value := Eval(exp.Value, env)
	env.Set(exp.Name.Value, value)
	return value
}

// evalMemberAssignmentExpression 成员赋值
func evalMemberAssignmentExpression(exp *ast.MemberAssignmentExpression, env *parser.Environment) object.Object {
	obj := Eval(exp.Object, env)
	val := Eval(exp.Value, env)
	if inst, ok := obj.(*object.StructInstance); ok {
		if inst.Fields == nil {
			inst.Fields = map[string]object.Object{}
		}
		inst.Fields[exp.Member.Value] = val
		return val
	}
	return NULL
}

// evalMemberAccessExpression 求值成员访问表达式
func evalMemberAccessExpression(exp *ast.MemberAccessExpression, env *parser.Environment) object.Object {
	val := Eval(exp.Object, env)

	// 检查是否是包名（内置包延后在调用处理）
	if ident, ok := exp.Object.(*ast.Identifier); ok {
		switch ident.Value {
		case "network", "io", "log", "config", "db":
			return NULL
		}
	}

	// 结构体实例：字段或绑定方法
	if inst, ok := val.(*object.StructInstance); ok {
		if v, exists := inst.Fields[exp.Member.Value]; exists {
			return v
		}
		return &object.BoundMethod{TypeName: inst.TypeName, Method: exp.Member.Value, Self: inst}
	}

	// 结构体类型：类型方法绑定（无 self 注入）
	if st, ok := val.(*object.StructType); ok {
		return &object.BoundMethod{TypeName: st.Name, Method: exp.Member.Value, Self: nil}
	}

	// 模块对象：从导出符号表中读取成员
	if mod, ok := val.(*object.ModuleObject); ok {
		if exported, exists := mod.Exports[exp.Member.Value]; exists {
			return exported
		}
		return NULL
	}

	return NULL
}

// evalCallExpression 求值函数调用表达式
func evalCallExpression(exp *ast.CallExpression, env *parser.Environment) object.Object {
	function := Eval(exp.Function, env)
	args := []object.Object{}
	for _, arg := range exp.Arguments {
		args = append(args, Eval(arg, env))
	}

	// 结构体构造：User("Alice", 1)
	if st, ok := function.(*object.StructType); ok {
		inst := &object.StructInstance{TypeName: st.Name, Fields: map[string]object.Object{}}
		for i, fname := range st.Fields {
			if i < len(args) {
				inst.Fields[fname] = args[i]
			} else {
				inst.Fields[fname] = NULL
			}
		}
		return inst
	}

	// 绑定方法：转为调用全局函数 Type_Method(self, ...)
	if bm, ok := function.(*object.BoundMethod); ok {
		globalName := bm.TypeName + "_" + bm.Method
		if any, ok := env.Get(globalName); ok {
			if fn, ok := any.(*FunctionObject); ok {
				if bm.Self != nil {
					return applyFunction(fn, append([]object.Object{bm.Self}, args...))
				}
				return applyFunction(fn, args)
			}
		}
		return NULL
	}

	if fn, ok := function.(*FunctionObject); ok {
		return applyFunction(fn, args)
	}

	if nf, ok := function.(*object.NativeFunction); ok {
		return nf.Fn(args)
	}

	return NULL
}
