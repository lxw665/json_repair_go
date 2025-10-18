# JSON Repair Go

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/lxw665/json_repair_go)](https://goreportcard.com/report/github.com/lxw665/json_repair_go)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/lxw665/json_repair_go)](https://pkg.go.dev/github.com/lxw665/json_repair_go)

Repair invalid JSON strings in Go. Especially handy for cleaning up JSON coming from Large Language Models (LLMs) or hand-written configs.

English | [中文](README_CN.md)

## Features

- Fix missing/wrong quotes on keys and strings
- Convert single quotes to double quotes
- Add missing commas, closing braces/brackets
- Remove trailing commas
- Normalize booleans (`True`/`False` ➜ `true`/`false`) and nulls (`None` ➜ `null`)
- Strip comments (`//`, `/* */`, `#`)
- Light escaping and number parsing fixes
- Zero external deps (stdlib only)

## Install

```bash
go get github.com/lxw665/json_repair_go
```

Note: make sure your import path matches the repository path you publish under. If you fork, replace `lxw665` with your GitHub username.

## Quick start

```go
package main

import (
    "fmt"
    "log"
    jsonrepair "github.com/lxw665/json_repair_go"
)

func main() {
    broken := `{name: 'John', age: 30, active: True`

    // 1) Repair to JSON string
    fixed, err := jsonrepair.RepairJSON(broken)
    if err != nil { log.Fatal(err) }
    fmt.Println(fixed) // {"name":"John","age":30,"active":true}

    // 2) Parse directly to Go value
    v, err := jsonrepair.Loads(broken)
    if err != nil { log.Fatal(err) }
    fmt.Printf("%+v\n", v)
}
```

## API

- `RepairJSON(jsonStr string) (string, error)` — Repair and return a JSON string
- `RepairJSONWithOptions(jsonStr string, options RepairOptions) (string, error)`
- `Loads(jsonStr string) (interface{}, error)` — Repair and parse to Go values
- `Unmarshal(data []byte, v any) error` — Drop-in for `json.Unmarshal` with repair
- `Load(r io.Reader, v any) error`

Options:
```go
type RepairOptions struct {
    SkipValidation bool // Skip initial json.Unmarshal check
    ReturnObjects  bool // Return parsed objects instead of JSON string
    Logging        bool // Enable repair logs
    StreamStable   bool // Keep repair stable for streaming use-cases
    EnsureASCII    bool // Escape non-ASCII
}
```

## Examples

Input ➜ Output

- `{name: 'John', age: 30` ➜ `{"name":"John","age":30}`
- `{'key': 'value',}` ➜ `{"key":"value"}`
- `[1, 2, 3,` ➜ `[1,2,3]`
- `{"active": True}` ➜ `{"active":true}`
- `{"value": None}` ➜ `{"value":null}`
- `{a: 1, b: 2` ➜ `{"a":1,"b":2}`

More runnable examples in `demo/`:
```bash
cd demo
go run .
```

## Testing

```bash
go test -v
go test -bench=.
go test -cover
```

## Related work

- Python: https://github.com/mangiucugna/json_repair
- TypeScript: https://github.com/josdejong/jsonrepair
- Rust: https://github.com/oramasearch/llm_json
- Ruby: https://github.com/sashazykov/json-repair-rb

## Contributing

PRs and issues are welcome.

1) Fork ➜ 2) create branch ➜ 3) commit ➜ 4) push ➜ 5) open PR

## License

MIT. See [LICENSE](LICENSE).

---

If this project helps you, please give it a ⭐️!
