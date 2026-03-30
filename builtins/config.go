package builtins

import (
	"memmole/pkg/ast"
	"memmole/pkg/config"
	"memmole/pkg/object"
	"memmole/pkg/parser"
	"fmt"
)

// EvalConfigCall 处理config包调用
func EvalConfigCall(exp *ast.MemberAccessExpression, args []object.Object, env *parser.Environment) object.Object {
	switch exp.Member.Value {
	case "load":
		return evalConfigLoad(args)
	case "load_from_file":
		return evalConfigLoadFromFile(args)
	case "get":
		return evalConfigGet(args)
	case "set":
		return evalConfigSet(args)
	case "get_int":
		return evalConfigGetInt(args)
	case "set_int":
		return evalConfigSetInt(args)
	case "get_float":
		return evalConfigGetFloat(args)
	case "set_float":
		return evalConfigSetFloat(args)
	case "get_bool":
		return evalConfigGetBool(args)
	case "set_bool":
		return evalConfigSetBool(args)
	case "has":
		return evalConfigHas(args)
	case "delete":
		return evalConfigDelete(args)
	case "get_all":
		return evalConfigGetAll(args)
	case "save_to_file":
		return evalConfigSaveToFile(args)
	case "clear":
		return evalConfigClear(args)
	default:
		return &object.String{Value: fmt.Sprintf("未知的config方法: %s", exp.Member.Value)}
	}
}

// evalConfigLoad 加载默认配置文件
func evalConfigLoad(args []object.Object) object.Object {
	if len(args) != 0 {
		return &object.String{Value: "load不需要参数"}
	}
	
	err := config.LoadConfig()
	if err != nil {
		return &object.String{Value: fmt.Sprintf("加载配置文件失败: %v", err)}
	}
	
	return &object.String{Value: "配置文件加载成功"}
}

// evalConfigLoadFromFile 从指定文件加载配置
func evalConfigLoadFromFile(args []object.Object) object.Object {
	if len(args) != 1 {
		return &object.String{Value: "load_from_file需要1个参数: filename"}
	}
	
	filename, ok := args[0].(*object.String)
	if !ok {
		return &object.String{Value: "参数类型错误，应为 string"}
	}
	
	err := config.LoadConfigFromFile(filename.Value)
	if err != nil {
		return &object.String{Value: fmt.Sprintf("加载配置文件失败: %v", err)}
	}
	
	return &object.String{Value: "配置文件加载成功"}
}

// evalConfigGet 获取配置项
func evalConfigGet(args []object.Object) object.Object {
	if len(args) != 1 {
		return &object.String{Value: "get需要1个参数: key"}
	}
	
	key, ok := args[0].(*object.String)
	if !ok {
		return &object.String{Value: "参数类型错误，应为 string"}
	}
	
	value := config.Get(key.Value)
	return &object.String{Value: value}
}

// evalConfigSet 设置配置项
func evalConfigSet(args []object.Object) object.Object {
	if len(args) != 2 {
		return &object.String{Value: "set需要2个参数: key, value"}
	}
	
	key, ok1 := args[0].(*object.String)
	value, ok2 := args[1].(*object.String)
	if !ok1 || !ok2 {
		return &object.String{Value: "参数类型错误，应为 (string, string)"}
	}
	
	config.Set(key.Value, value.Value)
	return &object.String{Value: "配置项设置成功"}
}

// evalConfigGetInt 获取整型配置项
func evalConfigGetInt(args []object.Object) object.Object {
	if len(args) != 1 {
		return &object.String{Value: "get_int需要1个参数: key"}
	}
	
	key, ok := args[0].(*object.String)
	if !ok {
		return &object.String{Value: "参数类型错误，应为 string"}
	}
	
	value, err := config.GetInt(key.Value)
	if err != nil {
		return &object.String{Value: fmt.Sprintf("获取配置项失败: %v", err)}
	}
	
	return &object.Integer{Value: int64(value)}
}

// evalConfigSetInt 设置整型配置项
func evalConfigSetInt(args []object.Object) object.Object {
	if len(args) != 2 {
		return &object.String{Value: "set_int需要2个参数: key, value"}
	}
	
	key, ok1 := args[0].(*object.String)
	value, ok2 := args[1].(*object.Integer)
	if !ok1 || !ok2 {
		return &object.String{Value: "参数类型错误，应为 (string, integer)"}
	}
	
	config.SetInt(key.Value, int(value.Value))
	return &object.String{Value: "整型配置项设置成功"}
}

// evalConfigGetFloat 获取浮点型配置项
func evalConfigGetFloat(args []object.Object) object.Object {
	if len(args) != 1 {
		return &object.String{Value: "get_float需要1个参数: key"}
	}
	
	key, ok := args[0].(*object.String)
	if !ok {
		return &object.String{Value: "参数类型错误，应为 string"}
	}
	
	value, err := config.GetFloat(key.Value)
	if err != nil {
		return &object.String{Value: fmt.Sprintf("获取配置项失败: %v", err)}
	}
	
	return &object.Float{Value: value}
}

// evalConfigSetFloat 设置浮点型配置项
func evalConfigSetFloat(args []object.Object) object.Object {
	if len(args) != 2 {
		return &object.String{Value: "set_float需要2个参数: key, value"}
	}
	
	key, ok1 := args[0].(*object.String)
	value, ok2 := args[1].(*object.Float)
	if !ok1 || !ok2 {
		return &object.String{Value: "参数类型错误，应为 (string, float)"}
	}
	
	config.SetFloat(key.Value, value.Value)
	return &object.String{Value: "浮点型配置项设置成功"}
}

// evalConfigGetBool 获取布尔型配置项
func evalConfigGetBool(args []object.Object) object.Object {
	if len(args) != 1 {
		return &object.String{Value: "get_bool需要1个参数: key"}
	}
	
	key, ok := args[0].(*object.String)
	if !ok {
		return &object.String{Value: "参数类型错误，应为 string"}
	}
	
	value, err := config.GetBool(key.Value)
	if err != nil {
		return &object.String{Value: fmt.Sprintf("获取配置项失败: %v", err)}
	}
	
	return &object.Boolean{Value: value}
}

// evalConfigSetBool 设置布尔型配置项
func evalConfigSetBool(args []object.Object) object.Object {
	if len(args) != 2 {
		return &object.String{Value: "set_bool需要2个参数: key, value"}
	}
	
	key, ok1 := args[0].(*object.String)
	value, ok2 := args[1].(*object.Boolean)
	if !ok1 || !ok2 {
		return &object.String{Value: "参数类型错误，应为 (string, boolean)"}
	}
	
	config.SetBool(key.Value, value.Value)
	return &object.String{Value: "布尔型配置项设置成功"}
}

// evalConfigHas 检查配置项是否存在
func evalConfigHas(args []object.Object) object.Object {
	if len(args) != 1 {
		return &object.String{Value: "has需要1个参数: key"}
	}
	
	key, ok := args[0].(*object.String)
	if !ok {
		return &object.String{Value: "参数类型错误，应为 string"}
	}
	
	exists := config.Has(key.Value)
	return &object.Boolean{Value: exists}
}

// evalConfigDelete 删除配置项
func evalConfigDelete(args []object.Object) object.Object {
	if len(args) != 1 {
		return &object.String{Value: "delete需要1个参数: key"}
	}
	
	key, ok := args[0].(*object.String)
	if !ok {
		return &object.String{Value: "参数类型错误，应为 string"}
	}
	
	config.Delete(key.Value)
	return &object.String{Value: "配置项删除成功"}
}

// evalConfigGetAll 获取所有配置项
func evalConfigGetAll(args []object.Object) object.Object {
	if len(args) != 0 {
		return &object.String{Value: "get_all不需要参数"}
	}
	
	allConfig := config.GetAll()
	
	// 将map转换为字符串表示
	result := "{"
	first := true
	for key, value := range allConfig {
		if !first {
			result += ", "
		}
		result += fmt.Sprintf("\"%s\": \"%s\"", key, value)
		first = false
	}
	result += "}"
	
	return &object.String{Value: result}
}

// evalConfigSaveToFile 保存配置到文件
func evalConfigSaveToFile(args []object.Object) object.Object {
	if len(args) != 1 {
		return &object.String{Value: "save_to_file需要1个参数: filename"}
	}
	
	filename, ok := args[0].(*object.String)
	if !ok {
		return &object.String{Value: "参数类型错误，应为 string"}
	}
	
	err := config.SaveToFile(filename.Value)
	if err != nil {
		return &object.String{Value: fmt.Sprintf("保存配置文件失败: %v", err)}
	}
	
	return &object.String{Value: "配置文件保存成功"}
}

// evalConfigClear 清空所有配置
func evalConfigClear(args []object.Object) object.Object {
	if len(args) != 0 {
		return &object.String{Value: "clear不需要参数"}
	}
	
	config.Clear()
	return &object.String{Value: "所有配置项已清空"}
}
