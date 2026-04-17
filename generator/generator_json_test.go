package generator

import (
	"context"
	"encoding/json"
	"testing"
)

func TestGenerateJSON(t *testing.T) {
	ctx := context.Background()

	t.Run("nil rows marshals to null", func(t *testing.T) {
		got, err := generateJSON(ctx, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(got) != "null" {
			t.Errorf("got %q, want %q", string(got), "null")
		}
	})

	t.Run("empty slice marshals to empty array", func(t *testing.T) {
		got, err := generateJSON(ctx, []map[string]interface{}{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(got) != "[]" {
			t.Errorf("got %q, want %q", string(got), "[]")
		}
	})

	t.Run("single row round-trips correctly", func(t *testing.T) {
		rows := []map[string]interface{}{
			{"id": 1, "name": "Alice"},
		}
		got, err := generateJSON(ctx, rows)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		var result []map[string]interface{}
		if err := json.Unmarshal(got, &result); err != nil {
			t.Fatalf("output is not valid JSON: %v\noutput: %s", err, got)
		}
		if len(result) != 1 {
			t.Errorf("expected 1 row, got %d", len(result))
		}
	})

	t.Run("multiple rows all present", func(t *testing.T) {
		rows := []map[string]interface{}{
			{"id": 1, "name": "Alice"},
			{"id": 2, "name": "Bob"},
			{"id": 3, "name": "Carol"},
		}
		got, err := generateJSON(ctx, rows)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		var result []map[string]interface{}
		if err := json.Unmarshal(got, &result); err != nil {
			t.Fatalf("output is not valid JSON: %v", err)
		}
		if len(result) != 3 {
			t.Errorf("expected 3 rows, got %d", len(result))
		}
	})

	t.Run("string values preserved", func(t *testing.T) {
		rows := []map[string]interface{}{{"name": "Vijay"}}
		got, err := generateJSON(ctx, rows)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		var result []map[string]interface{}
		json.Unmarshal(got, &result)
		if result[0]["name"] != "Vijay" {
			t.Errorf("expected name=Vijay, got %v", result[0]["name"])
		}
	})

	t.Run("numeric values preserved", func(t *testing.T) {
		rows := []map[string]interface{}{{"score": 42.5}}
		got, err := generateJSON(ctx, rows)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		var result []map[string]interface{}
		json.Unmarshal(got, &result)
		if result[0]["score"] != 42.5 {
			t.Errorf("expected score=42.5, got %v", result[0]["score"])
		}
	})

	t.Run("boolean values preserved", func(t *testing.T) {
		rows := []map[string]interface{}{{"active": true}, {"active": false}}
		got, err := generateJSON(ctx, rows)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		var result []map[string]interface{}
		json.Unmarshal(got, &result)
		if result[0]["active"] != true || result[1]["active"] != false {
			t.Errorf("boolean values not preserved: %v", result)
		}
	})

	t.Run("nil field value marshals to null", func(t *testing.T) {
		rows := []map[string]interface{}{{"name": nil}}
		got, err := generateJSON(ctx, rows)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		var result []map[string]interface{}
		json.Unmarshal(got, &result)
		if v, ok := result[0]["name"]; !ok || v != nil {
			t.Errorf("expected name=null, got %v", v)
		}
	})

	t.Run("special characters in strings are escaped", func(t *testing.T) {
		rows := []map[string]interface{}{{"msg": `say "hello" & <goodbye>`}}
		got, err := generateJSON(ctx, rows)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		var result []map[string]interface{}
		if err := json.Unmarshal(got, &result); err != nil {
			t.Fatalf("output is not valid JSON: %v\noutput: %s", err, got)
		}
		if result[0]["msg"] != `say "hello" & <goodbye>` {
			t.Errorf("special characters not round-tripped correctly: %v", result[0]["msg"])
		}
	})

	t.Run("nested map value marshals correctly", func(t *testing.T) {
		rows := []map[string]interface{}{
			{"meta": map[string]interface{}{"key": "val"}},
		}
		got, err := generateJSON(ctx, rows)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		var result []map[string]interface{}
		if err := json.Unmarshal(got, &result); err != nil {
			t.Fatalf("output is not valid JSON: %v", err)
		}
		meta, ok := result[0]["meta"].(map[string]interface{})
		if !ok || meta["key"] != "val" {
			t.Errorf("nested map not preserved: %v", result[0]["meta"])
		}
	})

	t.Run("slice value marshals correctly", func(t *testing.T) {
		rows := []map[string]interface{}{
			{"tags": []interface{}{"go", "test", "mock"}},
		}
		got, err := generateJSON(ctx, rows)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		var result []map[string]interface{}
		if err := json.Unmarshal(got, &result); err != nil {
			t.Fatalf("output is not valid JSON: %v", err)
		}
		tags, ok := result[0]["tags"].([]interface{})
		if !ok || len(tags) != 3 {
			t.Errorf("slice value not preserved: %v", result[0]["tags"])
		}
	})

	t.Run("output is indented", func(t *testing.T) {
		rows := []map[string]interface{}{{"a": 1}}
		got, err := generateJSON(ctx, rows)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// MarshalIndent with "  " produces newlines and spaces
		output := string(got)
		if output == "" || output[0] != '[' {
			t.Errorf("unexpected output start: %q", output)
		}
		hasIndent := false
		for _, line := range []string{output} {
			if len(line) > 2 && line[1] == ' ' {
				hasIndent = true
				break
			}
		}
		_ = hasIndent // presence of newlines confirms indentation
		compacted, _ := json.Marshal(rows)
		if string(got) == string(compacted) {
			t.Errorf("expected indented output but got compact JSON")
		}
	})

	t.Run("unmarshalable value returns error", func(t *testing.T) {
		rows := []map[string]interface{}{
			{"ch": make(chan int)}, // channels cannot be marshaled
		}
		_, err := generateJSON(ctx, rows)
		if err == nil {
			t.Error("expected error for unmarshalable value, got nil")
		}
	})

	t.Run("output is valid JSON for many rows", func(t *testing.T) {
		rows := make([]map[string]interface{}, 100)
		for i := range rows {
			rows[i] = map[string]interface{}{"id": i, "val": float64(i) * 1.5}
		}
		got, err := generateJSON(ctx, rows)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		var result []map[string]interface{}
		if err := json.Unmarshal(got, &result); err != nil {
			t.Fatalf("output is not valid JSON: %v", err)
		}
		if len(result) != 100 {
			t.Errorf("expected 100 rows, got %d", len(result))
		}
	})
}
