package main

import (
	"fmt"
	"php2go"
)

func main() {
	// 示例数据：每个元素是 map[string]interface{}
	data := []map[string]interface{}{
		{"id": 1, "name": "Alice", "age": 30},
		{"id": 2, "name": "Bob", "age": 25},
		{"id": 3, "name": "Charlie", "age": 35},
		{"id": 4, "name": "David", "age": 40},
	}

	//data := []map[string]interface{}{
	//	{"id": 1, "name": "Alice", "age": 30},
	//}
	//
	// 获取 "name" 列
	columnKey := "name"
	//var indexKey *string // 如果不需要 indexKey，可以传 nil
	//indexKey := "id"
	result, err := php2go.ArrayColumn(data, columnKey)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(result)

	// 打印结果
	for _, value := range result {
		fmt.Println(value)
	}
}
