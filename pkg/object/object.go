package object

import (
	"bytes"
	"fmt"
	"strings"
)

// ObjectType 对象类型
type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	FLOAT_OBJ        = "FLOAT"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_OBJ       = "RETURN"
	FUNCTION_OBJ     = "FUNCTION"
	STRING_OBJ       = "STRING"
	PACKAGE_OBJ      = "PACKAGE"
	NETWORK_OBJ      = "NETWORK"
	DB_OBJ           = "DATABASE"
	TABLE_OBJ        = "TABLE"
	ROW_OBJ          = "ROW"
	STRUCT_TYPE_OBJ  = "STRUCT_TYPE"
	STRUCT_INST_OBJ  = "STRUCT_INSTANCE"
	BOUND_METHOD_OBJ = "BOUND_METHOD"
	MODULE_OBJ       = "MODULE"
	NATIVE_FUNC_OBJ  = "NATIVE_FUNCTION"
)

// Object 对象接口
type Object interface {
	Type() ObjectType
	Inspect() string
}

// Integer 整数对象
type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

// Float 浮点数对象
type Float struct {
	Value float64
}

func (f *Float) Type() ObjectType { return FLOAT_OBJ }
func (f *Float) Inspect() string  { return fmt.Sprintf("%f", f.Value) }

// String 字符串对象
type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return fmt.Sprintf("%q", s.Value) }

// Boolean 布尔对象
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

// Null 空对象
type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

// ReturnValue 返回值对象
type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

// Function 函数对象
type Function struct {
	Parameters []*TypedParameter
	Body       *BlockStatement
	Env        interface{} // 使用interface{}来避免循环依赖
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer
	params := make([]string, len(f.Parameters))
	for i, p := range f.Parameters {
		params[i] = p.String()
	}
	out.WriteString("fn(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")
	return out.String()
}

// TypedParameter 带类型的参数
type TypedParameter struct {
	Type string
	Name *Identifier
}

func (tp *TypedParameter) String() string {
	return tp.Type + " " + tp.Name.String()
}

// Identifier 标识符
type Identifier struct {
	Value string
}

func (i *Identifier) String() string { return i.Value }

// BlockStatement 代码块语句
type BlockStatement struct {
	Statements []Statement
}

func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// Statement 语句接口
type Statement interface {
	String() string
}

// PackageObject 包对象
type PackageObject struct {
	Package *Package
}

func (po *PackageObject) Type() ObjectType { return PACKAGE_OBJ }
func (po *PackageObject) Inspect() string  { return fmt.Sprintf("package(%s)", po.Package.Name) }

// Package 包结构
type Package struct {
	Name        string
	Path        string
	Statements  []Statement
	Environment interface{} // 使用interface{}来避免循环依赖
	Imports     []*ImportStatement
}

// ImportStatement 导入语句
type ImportStatement struct {
	Path  *Identifier
	Alias *Identifier
}

func (is *ImportStatement) String() string {
	var out bytes.Buffer
	out.WriteString("import ")
	out.WriteString(is.Path.String())
	if is.Alias != nil {
		out.WriteString(" as ")
		out.WriteString(is.Alias.String())
	}
	out.WriteString(";")
	return out.String()
}

// DatabaseObject 数据库对象
type DatabaseObject struct {
	Backend  string // "mysql", "file", "memory"（原来的 Type 改名）
	Config   map[string]interface{}
	Instance interface{} // 实际的数据库实例
}

func (db *DatabaseObject) Type() ObjectType { return DB_OBJ }
func (db *DatabaseObject) Inspect() string {
	return fmt.Sprintf("database(%s, %v)", db.Backend, db.Config)
}

// TableObject 表对象
type TableObject struct {
	Name   string
	Schema map[string]string // 字段名 -> 类型
	DB     *DatabaseObject
}

func (t *TableObject) Type() ObjectType { return TABLE_OBJ }
func (t *TableObject) Inspect() string {
	return fmt.Sprintf("table(%s, %v)", t.Name, t.Schema)
}

// RowObject 行对象
type RowObject struct {
	Data map[string]interface{}
}

func (r *RowObject) Type() ObjectType { return ROW_OBJ }
func (r *RowObject) Inspect() string {
	return fmt.Sprintf("row(%v)", r.Data)
}

// StructType 结构体类型对象
type StructType struct {
	Name   string
	Fields []string          // 声明顺序
	Types  map[string]string // 字段名 -> 类型名
}

func (st *StructType) Type() ObjectType { return STRUCT_TYPE_OBJ }
func (st *StructType) Inspect() string  { return fmt.Sprintf("struct %s %v", st.Name, st.Fields) }

// StructInstance 结构体实例
type StructInstance struct {
	TypeName string
	Fields   map[string]Object
}

func (si *StructInstance) Type() ObjectType { return STRUCT_INST_OBJ }
func (si *StructInstance) Inspect() string  { return fmt.Sprintf("%s%v", si.TypeName, si.Fields) }

// BoundMethod 绑定方法：保存类型名、方法名与 self
type BoundMethod struct {
	TypeName string
	Method   string
	Self     *StructInstance
}

func (bm *BoundMethod) Type() ObjectType { return BOUND_METHOD_OBJ }
func (bm *BoundMethod) Inspect() string  { return fmt.Sprintf("%s.%s", bm.TypeName, bm.Method) }

// ModuleObject 模块命名空间对象
type ModuleObject struct {
	Name    string
	Exports map[string]Object
}

func (m *ModuleObject) Type() ObjectType { return MODULE_OBJ }
func (m *ModuleObject) Inspect() string {
	return fmt.Sprintf("module(%s)", m.Name)
}

// NativeFunction Go 原生函数对象
type NativeFunction struct {
	Name string
	Fn   func(args []Object) Object
}

func (nf *NativeFunction) Type() ObjectType { return NATIVE_FUNC_OBJ }
func (nf *NativeFunction) Inspect() string  { return fmt.Sprintf("native_fn(%s)", nf.Name) }
