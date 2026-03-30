# MML 数据库使用示例

本文档提供了MML语言中数据库功能的详细使用示例。

## 支持的数据库类型

MML支持三种数据库类型：

1. **内存数据库 (memory)** - 数据存储在内存中，程序结束后数据丢失
2. **文件数据库 (file)** - 数据持久化到JSON文件，程序结束后数据保留
3. **MySQL数据库 (mysql)** - 数据存储在MySQL服务器中，支持复杂SQL查询

## 基本用法

### 1. 连接数据库

```mml
// 连接内存数据库
io.println(db.connect("memory", "{}"));

// 连接文件数据库
io.println(db.connect("file", "{\"path\":\"./data\"}"));

// 连接MySQL数据库
let mysql_config = "{\"host\":\"localhost\", \"port\":\"3306\", \"user\":\"root\", \"password\":\"1234\", \"database\":\"test_db\"}";
io.println(db.connect("mysql", mysql_config));
```

### 2. 创建表

```mml
// 创建用户表
let schema = "{\"name\":\"varchar(100)\", \"age\":\"int\", \"email\":\"varchar(255)\"}";
io.println(db.create_table("users", schema));
```

### 3. 插入数据

```mml
// 插入用户数据
io.println(db.insert("users", "{\"name\":\"张三\", \"age\":25, \"email\":\"zhangsan@example.com\"}"));
```

### 4. 查询数据

```mml
// 查询所有用户
io.println(db.query("users", "{}", 10));

// 查询特定条件的用户
io.println(db.query("users", "{\"age\":25}", 10));
```

### 5. 更新数据

```mml
// 更新用户信息
io.println(db.update("users", "{\"name\":\"张三\"}", "{\"age\":26, \"email\":\"zhangsan_new@example.com\"}"));
```

### 6. 删除数据

```mml
// 删除用户
io.println(db.delete("users", "{\"name\":\"张三\"}"));
```

### 7. 切换数据库连接

```mml
// 连接多个数据库
let memory_conn = db.connect("memory", "{}");
let file_conn = db.connect("file", "{\"path\":\"./data\"}");

// 切换到内存数据库
db.use(memory_conn);

// 切换到文件数据库
db.use(file_conn);
```

## 示例文件

### 1. memory_db_demo.mml
内存数据库完整示例，包括：
- 连接内存数据库
- 创建用户表
- 插入、查询、更新、删除用户数据

### 2. file_db_demo.mml
文件数据库完整示例，包括：
- 连接文件数据库
- 创建文章表
- 插入、查询、更新、删除文章数据

### 3. mysql_db_demo.mml
MySQL数据库完整示例，包括：
- 连接MySQL数据库
- 创建产品表
- 插入、查询、更新、删除产品数据
- 使用原始SQL查询

### 4. multi_db_demo.mml
多数据库连接示例，展示如何：
- 同时连接多个数据库
- 在不同数据库间切换
- 比较不同数据库的特性

## 运行示例

```bash
# 运行内存数据库示例
.\mml.exe .\examples\memory_db_demo.mml

# 运行文件数据库示例
.\mml.exe .\examples\file_db_demo.mml

# 运行MySQL数据库示例（需要先启动MySQL服务）
.\mml.exe .\examples\mysql_db_demo.mml

# 运行多数据库示例
.\mml.exe .\examples\multi_db_demo.mml
```

## 注意事项

1. **MySQL数据库**：运行MySQL示例前，请确保：
   - MySQL服务已启动
   - 数据库 `test_db` 已创建
   - 用户名和密码正确

2. **文件数据库**：数据文件会保存在指定的目录中（如 `./data/`）

3. **内存数据库**：数据仅在程序运行期间有效，程序结束后数据丢失

4. **JSON格式**：所有配置和数据都使用JSON格式，注意转义字符的使用

## 数据库特性对比

| 特性 | 内存数据库 | 文件数据库 | MySQL数据库 |
|------|------------|------------|-------------|
| 数据持久性 | ❌ | ✅ | ✅ |
| 性能 | 最快 | 中等 | 中等 |
| 并发支持 | 有限 | 有限 | 完整 |
| SQL查询 | ❌ | ❌ | ✅ |
| 事务支持 | 模拟 | 模拟 | 完整 |
| 适用场景 | 临时数据、缓存 | 小型应用、原型 | 生产环境 |
