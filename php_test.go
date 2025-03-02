package php2go

import (
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArray(t *testing.T) {
	type testCase[T any] struct {
		name string
		args []T
		want []T
	}

	var tests = []testCase[any]{
		{
			name: "Test with int slice",
			args: []any{1, 2, 3, 4, 5}, // 正确初始化 args[int] 类型
			want: []any{1, 2, 3, 4, 5},
		},
		{
			name: "Test with string slice",
			args: []any{"a", "b", "c"},
			want: []any{"a", "b", "c"},
		},
		{
			name: "Test with mixed types",
			args: []any{1, "a", true, 3.14},
			want: []any{1, "a", true, 3.14},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Array(tt.args...) // 展开切片传入变长参数
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Array() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArrayChangeKeyCase(t *testing.T) {
	type args struct {
		input   map[string]any
		toUpper bool
	}
	tests := []struct {
		name string
		args args
		want map[string]any
	}{
		{
			name: "Test with map to upper",
			args: args{
				input: map[string]any{
					"a": 1,
					"b": 2,
				},
				toUpper: true,
			},
			want: map[string]any{
				"A": 1,
				"B": 2,
			},
		},
		{
			name: "Test with map to lower",
			args: args{
				input: map[string]any{
					"A": 1,
					"B": 2,
				},
				toUpper: false,
			},
			want: map[string]any{
				"a": 1,
				"b": 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ArrayChangeKeyCase(tt.args.input, tt.args.toUpper); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArrayChangeKeyCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArrayChunk(t *testing.T) {
	type args[T any] struct {
		v    []T
		size int
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want [][]T
	}
	tests := []testCase[any]{
		{name: "Test with size 2", args: args[any]{v: []any{1, 2, 3, 4, 5}, size: 2}, want: [][]any{{1, 2}, {3, 4}, {5}}},
		{name: "Test with empty slice", args: args[any]{v: []any{}, size: 2}, want: [][]any{}},
		{name: "Test with overflow ", args: args[any]{v: []any{1, 2, 3}, size: 10}, want: [][]any{{1, 2, 3}}},
		//{name: "Test with negative size", args: args[any]{v: []any{1, 2, 3}, size: -1}, want: [][]any{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ArrayChunk(tt.args.v, tt.args.size); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArrayChunk() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArrayColumn(t *testing.T) {
	type args[T any] struct {
		input     []map[string]T
		columnKey string
	}
	type testCase[T any] struct {
		name    string
		args    args[T]
		want    []T
		wantErr bool
	}
	tests := []testCase[any]{
		{
			name: "Test with column_key exist", args: args[any]{
			input: []map[string]interface{}{
				{"id": 1, "name": "Alice", "age": 30},
				{"id": 2, "name": "Bob", "age": 25},
			},
			columnKey: "name",
		},
			want:    []any{"Alice", "Bob"},
			wantErr: false,
		},
		{
			name: "Test with column_key not exist", args: args[any]{
			input: []map[string]interface{}{
				{"id": 1, "name": "Alice", "age": 30},
				{"id": 2, "name": "Bob", "age": 25},
			},
			columnKey: "cc",
		},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ArrayColumn(tt.args.input, tt.args.columnKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("ArrayColumn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArrayColumn() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestArrayColumnByIndexKey(t *testing.T) {
	type args[T any] struct {
		input     []map[string]T
		columnKey string
		indexKey  *string
	}
	type testCase[T any] struct {
		name    string
		args    args[T]
		want    map[string][]T
		wantErr bool
	}

	ageKey := "age"
	notExistIndexKey := "notExistIndexKey"

	tests := []testCase[any]{
		{
			name: "Test with columnKey and indexKey", args: args[any]{
			input: []map[string]interface{}{
				{"id": 1, "name": "Alice", "age": 30},
				{"id": 2, "name": "Bob", "age": 25},
				{"id": 3, "name": "Charlie", "age": 30},
			},
			columnKey: "name",
			indexKey:  &ageKey,
		},
			want: map[string][]any{
				"30": {"Alice", "Charlie"},
				"25": {"Bob"},
			},
			wantErr: false,
		},
		{
			name: "Test with columnKey only (no indexKey)",
			args: args[any]{
				input: []map[string]interface{}{
					{"id": 1, "name": "Alice", "age": 30},
					{"id": 2, "name": "Bob", "age": 25},
				},
				columnKey: "name",
			},
			want: map[string][]any{
				"default": {"Alice", "Bob"},
			},
			wantErr: false,
		},
		{
			name: "Test with non-existing columnKey",
			args: args[any]{
				input: []map[string]interface{}{
					{"id": 1, "name": "Alice", "age": 30},
				},
				columnKey: "nonExistingKey", // 该列不存在
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Test with non-existing indexKey",
			args: args[any]{
				input: []map[string]interface{}{
					{"id": 1, "name": "Alice", "age": 30},
					{"id": 2, "name": "Bob", "age": 25},
				},
				columnKey: "name",
				indexKey:  &notExistIndexKey,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ArrayColumnByIndexKey(tt.args.input, tt.args.columnKey, tt.args.indexKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("ArrayColumnByIndexKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArrayColumnByIndexKey() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetcwd(t *testing.T) {
	dir, err := Getcwd()

	// 断言错误为 nil
	require.NoError(t, err)

	// 断言返回的目录不为空
	assert.NotEmpty(t, dir)
}

func TestIsDir(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "test_dir")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// 测试目录存在且是目录
	isDir := IsDir(tempDir)
	assert.True(t, isDir)

	// 测试文件不是目录
	tempFile := filepath.Join(tempDir, "test_file.txt")
	_, err = os.Create(tempFile)
	assert.NoError(t, err)

	isDir = IsDir(tempFile)
	assert.False(t, isDir)

	// 测试不存在的路径
	nonExistentPath := filepath.Join(tempDir, "non_existent")
	isDir = IsDir(nonExistentPath)
	assert.False(t, isDir)
}

func TestIsFile(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "test_dir")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// 测试目录存在且是目录
	isDir := IsFile(tempDir)
	assert.False(t, isDir)

	// 测试文件不是目录
	tempFile := filepath.Join(tempDir, "test_file.txt")
	_, err = os.Create(tempFile)
	assert.NoError(t, err)

	isDir = IsFile(tempFile)
	assert.True(t, isDir)

	// 测试不存在的路径
	nonExistentPath := filepath.Join(tempDir, "non_existent")
	isDir = IsFile(nonExistentPath)
	assert.False(t, isDir)
}

func TestIsLink(t *testing.T) {
	// 跳过 Windows 上的符号链接测试
	if runtime.GOOS == "windows" {
		t.Skip("Windows 上普通用户无法创建符号链接，需要管理员权限")
	}

	// 测试存在的链接
	linkPath := "test_link"
	err := os.Symlink("target", linkPath)
	assert.NoError(t, err)
	defer os.Remove(linkPath)

	assert.True(t, IsLink(linkPath))

	// 测试不存在的路径
	nonExistentPath := "non_existent_path"
	assert.False(t, IsLink(nonExistentPath))

	// 测试普通文件
	filePath := "test_file"
	err = os.WriteFile(filePath, []byte("test"), 0644)
	assert.NoError(t, err)
	defer os.Remove(filePath)

	assert.False(t, IsLink(filePath))
}

func TestIsReadable(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "test_dir")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// 测试目录存在且是目录
	isDir := IsReadable(tempDir)
	assert.True(t, isDir)

	// 测试文件不是目录
	tempFile := filepath.Join(tempDir, "test_file.txt")
	_, err = os.Create(tempFile)
	assert.NoError(t, err)

	isDir = IsReadable(tempFile)
	assert.True(t, isDir)

	// 测试不存在的路径
	nonExistentPath := filepath.Join(tempDir, "non_existent")
	isDir = IsReadable(nonExistentPath)
	assert.False(t, isDir)
}

func TestIsWritable(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "test_dir")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// 测试目录存在且是目录
	isDir := IsWritable(tempDir)
	assert.False(t, isDir)

	// 测试文件不是目录
	tempFile := filepath.Join(tempDir, "test_file.txt")
	_, err = os.Create(tempFile)
	assert.NoError(t, err)

	isDir = IsWritable(tempFile)
	assert.True(t, isDir)

	// 测试不存在的路径
	nonExistentPath := filepath.Join(tempDir, "non_existent")
	isDir = IsWritable(nonExistentPath)
	assert.False(t, isDir)
}
