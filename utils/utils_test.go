package utils

import (
	"testing"
)

// ---- IsBlank ----

func TestIsBlank(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"empty string", "", true},
		{"spaces only", "   ", true},
		{"tabs and spaces", "\t  \t", true},
		{"newline only", "\n", true},
		{"non-blank", "hello", false},
		{"spaces around text", "  hello  ", false},
		{"single char", "a", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsBlank(tt.input); got != tt.want {
				t.Errorf("IsBlank(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

// ---- Includes ----

func TestIncludes(t *testing.T) {
	tests := []struct {
		name string
		arr  []string
		str  string
		want bool
	}{
		{"found in middle", []string{"a", "b", "c"}, "b", true},
		{"found at start", []string{"a", "b", "c"}, "a", true},
		{"found at end", []string{"a", "b", "c"}, "c", true},
		{"not found", []string{"a", "b", "c"}, "d", false},
		{"empty array", []string{}, "a", false},
		{"nil array", nil, "a", false},
		{"empty string match", []string{"", "a"}, "", true},
		{"case sensitive", []string{"Hello"}, "hello", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Includes(tt.arr, tt.str); got != tt.want {
				t.Errorf("Includes(%v, %q) = %v, want %v", tt.arr, tt.str, got, tt.want)
			}
		})
	}
}

// ---- RandomOneOf ----

func TestRandomOneOf(t *testing.T) {
	t.Run("int - result is one of input values", func(t *testing.T) {
		values := []int{10, 20, 30, 40, 50}
		for i := 0; i < 50; i++ {
			got := RandomOneOf(values...)
			if !func() bool {
				for _, v := range values {
					if v == got {
						return true
					}
				}
				return false
			}() {
				t.Errorf("RandomOneOf returned %v which is not in input", got)
			}
		}
	})

	t.Run("string - result is one of input values", func(t *testing.T) {
		values := []string{"foo", "bar", "baz"}
		for i := 0; i < 50; i++ {
			got := RandomOneOf(values...)
			if !Includes(values, got) {
				t.Errorf("RandomOneOf returned %q which is not in input", got)
			}
		}
	})

	t.Run("single value always returns that value", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			if got := RandomOneOf("only"); got != "only" {
				t.Errorf("expected 'only', got %q", got)
			}
		}
	})

	t.Run("all values eventually returned", func(t *testing.T) {
		seen := map[int]bool{}
		for i := 0; i < 1000; i++ {
			seen[RandomOneOf(1, 2, 3)] = true
		}
		for _, v := range []int{1, 2, 3} {
			if !seen[v] {
				t.Errorf("value %d was never returned in 1000 iterations", v)
			}
		}
	})
}

// ---- SplitToInt ----

func TestSplitToInt(t *testing.T) {
	tests := []struct {
		name string
		s    string
		sep  string
		want []int64
	}{
		{"comma separated", "1,2,3", ",", []int64{1, 2, 3}},
		{"pipe separated", "10|20|30", "|", []int64{10, 20, 30}},
		{"single value", "42", ",", []int64{42}},
		{"negative values", "-1,-2,-3", ",", []int64{-1, -2, -3}},
		{"invalid entries become zero", "1,abc,3", ",", []int64{1, 0, 3}},
		{"empty string yields zero", "", ",", []int64{0}},
		{"spaces not trimmed", "1, 2, 3", ",", []int64{1, 0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SplitToInt(tt.s, tt.sep)
			if len(got) != len(tt.want) {
				t.Fatalf("SplitToInt(%q, %q) len = %d, want %d", tt.s, tt.sep, len(got), len(tt.want))
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("SplitToInt(%q, %q)[%d] = %d, want %d", tt.s, tt.sep, i, got[i], tt.want[i])
				}
			}
		})
	}
}

// ---- SplitToFloat ----

func TestSplitToFloat(t *testing.T) {
	tests := []struct {
		name string
		s    string
		sep  string
		want []float64
	}{
		{"comma separated", "1.1,2.2,3.3", ",", []float64{1.1, 2.2, 3.3}},
		{"pipe separated", "1.5|2.5", "|", []float64{1.5, 2.5}},
		{"integer strings", "1,2,3", ",", []float64{1.0, 2.0, 3.0}},
		{"single value", "3.14", ",", []float64{3.14}},
		{"invalid entries become zero", "1.1,abc,3.3", ",", []float64{1.1, 0, 3.3}},
		{"empty string yields zero", "", ",", []float64{0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SplitToFloat(tt.s, tt.sep)
			if len(got) != len(tt.want) {
				t.Fatalf("SplitToFloat(%q, %q) len = %d, want %d", tt.s, tt.sep, len(got), len(tt.want))
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("SplitToFloat(%q, %q)[%d] = %v, want %v", tt.s, tt.sep, i, got[i], tt.want[i])
				}
			}
		})
	}
}
