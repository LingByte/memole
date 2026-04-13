# 官方模块扩展指南（Go 实现）

这份文档回答一个问题：如果要给 Memole 增加“官方模块”，应该如何用 Go 落地。

当前仓库已落地一个最小官方模块：`std`（Go 实现）。

---

## 当前机制（先认识现状）

现在 `import` 采用的是 **`.mml` 文件模块加载**：

- 在 `pkg/evaluator/module_loader.go` 中读取模块文件并执行
- 导出变量/函数存入 `ModuleObject.Exports`
- 通过 `mod.member` 访问

这适合纯脚本模块，但“官方模块”通常更希望由 Go 提供能力（例如时间、随机数、文件、网络）。

---

## 推荐方案：内置 Go 模块注册表

建议新增一层“官方模块注册表”：

1. 先查内置注册表（Go 实现）
2. 未命中再走 `.mml` 文件模块加载

这样可以同时支持：

- `import std;`（Go 官方模块）
- `import mylib;`（用户 `.mml` 模块）

---

## 目录建议

建议新增目录：

- `pkg/builtinmodules/registry.go`：注册表和查找逻辑
- `pkg/builtinmodules/std.go`：示例官方模块
- （可选）`pkg/object/native_function.go`：原生函数对象定义

---

## 最小实现步骤

### 1. 定义原生函数对象（Go 可调用）

在 `pkg/object` 增加一种 `NativeFunction` 类型，核心是：

- 入参：`[]object.Object`
- 出参：`object.Object`

这样 evaluator 在调用函数时，既能调用 Memole `FunctionObject`，也能调用 Go `NativeFunction`。

### 2. 实现模块注册表

在 `pkg/builtinmodules/registry.go` 中维护：

- `map[string]*object.ModuleObject`
- `Get(name string) (*object.ModuleObject, bool)`

### 3. 在 import 处优先检查官方模块

在 `pkg/evaluator/statements.go` 的 `evalImportStatement` 里：

1. `builtinmodules.Get(stmt.Path.Value)`
2. 命中则直接 `env.Set(aliasOrName, mod)` 并返回
3. 未命中再走现有 `module_loader.go`

### 4. 在调用表达式支持 NativeFunction

在 `pkg/evaluator/expressions.go` 的 `evalCallExpression` 中，增加：

- `if nf, ok := function.(*object.NativeFunction); ok { return nf.Fn(args) }`

### 5. 编写第一个官方模块（示例）

例如 `std`：

- `std.nowUnix()` -> 返回当前时间戳（Integer）
- `std.typeOf(x)` -> 返回对象类型字符串

---

## 示例：你期望的用户体验

新增官方模块后，用户可以直接写：

```mml
import std;

fn main() {
    return std.nowUnix();
}
```

当前实现可直接使用：

```mml
import std;

fn main() {
    std.println("hello");
    std.println(std.typeOf(123));
    return std.nowUnix();
}
```

运行方式保持文件执行：

```bash
go run . run app.mml
```

---

## 设计建议（避免后续返工）

- **命名空间隔离**：官方模块统一 `std*` 前缀，避免和用户模块同名
- **错误对象统一**：原生函数参数错误应返回统一错误对象（建议后续补 `ERROR_OBJ`）
- **类型检查前置**：每个原生函数自行校验参数数量和类型
- **文档同步**：每加一个官方模块，同步更新 `docs/interpreter-runtime.md` 与示例

