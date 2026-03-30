# MemMole - 自定义虚拟机语言

MemMole是一个用Go语言实现的自定义虚拟机语言，支持基本的编程语言功能，包括变量声明、函数定义、控制流、网络操作等。

## 🚀 特性

- **简洁的语法**：类似C语言的语法设计
- **类型系统**：支持整数、字符串、布尔值等基本类型
- **函数支持**：支持函数定义和调用
- **控制流**：支持if-else条件语句和while循环
- **网络功能**：内置HTTP、TCP、UDP等网络操作，支持HTTP服务器开发
- **包管理**：支持模块化开发和包导入
- **交互式REPL**：提供命令行交互环境

## 📁 项目结构

```
MemMole/
├── ast/                    # 抽象语法树节点定义
│   ├── node.go            # 基础节点接口
│   ├── statements.go      # 语句类型定义
│   └── expressions.go     # 表达式类型定义
├── object/                # 对象系统
│   └── object.go          # 对象接口和类型定义
├── parser/                # 语法解析器
│   ├── parser.go          # 主解析器逻辑
│   ├── statements.go      # 语句解析
│   ├── expressions.go     # 表达式解析
│   ├── environment.go     # 环境管理
│   └── package_manager.go # 包管理
├── evaluator/             # 代码求值器
│   ├── evaluator.go       # 主求值逻辑
│   ├── statements.go      # 语句求值
│   └── expressions.go     # 表达式求值
├── builtins/              # 内置函数
│   └── network.go         # 网络功能
├── repl/                  # 交互式环境
│   └── repl.go            # REPL实现
├── lexer/                 # 词法分析器
├── network/               # 网络包
├── examples/              # 示例代码
└── main.go                # 程序入口
```

## 🏗️ 架构设计

### 设计原则

- **单一职责原则**：每个模块专注于特定功能
- **模块化设计**：功能分离，降低耦合度
- **清晰的依赖关系**：避免循环依赖
- **可扩展性**：易于添加新功能

### 核心组件

1. **词法分析器 (Lexer)**：将源代码转换为Token序列
2. **语法分析器 (Parser)**：构建抽象语法树(AST)
3. **求值器 (Evaluator)**：执行AST并产生结果
4. **对象系统 (Object)**：定义语言中的数据类型
5. **环境管理 (Environment)**：管理变量作用域
6. **内置函数 (Builtins)**：提供网络等系统功能

## 📖 语言语法

### 基本语法

```mml
// 变量声明
let x = 42;
let message = "Hello, World!";

// 函数定义
fn void main() {
    return 0;
}

// 条件语句
if (x > 10) {
    return true;
} else {
    return false;
}

// 循环语句
while (x > 0) {
    x = x - 1;
}

// 函数调用
let result = add(5, 3);
```

### 网络操作

#### 客户端功能

```mml
// HTTP GET请求
let response = network.http_get("https://api.example.com/data");

// HTTP POST请求
let result = network.http_post("https://api.example.com/submit", "{\"key\":\"value\"}");

// TCP连接
let conn = network.tcp_connect("localhost", 8080);
network.tcp_send(conn, "Hello Server");
let data = network.tcp_receive(conn);
network.tcp_close(conn);

// UDP发送
network.udp_send("localhost", 8080, "Hello UDP");

// DNS解析
let ip = network.resolve_dns("example.com");

// Ping测试
let status = network.ping("google.com");
```

#### 服务器功能

```mml
// 创建HTTP服务器
let server_id = network.create_server(8080);

// 添加路由处理器
network.add_route(server_id, "GET", "/", "欢迎访问MemMole服务器！");
network.add_route(server_id, "GET", "/api/users", "{\"users\":[{\"id\":1,\"name\":\"Alice\"}]}");
network.add_route(server_id, "POST", "/api/users", "用户创建成功");

// 启动服务器（非阻塞模式）
let result = network.start_server(server_id);

// 启动服务器（阻塞模式，程序会保持运行）
let result = network.start_server_and_wait(server_id);

// 检查服务器状态
let is_running = network.is_server_running(server_id);
let port = network.get_server_port(server_id);

// 停止服务器
network.stop_server(server_id);
```

### 包管理

```mml
// 包声明
package myapp;

// 导入包
import math;
import utils as util;

// 使用导入的函数
let result = math.add(5, 3);
let formatted = util.format("Hello");
```

### 服务器运行模式

CVM支持两种服务器运行模式：

1. **非阻塞模式** (`network.start_server`)：
   - 服务器在后台运行
   - 程序立即返回，继续执行后续代码
   - 适合需要同时执行其他任务的场景

2. **阻塞模式** (`network.start_server_and_wait`)：
   - 服务器在前台运行
   - 程序会阻塞等待，直到服务器停止
   - 适合长期运行的Web服务

## 🛠️ 安装和使用

### 环境要求

- Go 1.19 或更高版本

### 编译

```bash
git clone <repository-url>
cd MemMole
go build
```

### 运行

#### 交互式模式
```bash
./mml
```

#### 文件模式
```bash
./mml examples/test.mml
```

## 📝 示例代码

### 基础计算

```mml
fn int add(int a, int b) {
    return a + b;
}

fn void main() {
    let x = 10;
    let y = 20;
    let result = add(x, y);
    return result;
}
```
```

### 网络客户端

```mml
fn void main() {
    // 获取网页内容
    let response = network.http_get("https://httpbin.org/get");
    
    // 发送POST请求
    let data = "{\"name\":\"MemMole\",\"version\":\"1.0\"}";
    let result = network.http_post("https://httpbin.org/post", data);
    
    return 0;
}
```
```

### HTTP服务器

#### 非阻塞模式（适合后台服务）

```mml
fn void main() {
    // 创建HTTP服务器
    let server_id = network.create_server(8080);
    
    // 添加路由
    network.add_route(server_id, "GET", "/", "欢迎使用MemMole HTTP服务器！");
    network.add_route(server_id, "GET", "/hello", "Hello, World! 来自MemMole服务器");
    network.add_route(server_id, "GET", "/api/info", "{\"name\":\"MemMole Server\",\"version\":\"1.0\"}");
    
    // 启动服务器（非阻塞）
    let result = network.start_server(server_id);
    
    // 检查状态
    let is_running = network.is_server_running(server_id);
    let port = network.get_server_port(server_id);
    
    return 0;
}
```
```

#### 阻塞模式（适合长期运行的服务）

```mml
fn void main() {
    // 创建HTTP服务器
    let server_id = network.create_server(8080);
    
    // 添加路由
    network.add_route(server_id, "GET", "/", "欢迎使用MemMole阻塞式HTTP服务器！");
    network.add_route(server_id, "GET", "/hello", "Hello, World! 服务器正在运行中...");
    network.add_route(server_id, "GET", "/status", "服务器状态：运行中");
    
    // 启动服务器并阻塞等待（程序会保持运行）
    let result = network.start_server_and_wait(server_id);
    
    // 这行代码只有在服务器停止后才会执行
    return 0;
}
```
```

### 循环和条件

```mml
fn int factorial(int n) {
    let result = 1;
    let i = 1;
    
    while (i <= n) {
        result = result * i;
        i = i + 1;
    }
    
    return result;
}

fn void main() {
    let n = 5;
    let fact = factorial(n);
    
    if (fact > 100) {
        return 1;
    } else {
        return 0;
    }
}
```
```

## 🔧 开发指南

### 添加新的语法特性

1. **词法分析器**：在 `lexer/lexer.go` 中添加新的Token类型
2. **语法分析器**：在 `parser/` 目录中添加相应的解析逻辑
3. **AST节点**：在 `ast/` 目录中定义新的节点类型
4. **求值器**：在 `evaluator/` 目录中实现求值逻辑

### 添加新的内置函数

1. 在 `builtins/` 目录中创建新的功能模块
2. 在 `evaluator/expressions.go` 中添加调用逻辑
3. 更新文档和示例

### 代码风格

- 遵循Go语言的代码规范
- 使用清晰的函数和变量命名
- 添加适当的注释
- 保持模块间的低耦合

## 🧪 测试

运行测试：
```bash
go test ./...
```

运行特定包的测试：
```bash
go test ./lexer
go test ./parser
```

## 📄 许可证

本项目采用 MIT 许可证。

## 🤝 贡献

欢迎提交Issue和Pull Request！

### 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 📞 联系方式

如有问题或建议，请通过以下方式联系：

- 提交 Issue
- 发送邮件
- 参与讨论

---

**CVM** - 让编程更简单，让学习更有趣！ 🚀
