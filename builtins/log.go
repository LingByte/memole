package builtins

import (
	"memmole/pkg/ast"
	"memmole/pkg/parser"
	"memmole/pkg/object"
	"memmole/pkg/logger"
)

// EvalLogCall 处理日志包调用
func EvalLogCall(exp *ast.MemberAccessExpression, args []object.Object, env *parser.Environment) object.Object {
	switch exp.Member.Value {
	case "debug":
		if len(args) != 1 {
			return &object.String{Value: "log.debug需要1个参数: message"}
		}
		message, ok := args[0].(*object.String)
		if !ok {
			return &object.String{Value: "log.debug参数必须是字符串"}
		}
		logger.Debug(message.Value)
		return &object.String{Value: "日志已输出"}

	case "info":
		if len(args) != 1 {
			return &object.String{Value: "log.info需要1个参数: message"}
		}
		message, ok := args[0].(*object.String)
		if !ok {
			return &object.String{Value: "log.info参数必须是字符串"}
		}
		logger.Info(message.Value)
		return &object.String{Value: "日志已输出"}

	case "warn":
		if len(args) != 1 {
			return &object.String{Value: "log.warn需要1个参数: message"}
		}
		message, ok := args[0].(*object.String)
		if !ok {
			return &object.String{Value: "log.warn参数必须是字符串"}
		}
		logger.Warn(message.Value)
		return &object.String{Value: "日志已输出"}

	case "error":
		if len(args) != 1 {
			return &object.String{Value: "log.error需要1个参数: message"}
		}
		message, ok := args[0].(*object.String)
		if !ok {
			return &object.String{Value: "log.error参数必须是字符串"}
		}
		logger.Error(message.Value)
		return &object.String{Value: "日志已输出"}

	case "fatal":
		if len(args) != 1 {
			return &object.String{Value: "log.fatal需要1个参数: message"}
		}
		message, ok := args[0].(*object.String)
		if !ok {
			return &object.String{Value: "log.fatal参数必须是字符串"}
		}
		logger.Fatal(message.Value)
		return &object.String{Value: "程序已退出"}

	case "set_level":
		if len(args) != 1 {
			return &object.String{Value: "log.set_level需要1个参数: level (0=DEBUG, 1=INFO, 2=WARN, 3=ERROR, 4=FATAL)"}
		}
		level, ok := args[0].(*object.Integer)
		if !ok {
			return &object.String{Value: "log.set_level参数必须是整数"}
		}
		logger.SetLevel(logger.LogLevel(level.Value))
		return &object.String{Value: "日志级别已设置"}

	case "get_level":
		if len(args) != 0 {
			return &object.String{Value: "log.get_level不需要参数"}
		}
		level := logger.GetLevel()
		return &object.Integer{Value: int64(level)}

	default:
		return &object.String{Value: "未知的日志方法: " + exp.Member.Value}
	}
}
