package generator

import (
	"bytes"
	"context"
	"encoding/csv"
	"strings"
	"testing"
)

func TestGenerateCSV(t *testing.T) {
	ctx := context.Background()

	t.Run("empty rows produces header only", func(t *testing.T) {
		table := makeTable("users", "id", "name")
		got, err := generateCSV(ctx, table, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		records := parseCSV(t, got)
		if len(records) != 1 {
			t.Fatalf("expected 1 row (header), got %d", len(records))
		}
		if records[0][0] != "id" || records[0][1] != "name" {
			t.Errorf("unexpected header: %v", records[0])
		}
	})

	t.Run("single row produces header + 1 data row", func(t *testing.T) {
		table := makeTable("users", "id", "name")
		rows := []map[string]interface{}{{"id": 1, "name": "Alice"}}
		got, err := generateCSV(ctx, table, rows)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		records := parseCSV(t, got)
		if len(records) != 2 {
			t.Fatalf("expected 2 rows, got %d", len(records))
		}
	})

	t.Run("multiple rows all present", func(t *testing.T) {
		table := makeTable("users", "id", "name")
		rows := []map[string]interface{}{
			{"id": 1, "name": "Alice"},
			{"id": 2, "name": "Bob"},
			{"id": 3, "name": "Carol"},
		}
		got, err := generateCSV(ctx, table, rows)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		records := parseCSV(t, got)
		if len(records) != 4 { // header + 3 rows
			t.Fatalf("expected 4 rows, got %d", len(records))
		}
	})

	t.Run("header uses column names as-is", func(t *testing.T) {
		table := makeTable("t", "UserID", "FirstName", "created_at")
		got, err := generateCSV(ctx, table, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		records := parseCSV(t, got)
		header := records[0]
		expected := []string{"UserID", "FirstName", "created_at"}
		for i, h := range expected {
			if header[i] != h {
				t.Errorf("header[%d] = %q, want %q", i, header[i], h)
			}
		}
	})

	t.Run("column order follows table.Columns order", func(t *testing.T) {
		table := makeTable("t", "z", "a", "m")
		rows := []map[string]interface{}{{"z": "1", "a": "2", "m": "3"}}
		got, err := generateCSV(ctx, table, rows)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		records := parseCSV(t, got)
		if records[0][0] != "z" || records[0][1] != "a" || records[0][2] != "m" {
			t.Errorf("unexpected header order: %v", records[0])
		}
		if records[1][0] != "1" || records[1][1] != "2" || records[1][2] != "3" {
			t.Errorf("unexpected data order: %v", records[1])
		}
	})

	t.Run("string values written correctly", func(t *testing.T) {
		table := makeTable("t", "name")
		rows := []map[string]interface{}{{"name": "Vijay"}}
		got, err := generateCSV(ctx, table, rows)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		records := parseCSV(t, got)
		if records[1][0] != "Vijay" {
			t.Errorf("got %q, want 'Vijay'", records[1][0])
		}
	})

	t.Run("integer values converted to string", func(t *testing.T) {
		table := makeTable("t", "count")
		rows := []map[string]interface{}{{"count": 42}}
		got, err := generateCSV(ctx, table, rows)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		records := parseCSV(t, got)
		if records[1][0] != "42" {
			t.Errorf("got %q, want '42'", records[1][0])
		}
	})

	t.Run("float values converted to string", func(t *testing.T) {
		table := makeTable("t", "price")
		rows := []map[string]interface{}{{"price": 9.99}}
		got, err := generateCSV(ctx, table, rows)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		records := parseCSV(t, got)
		if records[1][0] != "9.99" {
			t.Errorf("got %q, want '9.99'", records[1][0])
		}
	})

	t.Run("bool values converted to string", func(t *testing.T) {
		table := makeTable("t", "active")
		rows := []map[string]interface{}{{"active": true}}
		got, err := generateCSV(ctx, table, rows)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		records := parseCSV(t, got)
		if records[1][0] != "true" {
			t.Errorf("got %q, want 'true'", records[1][0])
		}
	})

	t.Run("nil value written as empty string (raw bytes)", func(t *testing.T) {
		table := makeTable("t", "name")
		rows := []map[string]interface{}{{"name": nil}}
		got, err := generateCSV(ctx, table, rows)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// csv.Writer writes an empty-string field as a bare newline; verify the
		// output has two newlines: one for the header row and one for the empty data row.
		if strings.Count(string(got), "\n") < 2 {
			t.Errorf("expected at least 2 newlines for header+empty row, got: %q", string(got))
		}
	})

	t.Run("missing key written as empty string", func(t *testing.T) {
		table := makeTable("t", "id", "name")
		rows := []map[string]interface{}{{"id": 1}} // name missing
		got, err := generateCSV(ctx, table, rows)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// The data row should contain "1," — id value followed by comma and empty name field.
		if !strings.Contains(string(got), "1,") {
			t.Errorf("expected '1,' in output for missing name field, got: %q", string(got))
		}
	})

	t.Run("value with comma is quoted by csv encoder", func(t *testing.T) {
		table := makeTable("t", "addr")
		rows := []map[string]interface{}{{"addr": "123 Main St, Springfield"}}
		got, err := generateCSV(ctx, table, rows)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// parse back — csv.Reader handles quoting transparently
		records := parseCSV(t, got)
		if records[1][0] != "123 Main St, Springfield" {
			t.Errorf("got %q, want '123 Main St, Springfield'", records[1][0])
		}
	})

	t.Run("value with newline is quoted by csv encoder", func(t *testing.T) {
		table := makeTable("t", "note")
		rows := []map[string]interface{}{{"note": "line1\nline2"}}
		got, err := generateCSV(ctx, table, rows)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		records := parseCSV(t, got)
		if records[1][0] != "line1\nline2" {
			t.Errorf("got %q, want 'line1\\nline2'", records[1][0])
		}
	})

	t.Run("value with double-quote is escaped by csv encoder", func(t *testing.T) {
		table := makeTable("t", "label")
		rows := []map[string]interface{}{{"label": `say "hi"`}}
		got, err := generateCSV(ctx, table, rows)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		records := parseCSV(t, got)
		if records[1][0] != `say "hi"` {
			t.Errorf("got %q, want 'say \"hi\"'", records[1][0])
		}
	})

	t.Run("output is valid UTF-8 CSV", func(t *testing.T) {
		table := makeTable("users", "id", "name", "score")
		rows := make([]map[string]interface{}, 50)
		for i := range rows {
			rows[i] = map[string]interface{}{"id": i, "name": "user", "score": float64(i) * 1.1}
		}
		got, err := generateCSV(ctx, table, rows)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		records := parseCSV(t, got)
		if len(records) != 51 { // header + 50 rows
			t.Errorf("expected 51 rows, got %d", len(records))
		}
		for _, r := range records {
			if len(r) != 3 {
				t.Errorf("expected 3 columns per row, got %d: %v", len(r), r)
			}
		}
	})

	t.Run("each row has same number of columns as header", func(t *testing.T) {
		table := makeTable("t", "a", "b", "c")
		rows := []map[string]interface{}{
			{"a": 1, "b": "x", "c": true},
			{"a": 2, "b": "y", "c": false},
		}
		got, err := generateCSV(ctx, table, rows)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		records := parseCSV(t, got)
		for i, r := range records {
			if len(r) != 3 {
				t.Errorf("row[%d] has %d columns, want 3", i, len(r))
			}
		}
	})

	t.Run("output ends with newline", func(t *testing.T) {
		table := makeTable("t", "id")
		rows := []map[string]interface{}{{"id": 1}}
		got, err := generateCSV(ctx, table, rows)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.HasSuffix(string(got), "\n") {
			t.Errorf("expected output to end with newline")
		}
	})
}

// parseCSV is a test helper that parses CSV bytes and fails the test on error.
func parseCSV(t *testing.T, data []byte) [][]string {
	t.Helper()
	r := csv.NewReader(bytes.NewReader(data))
	records, err := r.ReadAll()
	if err != nil {
		t.Fatalf("failed to parse CSV output: %v\noutput: %s", err, data)
	}
	return records
}
