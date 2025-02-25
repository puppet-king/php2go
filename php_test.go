package php2go

import (
	"reflect"
	"testing"
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
