package utils

import (
	"testing"
	"time"
)

// ---- ToString ----

func TestToString(t *testing.T) {
	ts := time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC)
	tests := []struct {
		name  string
		input interface{}
		want  string
	}{
		{"nil", nil, ""},
		{"string", "hello", "hello"},
		{"int", 42, "42"},
		{"int64", int64(100), "100"},
		{"float64", 3.14, "3.14"},
		{"bool true", true, "true"},
		{"bool false", false, "false"},
		{"time.Time", ts, ts.Format(time.RFC3339)},
		{"empty string", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToString(tt.input); got != tt.want {
				t.Errorf("ToString(%v) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// ---- ToInt64 ----

func TestToInt(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		want  int64
	}{
		{"nil", nil, 0},
		{"string valid", "42", 42},
		{"string zero", "0", 0},
		{"string negative", "-10", -10},
		{"string invalid", "abc", 0},
		{"int", int(7), 7},
		{"int8", int8(8), 8},
		{"int16", int16(16), 16},
		{"int32", int32(32), 32},
		{"int64", int64(64), 64},
		{"uint", uint(5), 5},
		{"uint8", uint8(8), 8},
		{"uint16", uint16(16), 16},
		{"uint32", uint32(32), 32},
		{"uint64", uint64(64), 64},
		{"float64", float64(3.9), 3},
		{"float32", float32(2.7), 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToInt64(tt.input); got != tt.want {
				t.Errorf("ToInt64(%v) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

// ---- ToBool ----

func TestToBool(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		want  bool
	}{
		{"nil", nil, false},
		{"bool true", true, true},
		{"bool false", false, false},
		{"string true", "true", true},
		{"string TRUE uppercase", "TRUE", true},
		{"string True mixed", "True", true},
		{"string y", "y", true},
		{"string yes", "yes", true},
		{"string YES", "YES", true},
		{"string false", "false", false},
		{"string no", "no", false},
		{"string 1", "1", false},
		{"string empty", "", false},
		{"int 1", 1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToBool(tt.input); got != tt.want {
				t.Errorf("ToBool(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

// ---- ToFloat ----

func TestToFloat(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		want  float64
	}{
		{"nil", nil, 0},
		{"string valid", "3.14", 3.14},
		{"string integer", "10", 10.0},
		{"string invalid", "abc", 0},
		{"string empty", "", 0},
		{"float64", float64(2.71), 2.71},
		{"float32", float32(1.5), float64(float32(1.5))},
		{"int (unsupported kind)", 42, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToFloat(tt.input); got != tt.want {
				t.Errorf("ToFloat(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

// ---- ToTime / ToTimeFormat ----

func TestToTime(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		want    time.Time
		wantErr bool
	}{
		{"nil", nil, time.Time{}, false},
		{"valid date string", "2024-06-15", time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC), false},
		{"invalid date string", "not-a-date", time.Time{}, true},
		{"time.Time passthrough", time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToTime(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToTime(%v) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && !got.Equal(tt.want) {
				t.Errorf("ToTime(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestToTimeFormat(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		format  string
		want    time.Time
		wantErr bool
	}{
		{"nil", nil, time.RFC3339, time.Time{}, false},
		{"RFC3339 string", "2024-06-15T10:30:00Z", time.RFC3339, time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC), false},
		{"custom format", "15/06/2024", "02/01/2006", time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC), false},
		{"invalid string", "bad", time.RFC3339, time.Time{}, true},
		{"time.Time passthrough", time.Date(2023, 3, 10, 0, 0, 0, 0, time.UTC), DateFormatYMD, time.Date(2023, 3, 10, 0, 0, 0, 0, time.UTC), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToTimeFormat(tt.input, tt.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToTimeFormat(%v, %q) error = %v, wantErr %v", tt.input, tt.format, err, tt.wantErr)
				return
			}
			if !tt.wantErr && !got.Equal(tt.want) {
				t.Errorf("ToTimeFormat(%v, %q) = %v, want %v", tt.input, tt.format, got, tt.want)
			}
		})
	}
}
