# JSON Repair Go （中文）

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/lxw665/json_repair_go)](https://goreportcard.com/report/github.com/lxw665/json_repair_go)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/lxw665/json_repair_go)](https://pkg.go.dev/github.com/lxw665/json_repair_go)

一个用于修复无效 JSON 字符串的 Go 库，特别适合清理来自大语言模型（LLM）或手写配置中的非标准 JSON。

[English](README.md) | 中文

## 特性

- 修复键名/字符串缺少或错误的引号
- 单引号自动转换为双引号
- 自动补齐缺失的逗号、右花括号/方括号
- 移除尾随逗号
- 规范化布尔和空值（`True/False` ➜ `true/false`，`None` ➜ `null`）
- 去除注释（`//`、`/* */`、`#`）
- 轻量的转义与数字解析修复
- 零外部依赖（仅标准库）

## 安装

```bash
go get github.com/lxw665/json_repair_go
```

注意：请确保你的导入路径与仓库路径一致。如果你在自己的账号或组织下发布，请把 `lxw665` 替换为你的 GitHub 用户名，并同步更新 go.mod 的 module 路径。

## 快速开始

```go
package main

import (
    "fmt"
    "log"
    jsonrepair "github.com/lxw665/json_repair_go"
)

func main() {
    broken := `{name: 'John', age: 30, active: True`

    // 1) 修复为 JSON 字符串
    fixed, err := jsonrepair.RepairJSON(broken)
    if err != nil { log.Fatal(err) }
    fmt.Println(fixed) // {"name":"John","age":30,"active":true}

    // 2) 直接解析为 Go 值
    v, err := jsonrepair.Loads(broken)
    if err != nil { log.Fatal(err) }
    fmt.Printf("%+v\n", v)
}
```

## API

- `RepairJSON(jsonStr string) (string, error)` — 修复并返回 JSON 字符串
- `RepairJSONWithOptions(jsonStr string, options RepairOptions) (string, error)` — 携带选项的修复
- `Loads(jsonStr string) (interface{}, error)` — 修复并解析为 Go 值
- `Unmarshal(data []byte, v any) error` — 在 `json.Unmarshal` 之上自动修复
- `Load(r io.Reader, v any) error`

选项结构：
```go
type RepairOptions struct {
    SkipValidation bool // 跳过初始的 json.Unmarshal 校验
    ReturnObjects  bool // 返回对象而不是 JSON 字符串
    Logging        bool // 打开修复日志
    StreamStable   bool // 面向流式场景的稳定修复
    EnsureASCII    bool // 转义非 ASCII 字符
}
```

## 示例

输入 ➜ 输出：

- `{name: 'John', age: 30` ➜ `{"name":"John","age":30}`
- `{'key': 'value',}` ➜ `{"key":"value"}`
- `[1, 2, 3,` ➜ `[1,2,3]`
- `{"active": True}` ➜ `{"active":true}`
- `{"value": None}` ➜ `{"value":null}`
- `{a: 1, b: 2` ➜ `{"a":1,"b":2}`

更多可运行示例在 `demo/`：
```bash
cd demo
go run .
```

## 测试

```bash
go test -v
go test -bench=.
go test -cover
```

## 相关实现

- Python: https://github.com/mangiucugna/json_repair
- TypeScript: https://github.com/josdejong/jsonrepair
- Rust: https://github.com/oramasearch/llm_json
- Ruby: https://github.com/sashazykov/json-repair-rb

## 贡献

欢迎 PR 和 Issue。

1) Fork ➜ 2) 新建分支 ➜ 3) 提交 ➜ 4) 推送 ➜ 5) 提 PR

## 许可证

MIT，详见 [LICENSE](LICENSE)。

---

如果这个项目对你有帮助，欢迎点个 ⭐️ 支持！