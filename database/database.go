package database

import (
	"encoding/json"
	"fmt"
	"memmole/pkg/object"
	"memmole/pkg/storage"
)

// DBPackage 负责管理多个连接
type DBPackage struct {
	factory storage.StorageFactory
	conns   map[int64]storage.StorageInterface
	current int64
	nextID  int64
}

func NewDBPackage() *DBPackage {
	return &DBPackage{
		factory: storage.StorageFactory{},
		conns:   make(map[int64]storage.StorageInterface),
		nextID:  1,
	}
}

// parseJSON 辅助：把 JSON 字符串转 map
func parseJSON(s string) (map[string]interface{}, error) {
	var m map[string]interface{}
	if s == "" {
		return map[string]interface{}{}, nil
	}
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		return nil, fmt.Errorf("JSON解析失败: %v", err)
	}
	return m, nil
}

// Connect 连接并返回连接ID，且设为当前连接
func (p *DBPackage) Connect(backend string, cfgJSON string) *object.Integer {
	cfg, err := parseJSON(cfgJSON)
	if err != nil {
		return &object.Integer{Value: -1}
	}

	var st storage.StorageInterface
	switch backend {
	case "memory":
		st = p.factory.NewStorage(storage.StorageTypeMemory)
	case "file":
		st = p.factory.NewStorage(storage.StorageTypeFile)
	case "mysql":
		st = p.factory.NewStorage(storage.StorageTypeMySQL)
	default:
		return &object.Integer{Value: -1}
	}

	if st == nil {
		return &object.Integer{Value: -1}
	}
	if err := st.Connect(cfg); err != nil {
		return &object.Integer{Value: -1}
	}

	id := p.nextID
	p.nextID++
	p.conns[id] = st
	p.current = id
	return &object.Integer{Value: id}
}

func (p *DBPackage) Use(id int64) *object.String {
	if _, ok := p.conns[id]; !ok {
		return &object.String{Value: "连接不存在"}
	}
	p.current = id
	return &object.String{Value: "ok"}
}

func (p *DBPackage) cur() (storage.StorageInterface, error) {
	st, ok := p.conns[p.current]
	if !ok || st == nil {
		return nil, fmt.Errorf("当前无有效数据库连接")
	}
	return st, nil
}

func (p *DBPackage) CreateTable(name string, schemaJSON string) *object.String {
	st, err := p.cur()
	if err != nil {
		return &object.String{Value: err.Error()}
	}
	schema, err := parseJSON(schemaJSON)
	if err != nil {
		return &object.String{Value: err.Error()}
	}

	// 将 interface{} 的值转成 string（字段类型）
	m := make(map[string]string, len(schema))
	for k, v := range schema {
		m[k] = fmt.Sprintf("%v", v)
	}
	if err := st.CreateTable(name, m); err != nil {
		return &object.String{Value: err.Error()}
	}
	return &object.String{Value: "ok"}
}

func (p *DBPackage) DropTable(name string) *object.String {
	st, err := p.cur()
	if err != nil {
		return &object.String{Value: err.Error()}
	}
	if err := st.DropTable(name); err != nil {
		return &object.String{Value: err.Error()}
	}
	return &object.String{Value: "ok"}
}

func (p *DBPackage) Insert(table string, rowJSON string) *object.Integer {
	st, err := p.cur()
	if err != nil {
		return &object.Integer{Value: -1}
	}
	row, err := parseJSON(rowJSON)
	if err != nil {
		return &object.Integer{Value: -1}
	}
	id, err := st.Insert(table, row)
	if err != nil {
		return &object.Integer{Value: -1}
	}
	return &object.Integer{Value: id}
}

func (p *DBPackage) Update(table string, whereJSON, dataJSON string) *object.Integer {
	st, err := p.cur()
	if err != nil {
		return &object.Integer{Value: -1}
	}
	where, err := parseJSON(whereJSON)
	if err != nil {
		return &object.Integer{Value: -1}
	}
	data, err := parseJSON(dataJSON)
	if err != nil {
		return &object.Integer{Value: -1}
	}
	n, err := st.Update(table, where, data)
	if err != nil {
		return &object.Integer{Value: -1}
	}
	return &object.Integer{Value: n}
}

func (p *DBPackage) Delete(table string, whereJSON string) *object.Integer {
	st, err := p.cur()
	if err != nil {
		return &object.Integer{Value: -1}
	}
	where, err := parseJSON(whereJSON)
	if err != nil {
		return &object.Integer{Value: -1}
	}
	n, err := st.Delete(table, where)
	if err != nil {
		return &object.Integer{Value: -1}
	}
	return &object.Integer{Value: n}
}

func (p *DBPackage) Query(table string, whereJSON string, limit int64) *object.String {
	st, err := p.cur()
	if err != nil {
		return &object.String{Value: err.Error()}
	}
	where, err := parseJSON(whereJSON)
	if err != nil {
		return &object.String{Value: err.Error()}
	}
	rows, err := st.Query(table, where, int(limit))
	if err != nil {
		return &object.String{Value: err.Error()}
	}

	// 返回 JSON 字符串，方便在 MML 里打印/透传
	b, _ := json.Marshal(rows)
	return &object.String{Value: string(b)}
}

func (p *DBPackage) QueryRaw(sql string) *object.String {
	st, err := p.cur()
	if err != nil {
		return &object.String{Value: err.Error()}
	}
	rows, err := st.QueryRaw(sql)
	if err != nil {
		return &object.String{Value: err.Error()}
	}

	// 返回 JSON 字符串，方便在 MML 里打印/透传
	b, _ := json.Marshal(rows)
	return &object.String{Value: string(b)}
}

func (p *DBPackage) ExecuteRaw(sql string) *object.Integer {
	st, err := p.cur()
	if err != nil {
		return &object.Integer{Value: -1}
	}
	n, err := st.ExecuteRaw(sql)
	if err != nil {
		return &object.Integer{Value: -1}
	}
	return &object.Integer{Value: n}
}

// DBObject：放进解释器环境中用的对象包装（仿照 network.NetworkObject）
type DBObject struct{ Package *DBPackage }

func (dbo *DBObject) Type() object.ObjectType { return "DATABASE" }
func (dbo *DBObject) Inspect() string         { return "db_package" }
