package set_test

import (
	"qqlx/base/helpers"
	"testing"
)

func TestDeduplicate(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		expected []int
	}{
		{
			name:     "empty slice",
			slice:    []int{},
			expected: []int{},
		},
		{
			name:     "no duplicates",
			slice:    []int{1, 2, 3},
			expected: []int{1, 2, 3},
		},
		{
			name:     "some duplicates",
			slice:    []int{1, 2, 2, 3},
			expected: []int{1, 2, 3},
		},
		{
			name:     "all duplicates",
			slice:    []int{1, 1, 1},
			expected: []int{1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := helpers.Deduplicate(tt.slice)
			if len(got) != len(tt.expected) {
				t.Errorf("Deduplicate() = %v, want %v", got, tt.expected)
				return
			}
			for i := range got {
				if got[i] != tt.expected[i] {
					t.Errorf("Deduplicate() = %v, want %v", got, tt.expected)
					return
				}
			}
		})
	}
}
