package parser

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/LingByte/memole/pkg/ast"
	"github.com/LingByte/memole/pkg/lexer"
)

// PackageManager 包管理器
type PackageManager struct {
	packages map[string]*Package
	basePath string
}

// Package 表示一个包
type Package struct {
	Name        string
	Path        string
	Statements  []ast.Statement
	Environment *Environment
	Imports     []*ast.ImportStatement
}

// NewPackageManager 创建新的包管理器
func NewPackageManager(basePath string) *PackageManager {
	return &PackageManager{
		packages: make(map[string]*Package),
		basePath: basePath,
	}
}

// LoadPackage 加载包
func (pm *PackageManager) LoadPackage(packagePath string) (*Package, error) {
	// 检查是否已经加载
	if pkg, exists := pm.packages[packagePath]; exists {
		return pkg, nil
	}

	// 构建完整的文件路径
	fullPath := filepath.Join(pm.basePath, packagePath+".mml")

	// 读取文件内容
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("无法读取包文件 %s: %v", fullPath, err)
	}

	// 创建词法分析器
	l := lexer.New(string(content))

	// 创建语法分析器
	p, err := New(l)
	if err != nil {
		return nil, fmt.Errorf("解析包文件失败: %v", err)
	}

	// 解析程序
	program := p.ParseProgram()

	// 检查语法错误
	if len(p.Errors()) > 0 {
		return nil, fmt.Errorf("包文件语法错误: %v", p.Errors())
	}

	// 创建包对象
	pkg := &Package{
		Name:        packagePath,
		Path:        fullPath,
		Statements:  program.Statements,
		Environment: NewEnvironment(nil),
		Imports:     []*ast.ImportStatement{},
	}

	// 处理包中的导入语句
	for _, stmt := range program.Statements {
		if importStmt, ok := stmt.(*ast.ImportStatement); ok {
			pkg.Imports = append(pkg.Imports, importStmt)
		}
	}

	// 缓存包
	pm.packages[packagePath] = pkg

	return pkg, nil
}

// GetPackage 获取已加载的包
func (pm *PackageManager) GetPackage(name string) (*Package, bool) {
	pkg, exists := pm.packages[name]
	return pkg, exists
}

// ResolveImports 解析包的所有导入
func (pm *PackageManager) ResolveImports(pkg *Package) error {
	for _, importStmt := range pkg.Imports {
		importPath := importStmt.Path.Value

		// 加载导入的包
		importedPkg, err := pm.LoadPackage(importPath)
		if err != nil {
			return fmt.Errorf("导入包 %s 失败: %v", importPath, err)
		}

		// 确定包名（使用别名或包名）
		packageName := importPath
		if importStmt.Alias != nil {
			packageName = importStmt.Alias.Value
		}

		// 将导入的包添加到当前包的环境
		// 暂时跳过包导入，避免类型问题
		_ = packageName
		_ = importedPkg
	}

	return nil
}
