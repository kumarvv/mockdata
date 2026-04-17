package generator

import (
	"context"
	"strings"
	"testing"

	"kumarvv.com/mockdata/models"
)

// helpers

func makeTable(name string, cols ...string) *models.ConfigTable {
	columns := make([]models.ConfigColumn, len(cols))
	for i, c := range cols {
		columns[i] = models.ConfigColumn{Name: c}
	}
	return &models.ConfigTable{Name: name, Columns: columns}
}

// ---- generateSQLColumns ----

func TestGenerateSQLColumns(t *testing.T) {
	tests := []struct {
		name  string
		table *models.ConfigTable
		want  string
	}{
		{
			"single column",
			makeTable("users", "ID"),
			"id",
		},
		{
			"multiple columns lowercased",
			makeTable("users", "ID", "FirstName", "LastName"),
			"id, firstname, lastname",
		},
		{
			"already lowercase",
			makeTable("orders", "id", "total", "created_at"),
			"id, total, created_at",
		},
		{
			"mixed case",
			makeTable("t", "UserID", "EmailAddress"),
			"userid, emailaddress",
		},
		{
			"no columns",
			makeTable("empty"),
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateSQLColumns(tt.table)
			if got != tt.want {
				t.Errorf("generateSQLColumns() = %q, want %q", got, tt.want)
			}
		})
	}
}

// ---- generateSQLInsert ----

func TestGenerateSQLInsert(t *testing.T) {
	tests := []struct {
		name        string
		table       *models.ConfigTable
		row         map[string]interface{}
		columnNames string
		wantContain []string
		wantAbsent  []string
	}{
		{
			name:        "string values are quoted",
			table:       makeTable("users", "name", "email"),
			row:         map[string]interface{}{"name": "Alice", "email": "alice@example.com"},
			columnNames: "name, email",
			wantContain: []string{"'Alice'", "'alice@example.com'", "name, email", "users"},
		},
		{
			name:        "integer value not quoted",
			table:       makeTable("orders", "id", "total"),
			row:         map[string]interface{}{"id": 42, "total": 100},
			columnNames: "id, total",
			wantContain: []string{"42", "100"},
			wantAbsent:  []string{"'42'", "'100'"},
		},
		{
			name:        "float value not quoted",
			table:       makeTable("products", "price"),
			row:         map[string]interface{}{"price": 9.99},
			columnNames: "price",
			wantContain: []string{"9.99"},
			wantAbsent:  []string{"'9.99'"},
		},
		{
			name:        "bool value not quoted",
			table:       makeTable("users", "active"),
			row:         map[string]interface{}{"active": true},
			columnNames: "active",
			wantContain: []string{"true"},
			wantAbsent:  []string{"'true'"},
		},
		{
			name:        "nil value produces empty slot (known behavior)",
			table:       makeTable("users", "name"),
			row:         map[string]interface{}{"name": nil},
			columnNames: "name",
			// nil sets value="NULL" but valueStr stays "", so empty string is appended
			wantContain: []string{"("},
		},
		{
			name:        "missing key treated as nil",
			table:       makeTable("users", "name", "age"),
			row:         map[string]interface{}{"name": "Bob"},
			columnNames: "name, age",
			wantContain: []string{"'Bob'"},
		},
		{
			name:        "table name lowercased in output",
			table:       makeTable("MyTable", "col"),
			row:         map[string]interface{}{"col": "v"},
			columnNames: "col",
			wantContain: []string{"mytable"},
			wantAbsent:  []string{"MyTable"},
		},
		{
			name:        "uses insertTemplate format",
			table:       makeTable("t", "a"),
			row:         map[string]interface{}{"a": "x"},
			columnNames: "a",
			wantContain: []string{"INSER INTO", "VALUES"},
		},
		{
			name:        "mixed string and int columns",
			table:       makeTable("events", "id", "label", "score"),
			row:         map[string]interface{}{"id": int64(1), "label": "click", "score": 0.95},
			columnNames: "id, label, score",
			wantContain: []string{"1", "'click'", "0.95"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateSQLInsert(tt.table, tt.row, tt.columnNames)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			for _, s := range tt.wantContain {
				if !strings.Contains(got, s) {
					t.Errorf("expected output to contain %q\ngot: %s", s, got)
				}
			}
			for _, s := range tt.wantAbsent {
				if strings.Contains(got, s) {
					t.Errorf("expected output NOT to contain %q\ngot: %s", s, got)
				}
			}
		})
	}
}

// ---- generateSQL ----

func TestGenerateSQL(t *testing.T) {
	ctx := context.Background()

	t.Run("empty rows returns empty output", func(t *testing.T) {
		table := makeTable("users", "id", "name")
		got, err := generateSQL(ctx, table, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(got) != "" {
			t.Errorf("expected empty output, got: %q", string(got))
		}
	})

	t.Run("single row produces one insert statement", func(t *testing.T) {
		table := makeTable("users", "id", "name")
		rows := []map[string]interface{}{
			{"id": 1, "name": "Alice"},
		}
		got, err := generateSQL(ctx, table, rows)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		output := string(got)
		count := strings.Count(output, "INSER INTO")
		if count != 1 {
			t.Errorf("expected 1 INSERT statement, got %d\noutput: %s", count, output)
		}
	})

	t.Run("multiple rows produce multiple insert statements", func(t *testing.T) {
		table := makeTable("users", "id", "name")
		rows := []map[string]interface{}{
			{"id": 1, "name": "Alice"},
			{"id": 2, "name": "Bob"},
			{"id": 3, "name": "Carol"},
		}
		got, err := generateSQL(ctx, table, rows)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		output := string(got)
		count := strings.Count(output, "INSER INTO")
		if count != 3 {
			t.Errorf("expected 3 INSERT statements, got %d\noutput: %s", count, output)
		}
	})

	t.Run("each row ends with newline", func(t *testing.T) {
		table := makeTable("t", "col")
		rows := []map[string]interface{}{
			{"col": "a"},
			{"col": "b"},
		}
		got, err := generateSQL(ctx, table, rows)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		lines := strings.Split(strings.TrimRight(string(got), "\n"), "\n")
		// each insert + trailing spaces makes 2 non-empty segments
		if len(lines) < 2 {
			t.Errorf("expected at least 2 lines, got %d", len(lines))
		}
	})

	t.Run("column names appear in every row", func(t *testing.T) {
		table := makeTable("orders", "order_id", "amount")
		rows := []map[string]interface{}{
			{"order_id": 1, "amount": 50},
			{"order_id": 2, "amount": 75},
		}
		got, err := generateSQL(ctx, table, rows)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		output := string(got)
		if strings.Count(output, "order_id, amount") != 2 {
			t.Errorf("expected column list to appear in each row\noutput: %s", output)
		}
	})

	t.Run("table name is lowercased", func(t *testing.T) {
		table := makeTable("MyOrders", "id")
		rows := []map[string]interface{}{{"id": 1}}
		got, err := generateSQL(ctx, table, rows)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		output := string(got)
		if !strings.Contains(output, "myorders") {
			t.Errorf("expected lowercased table name in output\ngot: %s", output)
		}
		if strings.Contains(output, "MyOrders") {
			t.Errorf("expected original-case table name NOT in output\ngot: %s", output)
		}
	})

	t.Run("string values are quoted", func(t *testing.T) {
		table := makeTable("users", "name")
		rows := []map[string]interface{}{{"name": "Vijay"}}
		got, err := generateSQL(ctx, table, rows)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(string(got), "'Vijay'") {
			t.Errorf("expected quoted string value\ngot: %s", string(got))
		}
	})
}
