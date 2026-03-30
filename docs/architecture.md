# MemMole 项目架构文档

## 项目概述

MemMole 是一个用Go语言实现的自定义虚拟机语言，采用经典的编译器架构设计，支持变量声明、函数定义、控制流、网络操作等功能。

## 整体架构流程

```mermaid
graph TD
    A[源代码 .mml文件] --> B[词法分析器 Lexer]
    B --> C[Token序列]
    C --> D[语法分析器 Parser]
    D --> E[抽象语法树 AST]
    E --> F[语义分析器]
    F --> G[环境管理器 Environment]
    G --> H[求值器 Evaluator]
    H --> I[对象系统 Object System]
    I --> J[执行结果]
    
    K[包管理器 PackageManager] --> D
    L[内置函数 Builtins] --> H
    M[存储系统 Storage] --> H
    N[网络模块 Network] --> H
    O[数据库模块 Database] --> H
```

## 核心组件详解

### 1. 词法分析器 (Lexer)

```mermaid
graph LR
    A[源代码字符串] --> B[字符扫描]
    B --> C[Token识别]
    C --> D[Token流]
    
    E[关键字识别] --> C
    F[标识符识别] --> C
    G[字面量识别] --> C
    H[运算符识别] --> C
```

**功能**:
- 将源代码转换为Token序列
- 支持关键字: `fn`, `if`, `else`, `while`, `return`, `package`, `import`
- 支持字面量: 整数、字符串、布尔值
- 支持运算符: `+`, `-`, `*`, `/`, `==`, `!=`, `>`, `<`
- 错误位置跟踪

### 2. 语法分析器 (Parser)

```mermaid
graph TD
    A[Token序列] --> B[递归下降解析]
    B --> C[AST节点构建]
    C --> D[语法树]
    
    E[表达式解析] --> C
    F[语句解析] --> C
    G[函数解析] --> C
    H[类型解析] --> C
```

**解析策略**:
- 递归下降解析 (Recursive Descent Parsing)
- 运算符优先级解析 (Operator Precedence Parsing)
- 前缀和中缀解析函数映射

**优先级层次**:
```
LOWEST < ASSIGNMENT < COMPARE < SUM < PRODUCT < PREFIX < CALL < MEMBER
```

### 3. 抽象语法树 (AST)

```mermaid
graph TD
    A[Program 根节点] --> B[Statement 语句节点]
    A --> C[Expression 表达式节点]
    
    B --> D[LetStatement 变量声明]
    B --> E[ReturnStatement 返回语句]
    B --> F[ExpressionStatement 表达式语句]
    B --> G[FunctionStatement 函数定义]
    B --> H[IfStatement 条件语句]
    B --> I[WhileStatement 循环语句]
    
    C --> J[Identifier 标识符]
    C --> K[IntegerLiteral 整数字面量]
    C --> L[StringLiteral 字符串字面量]
    C --> M[BooleanLiteral 布尔字面量]
    C --> N[FunctionLiteral 函数字面量]
    C --> O[CallExpression 函数调用]
    C --> P[InfixExpression 中缀表达式]
    C --> Q[PrefixExpression 前缀表达式]
```

### 4. 求值器 (Evaluator)

```mermaid
graph TD
    A[AST节点] --> B[节点类型匹配]
    B --> C[对应求值函数]
    C --> D[对象实例]
    
    E[环境管理] --> C
    F[内置函数调用] --> C
    G[类型检查] --> C
    H[错误处理] --> C
```

**求值策略**:
- 树遍历求值 (Tree Walking Evaluation)
- 环境链作用域管理
- 函数闭包支持

### 5. 对象系统 (Object System)

```mermaid
graph TD
    A[Object 接口] --> B[Integer 整数对象]
    A --> C[String 字符串对象]
    A --> D[Boolean 布尔对象]
    A --> E[Null 空对象]
    A --> F[Function 函数对象]
    A --> G[ReturnValue 返回值对象]
    A --> H[Error 错误对象]
```

## 数据流程图

### 完整执行流程

```mermaid
sequenceDiagram
    participant Main as main.go
    participant Lexer as lexer
    participant Parser as parser
    participant Evaluator as evaluator
    participant Env as environment
    participant Builtins as builtins
    
    Main->>Lexer: New(input)
    Lexer->>Lexer: 词法分析
    Lexer-->>Main: Token流
    
    Main->>Parser: New(lexer)
    Parser->>Parser: 语法分析
    Parser-->>Main: AST
    
    Main->>Env: NewEnvironment()
    Main->>Builtins: 注册内置包
    Builtins-->>Env: 设置网络、数据库包
    
    Main->>Evaluator: Eval(AST, Env)
    Evaluator->>Evaluator: 递归求值
    Evaluator-->>Main: 执行结果
```

### 解析树构建过程

```mermaid
graph TD
    A[源代码: let x = 5 + 3;] --> B[Token序列]
    B --> C[AST构建]
    
    subgraph "Token序列"
        B1[LET]
        B2[IDENTIFIER 'x']
        B3[ASSIGN '=']
        B4[INTEGER '5']
        B5[PLUS '+']
        B6[INTEGER '3']
        B7[SEMICOLON ';']
    end
    
    subgraph "AST结构"
        C1[LetStatement]
        C1 --> C2[Identifier: 'x']
        C1 --> C3[InfixExpression]
        C3 --> C4[IntegerLiteral: 5]
        C3 --> C5[Operator: '+']
        C3 --> C6[IntegerLiteral: 3]
    end
```

## 模块依赖关系

```mermaid
graph TD
    A[main.go] --> B[lexer]
    A --> C[parser]
    A --> D[evaluator]
    A --> E[repl]
    
    C --> F[ast]
    C --> B
    D --> F
    D --> G[object]
    D --> C
    
    H[builtins] --> D
    I[network] --> H
    J[database] --> H
    K[storage] --> J
    L[config] --> H
    M[logger] --> A
    
    C --> N[package_manager]
    C --> O[environment]
```

## 包管理架构

```mermaid
graph TD
    A[PackageManager] --> B[包解析]
    B --> C[Import语句处理]
    C --> D[模块加载]
    D --> E[环境合并]
    
    F[本地包] --> D
    G[内置包] --> D
    H[第三方包] --> D
```

## 错误处理流程

```mermaid
graph TD
    A[错误发生] --> B[错误类型判断]
    B --> C[词法错误]
    B --> D[语法错误]
    B --> E[运行时错误]
    
    C --> F[位置报告]
    D --> G[语法提示]
    E --> H[堆栈跟踪]
    
    F --> I[错误对象]
    G --> I
    H --> I
    I --> J[用户输出]
```

## 性能优化点

1. **词法分析优化**
   - 字符预读减少重复扫描
   - Token缓存机制

2. **语法分析优化**
   - 运算符优先级表快速查找
   - 解析函数映射表

3. **求值优化**
   - 环境链缓存
   - 内置函数快速分发

## 扩展性设计

1. **新语法特性添加**
   - Lexer: 添加新Token类型
   - Parser: 添加解析函数
   - AST: 添加新节点类型
   - Evaluator: 添加求值逻辑

2. **新内置函数**
   - 在builtins包中实现
   - 注册到环境管理器

3. **新存储后端**
   - 实现Storage接口
   - 注册到数据库包

## 总结

MemMole 采用经典的三阶段编译器架构，具有良好的模块化设计和扩展性。通过清晰的分层架构，使得各个组件职责明确，便于维护和扩展。
