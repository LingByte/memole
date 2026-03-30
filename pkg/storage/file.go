package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// FileStorage 文件存储实现
type FileStorage struct {
	basePath string
	tables   map[string]*FileTable
	mutex    sync.RWMutex
}

// FileTable 文件表结构
type FileTable struct {
	Name   string                   `json:"name"`
	Schema map[string]string        `json:"schema"`
	Data   []map[string]interface{} `json:"data"`
	NextID int64                    `json:"next_id"`
}

func (f *FileStorage) Connect(config map[string]interface{}) error {
	basePath, ok := config["path"].(string)
	if !ok || basePath == "" {
		basePath = "./data"
	}
	f.basePath = basePath
	f.tables = make(map[string]*FileTable)

	if err := os.MkdirAll(basePath, 0755); err != nil {
		return fmt.Errorf("创建数据目录失败: %v", err)
	}
	return f.loadTables()
}

func (f *FileStorage) Disconnect() error {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	for _, t := range f.tables {
		if err := f.saveTable(t); err != nil {
			return err
		}
	}
	return nil
}

func (f *FileStorage) IsConnected() bool { return f.basePath != "" }

func (f *FileStorage) CreateTable(name string, schema map[string]string) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if f.tables[name] != nil {
		return fmt.Errorf("表 %s 已存在", name)
	}
	t := &FileTable{
		Name:   name,
		Schema: schema,
		Data:   make([]map[string]interface{}, 0),
		NextID: 1,
	}
	f.tables[name] = t
	return f.saveTable(t)
}

func (f *FileStorage) DropTable(name string) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if f.tables[name] == nil {
		return fmt.Errorf("表 %s 不存在", name)
	}
	delete(f.tables, name)
	return os.Remove(filepath.Join(f.basePath, name+".json"))
}

func (f *FileStorage) TableExists(name string) bool {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	return f.tables[name] != nil
}

func (f *FileStorage) GetTableSchema(name string) (map[string]string, error) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	t := f.tables[name]
	if t == nil {
		return nil, fmt.Errorf("表 %s 不存在", name)
	}
	return t.Schema, nil
}

func (f *FileStorage) Insert(table string, data map[string]interface{}) (int64, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	t := f.tables[table]
	if t == nil {
		return 0, fmt.Errorf("表 %s 不存在", table)
	}
	data["id"] = t.NextID
	t.NextID++
	t.Data = append(t.Data, data)
	return data["id"].(int64), f.saveTable(t)
}

func (f *FileStorage) Update(table string, where map[string]interface{}, data map[string]interface{}) (int64, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	t := f.tables[table]
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
	if updated > 0 {
		return updated, f.saveTable(t)
	}
	return 0, nil
}

func (f *FileStorage) Delete(table string, where map[string]interface{}) (int64, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	t := f.tables[table]
	if t == nil {
		return 0, fmt.Errorf("表 %s 不存在", table)
	}
	var deleted int64
	newData := make([]map[string]interface{}, 0, len(t.Data))
	for _, row := range t.Data {
		if matchesWhere(row, where) {
			deleted++
		} else {
			newData = append(newData, row)
		}
	}
	t.Data = newData
	if deleted > 0 {
		return deleted, f.saveTable(t)
	}
	return 0, nil
}

func (f *FileStorage) Query(table string, where map[string]interface{}, limit int) ([]map[string]interface{}, error) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	t := f.tables[table]
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

func (f *FileStorage) QueryRaw(sql string, args ...interface{}) ([]map[string]interface{}, error) {
	return nil, fmt.Errorf("文件存储不支持原始SQL查询")
}
func (f *FileStorage) ExecuteRaw(sql string, args ...interface{}) (int64, error) {
	return 0, fmt.Errorf("文件存储不支持原始SQL执行")
}

func (f *FileStorage) BeginTransaction() (TransactionInterface, error) {
	return &FileTransaction{storage: f}, nil
}

func (f *FileStorage) GetType() StorageType { return StorageTypeFile }

func (f *FileStorage) loadTables() error {
	entries, err := os.ReadDir(f.basePath)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".json" {
			continue
		}
		name := strings.TrimSuffix(e.Name(), ".json")
		data, err := os.ReadFile(filepath.Join(f.basePath, e.Name()))
		if err != nil {
			continue
		}
		var t FileTable
		if err := json.Unmarshal(data, &t); err != nil {
			continue
		}
		if t.NextID == 0 {
			// 兼容旧数据
			t.NextID = int64(len(t.Data) + 1)
		}
		f.tables[name] = &t
	}
	return nil
}

func (f *FileStorage) saveTable(t *FileTable) error {
	fp := filepath.Join(f.basePath, t.Name+".json")
	data, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(fp, data, 0644)
}

// FileTransaction：文件事务是“伪事务”，即时写入
type FileTransaction struct {
	storage *FileStorage
}

func (ft *FileTransaction) Commit() error   { return nil }
func (ft *FileTransaction) Rollback() error { return fmt.Errorf("文件存储不支持事务回滚") }
func (ft *FileTransaction) Insert(table string, data map[string]interface{}) (int64, error) {
	return ft.storage.Insert(table, data)
}
func (ft *FileTransaction) Update(table string, where map[string]interface{}, data map[string]interface{}) (int64, error) {
	return ft.storage.Update(table, where, data)
}
func (ft *FileTransaction) Delete(table string, where map[string]interface{}) (int64, error) {
	return ft.storage.Delete(table, where)
}
func (ft *FileTransaction) Query(table string, where map[string]interface{}, limit int) ([]map[string]interface{}, error) {
	return ft.storage.Query(table, where, limit)
}
