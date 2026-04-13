package builtinmodules

import "github.com/LingByte/memole/pkg/object"

var modules = map[string]*object.ModuleObject{
	"std": newStdModule(),
}

// Get 返回官方内置模块
func Get(name string) (*object.ModuleObject, bool) {
	mod, ok := modules[name]
	return mod, ok
}
