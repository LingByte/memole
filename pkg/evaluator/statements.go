package evaluator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/LingByte/memole/pkg/ast"
	"github.com/LingByte/memole/pkg/builtinmodules"
	"github.com/LingByte/memole/pkg/object"
	"github.com/LingByte/memole/pkg/parser"
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
	bindName := stmt.Path.Value
	if stmt.Alias != nil {
		bindName = stmt.Alias.Value
	}

	// 先查官方内置模块，再回退到 .mml 文件模块
	if builtinMod, ok := builtinmodules.Get(stmt.Path.Value); ok {
		env.Set(bindName, builtinMod)
		return builtinMod
	}

	baseDir, err := os.Getwd()
	if err != nil {
		baseDir = "."
	}
	if raw, ok := env.Get("__module_dir__"); ok {
		if s, ok := raw.(*object.String); ok && s.Value != "" {
			baseDir = s.Value
		}
	}

	mod, loadErr := defaultModLD.load(stmt.Path.Value, filepath.Clean(baseDir))
	if loadErr != nil {
		fmt.Printf("模块导入失败 %s: %v\n", stmt.Path.Value, loadErr)
		return NULL
	}

	env.Set(bindName, mod)
	return mod
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
