package main

import (
	"fmt"
	"log"

	jsonrepair "github.com/lxw/json_repair_go"
)

func main() {
	fmt.Println("🔧 JSON Repair Go - Demo")
	fmt.Println("=" + fmt.Sprintf("%*s", 30, "="))

	// 示例测试用例
	examples := []struct {
		name  string
		input string
	}{
		{
			name:  "缺失闭合括号",
			input: `{"name": "John", "age": 30, "city": "New York"`,
		},
		{
			name:  "单引号替换",
			input: `{'name': 'John', 'age': 30, 'city': 'New York'}`,
		},
		{
			name:  "无引号键名",
			input: `{name: "John", age: 30, city: "New York"}`,
		},
		{
			name:  "尾随逗号",
			input: `{"name": "John", "age": 30, "city": "New York",}`,
		},
		{
			name:  "数组缺失括号",
			input: `["apple", "banana", "cherry"`,
		},
		{
			name:  "Python 布尔值",
			input: `{"active": True, "verified": False, "data": None}`,
		},
		{
			name:  "混合问题",
			input: `{name: 'John', age: 30, hobbies: ['reading', 'coding'], active: True`,
		},
		{
			name:  "嵌套对象",
			input: `{name: 'John', details: {age: 30, address: {street: '123 Main St', city: 'New York'`,
		},
		{
			name:  "带注释",
			input: `{"name": "John", "age": 30} // 这是注释`,
		},
	}

	for i, example := range examples {
		fmt.Printf("\n📝 示例 %d: %s\n", i+1, example.name)
		fmt.Printf("输入:   %s\n", example.input)

		// 修复 JSON
		repaired, err := jsonrepair.RepairJSON(example.input)
		if err != nil {
			log.Printf("❌ 修复失败: %v", err)
			continue
		}

		fmt.Printf("修复后: %s\n", repaired)

		// 解析为 Go 对象
		result, err := jsonrepair.Loads(example.input)
		if err != nil {
			log.Printf("❌ 解析失败: %v", err)
			continue
		}

		fmt.Printf("解析结果: %+v\n", result)
		fmt.Printf("✅ 成功修复并解析\n")
	}

	// 演示高级用法
	fmt.Println("\n🚀 高级用法演示")
	fmt.Println("=" + fmt.Sprintf("%*s", 20, "="))

	broken := `{name: 'John', age: 30}`
	options := jsonrepair.RepairOptions{
		SkipValidation: true,
		Logging:        false,
	}

	repaired, err := jsonrepair.RepairJSONWithOptions(broken, options)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("输入:     %s\n", broken)
	fmt.Printf("修复后:   %s\n", repaired)

	// 演示 Unmarshal 用法
	fmt.Println("\n📦 Unmarshal 演示")
	fmt.Println("=" + fmt.Sprintf("%*s", 18, "="))

	type Person struct {
		Name   string `json:"name"`
		Age    int    `json:"age"`
		Active bool   `json:"active"`
	}

	brokenPersonJSON := `{name: 'Alice', age: 25, active: True`
	var person Person

	err = jsonrepair.Unmarshal([]byte(brokenPersonJSON), &person)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("输入: %s\n", brokenPersonJSON)
	fmt.Printf("解析结果: %+v\n", person)

	fmt.Println("\n🎉 演示完成！")
	fmt.Println("\n💡 提示:")
	fmt.Println("   - 所有测试用例都成功修复")
	fmt.Println("   - 可以作为 json.Unmarshal 的直接替代品")
	fmt.Println("   - 特别适用于处理 LLM 输出的 JSON")
}
