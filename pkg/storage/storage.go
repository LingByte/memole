package storage

// StorageType 存储类型
type StorageType string

const (
	StorageTypeMySQL  StorageType = "mysql"
	StorageTypeFile   StorageType = "file"
	StorageTypeMemory StorageType = "memory"
)

// StorageInterface 存储接口
type StorageInterface interface {
	// 连接管理
	Connect(config map[string]interface{}) error
	Disconnect() error
	IsConnected() bool

	// 表操作
	CreateTable(name string, schema map[string]string) error
	DropTable(name string) error
	TableExists(name string) bool
	GetTableSchema(name string) (map[string]string, error)

	// 数据操作
	Insert(table string, data map[string]interface{}) (int64, error)
	Update(table string, where map[string]interface{}, data map[string]interface{}) (int64, error)
	Delete(table string, where map[string]interface{}) (int64, error)
	Query(table string, where map[string]interface{}, limit int) ([]map[string]interface{}, error)
	QueryRaw(sql string, args ...interface{}) ([]map[string]interface{}, error)
	ExecuteRaw(sql string, args ...interface{}) (int64, error)

	// 事务支持
	BeginTransaction() (TransactionInterface, error)

	// 获取存储类型
	GetType() StorageType
}

// TransactionInterface 事务接口
type TransactionInterface interface {
	Commit() error
	Rollback() error
	Insert(table string, data map[string]interface{}) (int64, error)
	Update(table string, where map[string]interface{}, data map[string]interface{}) (int64, error)
	Delete(table string, where map[string]interface{}) (int64, error)
	Query(table string, where map[string]interface{}, limit int) ([]map[string]interface{}, error)
}

// StorageFactory 存储工厂
type StorageFactory struct{}

// NewStorage 根据类型创建存储实例
func (sf *StorageFactory) NewStorage(t StorageType) StorageInterface {
	switch t {
	case StorageTypeMySQL:
		return &MySQLStorage{}
	case StorageTypeFile:
		return &FileStorage{}
	case StorageTypeMemory:
		return &MemoryStorage{}
	default:
		return nil
	}
}

/*
以下为未来 ORM 的接口预留（保留占位即可，不引用外部包，以免引入未使用依赖）
*/
type ORMInterface interface {
	DefineModel(name string, fields map[string]string) error
	Query() QueryBuilder
	Where(field string, operator string, value interface{}) QueryBuilder
	OrderBy(field string, direction string) QueryBuilder
	Limit(limit int) QueryBuilder
	Offset(offset int) QueryBuilder

	Create(model interface{}) error
	Update(model interface{}) error
	Delete(model interface{}) error
	Find(id interface{}) (interface{}, error)
	FindAll() ([]interface{}, error)
}

type QueryBuilder interface {
	Where(field string, operator string, value interface{}) QueryBuilder
	OrderBy(field string, direction string) QueryBuilder
	Limit(limit int) QueryBuilder
	Offset(offset int) QueryBuilder
	Get() ([]map[string]interface{}, error)
	First() (map[string]interface{}, error)
	Count() (int64, error)
	Update(data map[string]interface{}) (int64, error)
	Delete() (int64, error)
}
