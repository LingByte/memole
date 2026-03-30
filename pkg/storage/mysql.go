package storage

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// MySQLStorage MySQL存储实现
type MySQLStorage struct {
	db     *sql.DB
	config map[string]interface{}
}

func (m *MySQLStorage) Connect(config map[string]interface{}) error {
	m.config = config

	host, _ := config["host"].(string)
	if host == "" {
		host = "localhost"
	}
	port, _ := config["port"].(string)
	if port == "" {
		port = "3306"
	}
	user, _ := config["user"].(string)
	if user == "" {
		user = "root"
	}
	password, _ := config["password"].(string)
	database, _ := config["database"].(string)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, database)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("连接MySQL失败: %v", err)
	}
	if err := db.Ping(); err != nil {
		return fmt.Errorf("MySQL连接测试失败: %v", err)
	}
	m.db = db
	return nil
}

func (m *MySQLStorage) Disconnect() error {
	if m.db != nil {
		return m.db.Close()
	}
	return nil
}

func (m *MySQLStorage) IsConnected() bool { return m.db != nil }

func (m *MySQLStorage) CreateTable(name string, schema map[string]string) error {
	if !m.IsConnected() {
		return fmt.Errorf("数据库未连接")
	}
	var cols []string
	for field, typ := range schema {
		cols = append(cols, fmt.Sprintf("`%s` %s", field, typ))
	}
	cols = append(cols, "`id` INT AUTO_INCREMENT PRIMARY KEY")
	q := fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s` (%s)", name, strings.Join(cols, ", "))
	_, err := m.db.Exec(q)
	return err
}

func (m *MySQLStorage) DropTable(name string) error {
	if !m.IsConnected() {
		return fmt.Errorf("数据库未连接")
	}
	_, err := m.db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS `%s`", name))
	return err
}

func (m *MySQLStorage) TableExists(name string) bool {
	if !m.IsConnected() {
		return false
	}
	q := "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?"
	var c int
	err := m.db.QueryRow(q, name).Scan(&c)
	return err == nil && c > 0
}

func (m *MySQLStorage) GetTableSchema(name string) (map[string]string, error) {
	if !m.IsConnected() {
		return nil, fmt.Errorf("数据库未连接")
	}
	q := "SELECT COLUMN_NAME, DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?"
	rows, err := m.db.Query(q, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	schema := make(map[string]string)
	for rows.Next() {
		var col, typ string
		if err := rows.Scan(&col, &typ); err != nil {
			return nil, err
		}
		if col != "id" {
			schema[col] = typ
		}
	}
	return schema, nil
}

func (m *MySQLStorage) Insert(table string, data map[string]interface{}) (int64, error) {
	if !m.IsConnected() {
		return 0, fmt.Errorf("数据库未连接")
	}
	var fields, placeholders []string
	var values []interface{}
	for k, v := range data {
		fields = append(fields, fmt.Sprintf("`%s`", k))
		placeholders = append(placeholders, "?")
		values = append(values, v)
	}
	q := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s)", table, strings.Join(fields, ", "), strings.Join(placeholders, ", "))
	res, err := m.db.Exec(q, values...)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (m *MySQLStorage) Update(table string, where map[string]interface{}, data map[string]interface{}) (int64, error) {
	if !m.IsConnected() {
		return 0, fmt.Errorf("数据库未连接")
	}
	var setParts []string
	var args []interface{}
	for k, v := range data {
		setParts = append(setParts, fmt.Sprintf("`%s` = ?", k))
		args = append(args, v)
	}
	var whereParts []string
	for k, v := range where {
		whereParts = append(whereParts, fmt.Sprintf("`%s` = ?", k))
		args = append(args, v)
	}
	q := fmt.Sprintf("UPDATE `%s` SET %s", table, strings.Join(setParts, ", "))
	if len(whereParts) > 0 {
		q += " WHERE " + strings.Join(whereParts, " AND ")
	}
	res, err := m.db.Exec(q, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (m *MySQLStorage) Delete(table string, where map[string]interface{}) (int64, error) {
	if !m.IsConnected() {
		return 0, fmt.Errorf("数据库未连接")
	}
	var whereParts []string
	var args []interface{}
	for k, v := range where {
		whereParts = append(whereParts, fmt.Sprintf("`%s` = ?", k))
		args = append(args, v)
	}
	q := fmt.Sprintf("DELETE FROM `%s`", table)
	if len(whereParts) > 0 {
		q += " WHERE " + strings.Join(whereParts, " AND ")
	}
	res, err := m.db.Exec(q, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (m *MySQLStorage) Query(table string, where map[string]interface{}, limit int) ([]map[string]interface{}, error) {
	if !m.IsConnected() {
		return nil, fmt.Errorf("数据库未连接")
	}
	var whereParts []string
	var args []interface{}
	for k, v := range where {
		whereParts = append(whereParts, fmt.Sprintf("`%s` = ?", k))
		args = append(args, v)
	}
	q := fmt.Sprintf("SELECT * FROM `%s`", table)
	if len(whereParts) > 0 {
		q += " WHERE " + strings.Join(whereParts, " AND ")
	}
	if limit > 0 {
		q += fmt.Sprintf(" LIMIT %d", limit)
	}
	rows, err := m.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRows(rows)
}

func (m *MySQLStorage) QueryRaw(sqlStr string, args ...interface{}) ([]map[string]interface{}, error) {
	if !m.IsConnected() {
		return nil, fmt.Errorf("数据库未连接")
	}
	rows, err := m.db.Query(sqlStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRows(rows)
}

func (m *MySQLStorage) ExecuteRaw(sqlStr string, args ...interface{}) (int64, error) {
	if !m.IsConnected() {
		return 0, fmt.Errorf("数据库未连接")
	}
	res, err := m.db.Exec(sqlStr, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (m *MySQLStorage) BeginTransaction() (TransactionInterface, error) {
	if !m.IsConnected() {
		return nil, fmt.Errorf("数据库未连接")
	}
	tx, err := m.db.Begin()
	if err != nil {
		return nil, err
	}
	return &MySQLTransaction{tx: tx}, nil
}

func (m *MySQLStorage) GetType() StorageType { return StorageTypeMySQL }

// 工具：扫描行
func scanRows(rows *sql.Rows) ([]map[string]interface{}, error) {
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	var out []map[string]interface{}
	for rows.Next() {
		vals := make([]interface{}, len(cols))
		ptrs := make([]interface{}, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			return nil, err
		}
		row := make(map[string]interface{}, len(cols))
		for i, c := range cols {
			if vals[i] != nil {
				row[c] = vals[i]
			}
		}
		out = append(out, row)
	}
	return out, nil
}

// MySQLTransaction
type MySQLTransaction struct {
	tx *sql.Tx
}

func (mt *MySQLTransaction) Commit() error   { return mt.tx.Commit() }
func (mt *MySQLTransaction) Rollback() error { return mt.tx.Rollback() }

func (mt *MySQLTransaction) Insert(table string, data map[string]interface{}) (int64, error) {
	var fields, placeholders []string
	var args []interface{}
	for k, v := range data {
		fields = append(fields, fmt.Sprintf("`%s`", k))
		placeholders = append(placeholders, "?")
		args = append(args, v)
	}
	q := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s)", table, strings.Join(fields, ", "), strings.Join(placeholders, ", "))
	res, err := mt.tx.Exec(q, args...)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (mt *MySQLTransaction) Update(table string, where map[string]interface{}, data map[string]interface{}) (int64, error) {
	var setParts []string
	var args []interface{}
	for k, v := range data {
		setParts = append(setParts, fmt.Sprintf("`%s` = ?", k))
		args = append(args, v)
	}
	var whereParts []string
	for k, v := range where {
		whereParts = append(whereParts, fmt.Sprintf("`%s` = ?", k))
		args = append(args, v)
	}
	q := fmt.Sprintf("UPDATE `%s` SET %s", table, strings.Join(setParts, ", "))
	if len(whereParts) > 0 {
		q += " WHERE " + strings.Join(whereParts, " AND ")
	}
	res, err := mt.tx.Exec(q, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (mt *MySQLTransaction) Delete(table string, where map[string]interface{}) (int64, error) {
	var whereParts []string
	var args []interface{}
	for k, v := range where {
		whereParts = append(whereParts, fmt.Sprintf("`%s` = ?", k))
		args = append(args, v)
	}
	q := fmt.Sprintf("DELETE FROM `%s`", table)
	if len(whereParts) > 0 {
		q += " WHERE " + strings.Join(whereParts, " AND ")
	}
	res, err := mt.tx.Exec(q, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (mt *MySQLTransaction) Query(table string, where map[string]interface{}, limit int) ([]map[string]interface{}, error) {
	var whereParts []string
	var args []interface{}
	for k, v := range where {
		whereParts = append(whereParts, fmt.Sprintf("`%s` = ?", k))
		args = append(args, v)
	}
	q := fmt.Sprintf("SELECT * FROM `%s`", table)
	if len(whereParts) > 0 {
		q += " WHERE " + strings.Join(whereParts, " AND ")
	}
	if limit > 0 {
		q += fmt.Sprintf(" LIMIT %d", limit)
	}
	rows, err := mt.tx.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRows(rows)
}
