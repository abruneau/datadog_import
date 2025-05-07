package utilities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapIntToInt(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		fn       func(int) int
		expected []int
	}{
		{
			name:     "empty slice",
			input:    []int{},
			fn:       func(x int) int { return x * 2 },
			expected: []int{},
		},
		{
			name:     "double numbers",
			input:    []int{1, 2, 3, 4, 5},
			fn:       func(x int) int { return x * 2 },
			expected: []int{2, 4, 6, 8, 10},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Map(tt.input, tt.fn)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMapIntToString(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		fn       func(int) string
		expected []string
	}{
		{
			name:     "convert to string",
			input:    []int{1, 2, 3},
			fn:       func(x int) string { return string(rune(x + 64)) },
			expected: []string{"A", "B", "C"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Map(tt.input, tt.fn)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFilter(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		fn       func(int) bool
		expected []int
	}{
		{
			name:     "empty slice",
			input:    []int{},
			fn:       func(x int) bool { return x > 0 },
			expected: []int{},
		},
		{
			name:     "filter even numbers",
			input:    []int{1, 2, 3, 4, 5, 6},
			fn:       func(x int) bool { return x%2 == 0 },
			expected: []int{2, 4, 6},
		},
		{
			name:     "filter positive numbers",
			input:    []int{-2, -1, 0, 1, 2},
			fn:       func(x int) bool { return x > 0 },
			expected: []int{1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Filter(tt.input, tt.fn)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUnique(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{
			name:     "empty slice",
			input:    []int{},
			expected: []int{},
		},
		{
			name:     "no duplicates",
			input:    []int{1, 2, 3, 4, 5},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "with duplicates",
			input:    []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4},
			expected: []int{1, 2, 3, 4},
		},
		{
			name:     "all duplicates",
			input:    []int{1, 1, 1, 1, 1},
			expected: []int{1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Unique(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFlatten(t *testing.T) {
	tests := []struct {
		name     string
		input    [][]int
		expected []int
	}{
		{
			name:     "empty slice",
			input:    [][]int{},
			expected: []int{},
		},
		{
			name:     "empty inner slices",
			input:    [][]int{{}, {}, {}},
			expected: []int{},
		},
		{
			name:     "single inner slice",
			input:    [][]int{{1, 2, 3}},
			expected: []int{1, 2, 3},
		},
		{
			name:     "multiple inner slices",
			input:    [][]int{{1, 2}, {3, 4}, {5, 6}},
			expected: []int{1, 2, 3, 4, 5, 6},
		},
		{
			name:     "mixed length inner slices",
			input:    [][]int{{1}, {2, 3}, {4, 5, 6}},
			expected: []int{1, 2, 3, 4, 5, 6},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Flatten(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFlatMapIntToString(t *testing.T) {
	tests := []struct {
		name     string
		input    [][]int
		fn       func(int) string
		expected []string
	}{
		{
			name:     "empty slice",
			input:    [][]int{},
			fn:       func(x int) string { return string(rune(x + 64)) },
			expected: []string{},
		},
		{
			name:     "convert numbers to letters",
			input:    [][]int{{1, 2}, {3, 4}, {5, 6}},
			fn:       func(x int) string { return string(rune(x + 64)) },
			expected: []string{"A", "B", "C", "D", "E", "F"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FlatMap(tt.input, tt.fn)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFlatMapIntToInt(t *testing.T) {
	tests := []struct {
		name     string
		input    [][]int
		fn       func(int) int
		expected []int
	}{
		{
			name:     "double numbers",
			input:    [][]int{{1, 2}, {3, 4}},
			fn:       func(x int) int { return x * 2 },
			expected: []int{2, 4, 6, 8},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FlatMap(tt.input, tt.fn)
			assert.Equal(t, tt.expected, result)
		})
	}
}
