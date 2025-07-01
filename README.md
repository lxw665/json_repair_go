# JSON Repair Go

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/lxw/json_repair_go)](https://goreportcard.com/report/github.com/lxw/json_repair_go)

一个用于修复无效 JSON 字符串的 Go 库，特别适用于处理大语言模型（LLM）输出的有问题的 JSON 数据。

**English** | [中文](README_CN.md)

## 特性

- 🔧 **智能修复**：自动修复常见的 JSON 语法错误
- 🚀 **高性能**：基于高效的解析器，性能优异
- 📦 **零依赖**：仅使用 Go 标准库，无外部依赖
- 🛡️ **类型安全**：提供强类型的 API 接口
- 🎯 **LLM 友好**：专门针对 LLM 输出的常见错误进行优化
- 🔄 **兼容性**：可作为 `json.Unmarshal` 的直接替代品

## 支持的修复类型

- ✅ 缺失的引号（键名和字符串值）
- ✅ 单引号转双引号
- ✅ 缺失的括号和方括号
- ✅ 缺失的逗号
- ✅ 尾随逗号
- ✅ Python 风格的布尔值（`True`/`False`）
- ✅ Python 风格的空值（`None`）
- ✅ 注释处理（`//`、`/* */`、`#`）
- ✅ 转义字符修复
- ✅ 数字格式修复

## 安装

```bash
go get github.com/lxw665/json_repair_go
```

## 快速开始

### 基本用法

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/lxw/json_repair_go"
)

func main() {
    // 修复破损的 JSON 字符串
    broken := `{name: 'John', age: 30, active: True`
    
    // 方法 1: 修复并返回 JSON 字符串
    repaired, err := jsonrepair.RepairJSON(broken)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(repaired) // 输出: {"name":"John","age":30,"active":true}
    
    // 方法 2: 直接解析为 Go 对象
    var person map[string]interface{}
    err = jsonrepair.Unmarshal([]byte(broken), &person)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("%+v\n", person) // 输出: map[active:true age:30 name:John]
    
    // 方法 3: 使用 Loads 函数
    result, err := jsonrepair.Loads(broken)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("%+v\n", result)
}
```

### 高级用法

```go
// 使用选项
options := jsonrepair.RepairOptions{
    SkipValidation: true,  // 跳过输入验证以提高性能
    Logging:        true,  // 启用详细日志
}

repaired, err := jsonrepair.RepairJSONWithOptions(broken, options)
```

## 修复示例

| 输入 | 输出 |
|------|------|
| `{name: 'John', age: 30` | `{"name":"John","age":30}` |
| `{'key': 'value',}` | `{"key":"value"}` |
| `[1, 2, 3,` | `[1,2,3]` |
| `{"active": True}` | `{"active":true}` |
| `{"value": None}` | `{"value":null}` |
| `{a: 1, b: 2` | `{"a":1,"b":2}` |

## API 文档

### 核心函数

#### `RepairJSON(jsonStr string) (string, error)`
修复 JSON 字符串并返回有效的 JSON。

#### `RepairJSONWithOptions(jsonStr string, options RepairOptions) (string, error)`
使用自定义选项修复 JSON 字符串。

#### `Loads(jsonStr string) (interface{}, error)`
解析（并修复）JSON 字符串，返回 Go 对象。

#### `Unmarshal(data []byte, v interface{}) error`
类似于 `json.Unmarshal`，但会在解析前自动修复 JSON。

#### `Load(reader io.Reader, v interface{}) error`
从 `io.Reader` 读取并解析 JSON 数据。

### 选项配置

```go
type RepairOptions struct {
    SkipValidation bool // 跳过输入 JSON 有效性检查
    ReturnObjects  bool // 返回解析后的对象而非 JSON 字符串
    Logging        bool // 启用详细的修复日志
    StreamStable   bool // 为流式 JSON 保持稳定的修复结果
    EnsureASCII    bool // 确保非 ASCII 字符被转义
}
```

## 性能

基准测试结果：

```
BenchmarkRepairJSON-8     100000    12034 ns/op     2048 B/op    15 allocs/op
BenchmarkLoads-8          150000     8567 ns/op     1536 B/op    12 allocs/op
```

## 使用场景

### 1. LLM 输出处理
```go
// ChatGPT、Claude 等模型有时会输出格式不完整的 JSON
llmOutput := `{"response": "Hello", "confidence": 0.95`
data, _ := jsonrepair.Loads(llmOutput)
```

### 2. 配置文件修复
```go
// 修复手写的配置文件
config := `{
    name: 'MyApp',
    version: '1.0.0',
    debug: True,
}`
var cfg Config
jsonrepair.Unmarshal([]byte(config), &cfg)
```

### 3. API 响应处理
```go
// 处理第三方 API 返回的非标准 JSON
response := `{'status': 'ok', 'data': {...}}`
var result APIResponse
jsonrepair.Unmarshal([]byte(response), &result)
```

## 测试

运行测试套件：

```bash
# 运行所有测试
go test -v

# 运行基准测试
go test -bench=.

# 查看测试覆盖率
go test -cover
```

## 与其他语言的实现

此项目参考了以下优秀的实现：

- **Python**: [json_repair](https://github.com/mangiucugna/json_repair) - 原始实现
- **TypeScript**: [jsonrepair](https://github.com/josdejong/jsonrepair)
- **Rust**: [llm_json](https://github.com/oramasearch/llm_json)
- **Ruby**: [json-repair-rb](https://github.com/sashazykov/json-repair-rb)

## 贡献

欢迎提交 Issue 和 Pull Request！

1. Fork 这个仓库
2. 创建你的特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交你的改动 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建一个 Pull Request

## 开发路线图

- [ ] 支持更多的 JSON 修复场景
- [ ] 添加流式处理支持
- [ ] 优化性能和内存使用
- [ ] 添加命令行工具
- [ ] 支持自定义修复规则
- [ ] 添加更多语言的示例

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 致谢

- 感谢 [mangiucugna/json_repair](https://github.com/mangiucugna/json_repair) 提供的原始 Python 实现
- 感谢所有为 JSON 修复生态系统做出贡献的开发者

---

如果这个项目对你有帮助，请给个 ⭐️ 支持一下！
