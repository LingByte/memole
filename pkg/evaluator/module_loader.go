package evaluator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/LingByte/memole/pkg/lexer"
	"github.com/LingByte/memole/pkg/object"
	"github.com/LingByte/memole/pkg/parser"
)

type moduleLoader struct {
	cache      map[string]*object.ModuleObject
	inProgress map[string]bool
}

func newModuleLoader() *moduleLoader {
	return &moduleLoader{
		cache:      make(map[string]*object.ModuleObject),
		inProgress: make(map[string]bool),
	}
}

func (m *moduleLoader) load(path string, baseDir string) (*object.ModuleObject, error) {
	fullPath := path
	if !strings.HasSuffix(fullPath, ".mml") {
		fullPath += ".mml"
	}
	if !filepath.IsAbs(fullPath) {
		fullPath = filepath.Join(baseDir, fullPath)
	}
	fullPath = filepath.Clean(fullPath)

	if mod, ok := m.cache[fullPath]; ok {
		return mod, nil
	}
	if m.inProgress[fullPath] {
		return nil, fmt.Errorf("检测到循环导入: %s", fullPath)
	}
	m.inProgress[fullPath] = true
	defer delete(m.inProgress, fullPath)

	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("读取模块失败 %s: %w", fullPath, err)
	}

	p, err := parser.New(lexer.New(string(data)))
	if err != nil {
		return nil, fmt.Errorf("解析模块失败 %s: %w", fullPath, err)
	}
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		return nil, fmt.Errorf("模块语法错误 %s: %v", fullPath, p.Errors())
	}

	moduleEnv := parser.NewEnvironment(nil)
	moduleEnv.Set("__module_dir__", &object.String{Value: filepath.Dir(fullPath)})
	Eval(program, moduleEnv)

	exports := map[string]object.Object{}
	for name, val := range moduleEnv.Snapshot() {
		if strings.HasPrefix(name, "__") {
			continue
		}
		exports[name] = val
	}

	moduleName := strings.TrimSuffix(filepath.Base(fullPath), filepath.Ext(fullPath))
	module := &object.ModuleObject{
		Name:    moduleName,
		Exports: exports,
	}
	m.cache[fullPath] = module
	return module, nil
}
