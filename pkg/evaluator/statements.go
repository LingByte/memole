package evaluator

import (
	"memmole/pkg/ast"
	"memmole/pkg/parser"
	"memmole/pkg/object"
)

// evalPackageStatement 求值包声明语句
func evalPackageStatement(stmt *ast.PackageStatement, env *parser.Environment) object.Object {
	// 包声明主要用于标识，不需要特殊处理
	// 可以在环境中记录包名
	env.Set("__package__", &object.String{Value: stmt.Name.Value})
	return NULL
}

// evalImportStatement 求值导入语句
func evalImportStatement(stmt *ast.ImportStatement, env *parser.Environment) object.Object {
	// 导入语句在包加载时已经处理，这里不需要额外处理
	return NULL
}

// evalWhileStatement 求值while循环语句
func evalWhileStatement(stmt *ast.WhileStatement, env *parser.Environment) object.Object {
	var result object.Object

	for {
		condition := Eval(stmt.Condition, env)
		if !isTruthy(condition) {
			break
		}

		result = Eval(stmt.Body, env)

		// 检查是否有return语句
		if returnObj, ok := result.(*object.ReturnValue); ok {
			return returnObj
		}
	}

	return result
}

// evalTypeStatement 注册结构体类型
func evalTypeStatement(stmt *ast.TypeStatement, env *parser.Environment) object.Object {
    if stmt.Kind != "stru" {
        return NULL
    }
    fieldNames := make([]string, 0, len(stmt.Fields))
    typesMap := make(map[string]string, len(stmt.Fields))
    for _, f := range stmt.Fields {
        fieldNames = append(fieldNames, f.Name.Value)
        typesMap[f.Name.Value] = f.Type
    }
    st := &object.StructType{Name: stmt.Name.Value, Fields: fieldNames, Types: typesMap}
    env.Set(stmt.Name.Value, st)
    return NULL
}