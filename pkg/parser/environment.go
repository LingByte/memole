package parser

import (
	"memmole/pkg/object"
)

// Environment 环境管理结构
type Environment struct {
	store map[string]object.Object
	outer *Environment
}

// NewEnvironment 创建新环境
func NewEnvironment(outer *Environment) *Environment {
	return &Environment{
		store: make(map[string]object.Object),
		outer: outer,
	}
}

// Get 从环境中获取变量
func (e *Environment) Get(name string) (object.Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		return e.outer.Get(name)
	}
	return obj, ok
}

// Set 在环境中设置变量
func (e *Environment) Set(name string, val object.Object) object.Object {
	e.store[name] = val
	return val
}
