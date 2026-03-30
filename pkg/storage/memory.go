package storage

import (
	"fmt"
	"sync"
)

// MemoryStorage 内存存储实现
type MemoryStorage struct {
	tables map[string]*MemoryTable
	mutex  sync.RWMutex
}

// MemoryTable 内存表结构
type MemoryTable struct {
	Name   string
	Schema map[string]string
	Data   []map[string]interface{}
	NextID int64
}

// Connect 连接内存存储
func (m *MemoryStorage) Connect(config map[string]interface{}) error {
	m.tables = make(map[string]*MemoryTable)
	return nil
}

// Disconnect 断开连接
func (m *MemoryStorage) Disconnect() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.tables = make(map[string]*MemoryTable)
	return nil
}

// IsConnected 检查是否已连接
func (m *MemoryStorage) IsConnected() bool {
	return true
}

// CreateTable 创建表
func (m *MemoryStorage) CreateTable(name string, schema map[string]string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.tables[name] != nil {
		return fmt.Errorf("表 %s 已存在", name)
	}

	m.tables[name] = &MemoryTable{
		Name:   name,
		Schema: schema,
		Data:   make([]map[string]interface{}, 0),
		NextID: 1,
	}

	return nil
}

// DropTable 删除表
func (m *MemoryStorage) DropTable(name string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.tables[name] == nil {
		return fmt.Errorf("表 %s 不存在", name)
	}
	delete(m.tables, name)
	return nil
}

// TableExists 检查表是否存在
func (m *MemoryStorage) TableExists(name string) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.tables[name] != nil
}

// GetTableSchema 获取表结构
func (m *MemoryStorage) GetTableSchema(name string) (map[string]string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	table := m.tables[name]
	if table == nil {
		return nil, fmt.Errorf("表 %s 不存在", name)
	}
	return table.Schema, nil
}

// Insert 插入数据
func (m *MemoryStorage) Insert(table string, data map[string]interface{}) (int64, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	t := m.tables[table]
	if t == nil {
		return 0, fmt.Errorf("表 %s 不存在", table)
	}

	data["id"] = t.NextID
	t.NextID++
	t.Data = append(t.Data, data)

	return data["id"].(int64), nil
}

// Update 更新数据
func (m *MemoryStorage) Update(table string, where map[string]interface{}, data map[string]interface{}) (int64, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	t := m.tables[table]
	if t == nil {
		return 0, fmt.Errorf("表 %s 不存在", table)
	}

	var updated int64
	for i, row := range t.Data {
		if matchesWhere(row, where) {
			id := row["id"]
			for k, v := range data {
				row[k] = v
			}
			row["id"] = id
			t.Data[i] = row
			updated++
		}
	}
	return updated, nil
}

// Delete 删除数据
func (m *MemoryStorage) Delete(table string, where map[string]interface{}) (int64, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	t := m.tables[table]
	if t == nil {
		return 0, fmt.Errorf("表 %s 不存在", table)
	}

	var deleted int64
	newData := make([]map[string]interface{}, 0, len(t.Data))
	for _, row := range t.Data {
		if matchesWhere(row, where) {
			deleted++
			continue
		}
		newData = append(newData, row)
	}
	t.Data = newData
	return deleted, nil
}

// Query 查询数据
func (m *MemoryStorage) Query(table string, where map[string]interface{}, limit int) ([]map[string]interface{}, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	t := m.tables[table]
	if t == nil {
		return nil, fmt.Errorf("表 %s 不存在", table)
	}

	var out []map[string]interface{}
	for _, row := range t.Data {
		if limit > 0 && len(out) >= limit {
			break
		}
		if matchesWhere(row, where) {
			cp := make(map[string]interface{}, len(row))
			for k, v := range row {
				cp[k] = v
			}
			out = append(out, cp)
		}
	}
	return out, nil
}

// QueryRaw/ExecuteRaw：内存存储不支持 SQL
func (m *MemoryStorage) QueryRaw(sql string, args ...interface{}) ([]map[string]interface{}, error) {
	return nil, fmt.Errorf("内存存储不支持原始SQL查询")
}
func (m *MemoryStorage) ExecuteRaw(sql string, args ...interface{}) (int64, error) {
	return 0, fmt.Errorf("内存存储不支持原始SQL执行")
}

// 事务（模拟）
func (m *MemoryStorage) BeginTransaction() (TransactionInterface, error) {
	return &MemoryTransaction{storage: m}, nil
}

func (m *MemoryStorage) GetType() StorageType {
	return StorageTypeMemory
}

func matchesWhere(row map[string]interface{}, where map[string]interface{}) bool {
	for k, v := range where {
		if row[k] != v {
			return false
		}
	}
	return true
}

// MemoryTransaction
type MemoryTransaction struct {
	storage *MemoryStorage
}

func (mt *MemoryTransaction) Commit() error   { return nil }
func (mt *MemoryTransaction) Rollback() error { return fmt.Errorf("内存存储不支持事务回滚") }

func (mt *MemoryTransaction) Insert(table string, data map[string]interface{}) (int64, error) {
	return mt.storage.Insert(table, data)
}
func (mt *MemoryTransaction) Update(table string, where map[string]interface{}, data map[string]interface{}) (int64, error) {
	return mt.storage.Update(table, where, data)
}
func (mt *MemoryTransaction) Delete(table string, where map[string]interface{}) (int64, error) {
	return mt.storage.Delete(table, where)
}
func (mt *MemoryTransaction) Query(table string, where map[string]interface{}, limit int) ([]map[string]interface{}, error) {
	return mt.storage.Query(table, where, limit)
}
