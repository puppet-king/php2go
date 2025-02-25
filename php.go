package php2go

import (
	"fmt"
	"strings"
)

// Array array()
func Array[T any](v ...T) []T {
	arr := make([]T, len(v))
	copy(arr, v)
	return arr
}

// ArrayChangeKeyCase array_change_key_case
func ArrayChangeKeyCase[T any](input map[string]T, toUpper bool) map[string]T {
	output := make(map[string]T)

	for k, v := range input {
		if toUpper {
			output[strings.ToUpper(k)] = v
		} else {
			output[strings.ToLower(k)] = v
		}
	}

	return output
}

// ArrayChunk 将输入的切片 v 分割成多个大小为 size 的切片, preserve_keys 为 false 即不保留原来的键名
func ArrayChunk[T any](v []T, size int) [][]T {
	chunks := make([][]T, 0)
	//if size <= 0 {
	//	return chunks
	//}

	// 遍历切片，将其分割成多个块
	for i := 0; i < len(v); i += size {
		end := i + size
		if end > len(v) {
			end = len(v)
		}
		chunks = append(chunks, v[i:end])
	}

	return chunks
}

// ArrayColumn 返回输入数组中指定列的值
func ArrayColumn[T any](input []map[string]T, columnKey string) ([]T, error) {
	// 用来存储结果的 map，键为 indexKey 的值，值为相应列的所有值
	result := make([]T, 0, len(input))

	// 遍历切片
	for i, r := range input {
		// 获取 columnKey 对应的值
		columnValue, ok := r[columnKey]
		if !ok {
			return nil, fmt.Errorf("columnKey '%s' not found in row %d", columnKey, i)
		}

		result = append(result, columnValue)
	}

	return result, nil
}

// ArrayColumnByIndexKey  返回输入数组中指定列的值, 支持 indexKey 实现按照 indexKey 分类聚合的功能
func ArrayColumnByIndexKey[T any](input []map[string]T, columnKey string, indexKey *string) (map[string][]T, error) {
	// 用来存储结果的 map，键为 indexKey 的值，值为相应列的所有值
	result := make(map[string][]T)

	// 遍历切片
	for i, r := range input {
		// 获取 columnKey 对应的值
		columnValue, ok := r[columnKey]
		if !ok {
			return nil, fmt.Errorf("columnKey '%s' not found in row %d", columnKey, i)
		}

		// 如果提供了 indexKey，检查该键是否存在并获取其值
		if indexKey != nil {
			indexValue, ok := r[*indexKey]
			if !ok {
				return nil, fmt.Errorf("indexKey '%s' not found in row %d", *indexKey, i)
			}

			// 将数据按 indexKey 值分组，indexKey 的值作为 map 的键
			indexKeyStr := fmt.Sprintf("%v", indexValue) // 将 indexValue 转换为字符串，作为 map 的键
			result[indexKeyStr] = append(result[indexKeyStr], columnValue)
		} else {
			// 如果没有提供 indexKey，则按原来逻辑直接添加到 result
			result["default"] = append(result["default"], columnValue)
		}
	}

	return result, nil
}
