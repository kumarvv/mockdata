package configs

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"kumarvv.com/mockdata/models"
)

// ---- helpers ----

func writeTempYAML(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "config-*.yml")
	if err != nil {
		t.Fatalf("could not create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("could not write temp file: %v", err)
	}
	_ = f.Close()
	return f.Name()
}

func errMessages(errs []error) []string {
	msgs := make([]string, len(errs))
	for i, e := range errs {
		msgs[i] = e.Error()
	}
	return msgs
}

func hasErrContaining(errs []error, substr string) bool {
	for _, e := range errs {
		if strings.Contains(e.Error(), substr) {
			return true
		}
	}
	return false
}

// ---- fillEnvVars ----

func TestFillEnvVars(t *testing.T) {
	t.Run("blank DbConnStr is a no-op", func(t *testing.T) {
		cfg := &models.Config{}
		fillEnvVars(cfg)
		if cfg.Target.DbConnStr != "" {
			t.Errorf("expected empty, got %q", cfg.Target.DbConnStr)
		}
	})

	t.Run("no %% tokens leaves string unchanged", func(t *testing.T) {
		cfg := &models.Config{Target: models.ConfigTarget{DbConnStr: "host=localhost port=5432"}}
		fillEnvVars(cfg)
		if cfg.Target.DbConnStr != "host=localhost port=5432" {
			t.Errorf("unexpected change: %q", cfg.Target.DbConnStr)
		}
	})

	t.Run("single token replaced with env var value", func(t *testing.T) {
		t.Setenv("DB_PASS", "secret")
		cfg := &models.Config{Target: models.ConfigTarget{DbConnStr: "host=localhost password=%%DB_PASS%%"}}
		fillEnvVars(cfg)
		if cfg.Target.DbConnStr != "host=localhost password=secret" {
			t.Errorf("got %q", cfg.Target.DbConnStr)
		}
	})

	t.Run("multiple tokens all replaced", func(t *testing.T) {
		t.Setenv("DB_HOST", "myhost")
		t.Setenv("DB_PASS", "mypass")
		cfg := &models.Config{Target: models.ConfigTarget{DbConnStr: "host=%%DB_HOST%% password=%%DB_PASS%%"}}
		fillEnvVars(cfg)
		if cfg.Target.DbConnStr != "host=myhost password=mypass" {
			t.Errorf("got %q", cfg.Target.DbConnStr)
		}
	})

	t.Run("token with unset env var replaced with empty string", func(t *testing.T) {
		_ = os.Unsetenv("UNSET_VAR")
		cfg := &models.Config{Target: models.ConfigTarget{DbConnStr: "host=%%UNSET_VAR%%"}}
		fillEnvVars(cfg)
		if cfg.Target.DbConnStr != "host=" {
			t.Errorf("got %q", cfg.Target.DbConnStr)
		}
	})

	t.Run("token at start of string is not processed (startx==0 guard)", func(t *testing.T) {
		// startx <= 0 breaks the loop, so %%VAR%% at position 0 is never extracted
		t.Setenv("VAR", "value")
		original := "%%VAR%%"
		cfg := &models.Config{Target: models.ConfigTarget{DbConnStr: original}}
		fillEnvVars(cfg)
		if cfg.Target.DbConnStr != original {
			t.Errorf("expected no change for token at start, got %q", cfg.Target.DbConnStr)
		}
	})

	t.Run("unclosed %% token stops extraction", func(t *testing.T) {
		cfg := &models.Config{Target: models.ConfigTarget{DbConnStr: "host=%%NOCLOSE"}}
		original := cfg.Target.DbConnStr
		fillEnvVars(cfg)
		if cfg.Target.DbConnStr != original {
			t.Errorf("expected no change, got %q", cfg.Target.DbConnStr)
		}
	})
}

// ---- validate ----

func TestValidate(t *testing.T) {
	t.Run("nil config returns error", func(t *testing.T) {
		errs := validate(nil)
		if len(errs) == 0 || !strings.Contains(errs[0].Error(), "config is nil") {
			t.Errorf("expected 'config is nil' error, got %v", errMessages(errs))
		}
	})

	t.Run("blank target type defaults to sql", func(t *testing.T) {
		cfg := &models.Config{
			Target: models.ConfigTarget{ToPath: "/tmp"},
			Tables: []*models.ConfigTable{{Name: "t", Mode: "append", RawColumns: []map[string]string{{"id": "uuid()"}}}},
		}
		errs := validate(cfg)
		if len(errs) > 0 {
			t.Errorf("unexpected errors: %v", errMessages(errs))
		}
		if cfg.Target.Type != "sql" {
			t.Errorf("expected target type defaulted to 'sql', got %q", cfg.Target.Type)
		}
	})

	t.Run("invalid target type returns error", func(t *testing.T) {
		cfg := &models.Config{
			Target: models.ConfigTarget{Type: "bogus", ToPath: "/tmp"},
			Tables: []*models.ConfigTable{{Name: "t", Mode: "append"}},
		}
		errs := validate(cfg)
		if !hasErrContaining(errs, "invalid target type") {
			t.Errorf("expected 'invalid target type' error, got %v", errMessages(errs))
		}
	})

	t.Run("target type db requires db_type", func(t *testing.T) {
		cfg := &models.Config{
			Target: models.ConfigTarget{Type: "db", DbConnStr: "conn"},
			Tables: []*models.ConfigTable{{Name: "t", Mode: "append"}},
		}
		errs := validate(cfg)
		if !hasErrContaining(errs, "db_type is required") {
			t.Errorf("expected 'db_type is required', got %v", errMessages(errs))
		}
	})

	t.Run("target type db with invalid db_type returns error", func(t *testing.T) {
		cfg := &models.Config{
			Target: models.ConfigTarget{Type: "db", DbType: "oracle", DbConnStr: "conn"},
			Tables: []*models.ConfigTable{{Name: "t", Mode: "append"}},
		}
		errs := validate(cfg)
		if !hasErrContaining(errs, "invalid db_type") {
			t.Errorf("expected 'invalid db_type', got %v", errMessages(errs))
		}
	})

	t.Run("target type db requires db_conn_str", func(t *testing.T) {
		cfg := &models.Config{
			Target: models.ConfigTarget{Type: "db", DbType: "sqlite"},
			Tables: []*models.ConfigTable{{Name: "t", Mode: "append"}},
		}
		errs := validate(cfg)
		if !hasErrContaining(errs, "db_conn_str is required") {
			t.Errorf("expected 'db_conn_str is required', got %v", errMessages(errs))
		}
	})

	t.Run("non-db target requires to_path", func(t *testing.T) {
		cfg := &models.Config{
			Target: models.ConfigTarget{Type: "json"},
			Tables: []*models.ConfigTable{{Name: "t", Mode: "append"}},
		}
		errs := validate(cfg)
		if !hasErrContaining(errs, "to_path is required") {
			t.Errorf("expected 'to_path is required', got %v", errMessages(errs))
		}
	})

	t.Run("empty tables returns error", func(t *testing.T) {
		cfg := &models.Config{
			Target: models.ConfigTarget{Type: "json", ToPath: "/tmp"},
		}
		errs := validate(cfg)
		if !hasErrContaining(errs, "at least one table") {
			t.Errorf("expected 'at least one table' error, got %v", errMessages(errs))
		}
	})

	t.Run("table with blank name returns error", func(t *testing.T) {
		cfg := &models.Config{
			Target: models.ConfigTarget{Type: "json", ToPath: "/tmp"},
			Tables: []*models.ConfigTable{{Name: "", Mode: "append"}},
		}
		errs := validate(cfg)
		if !hasErrContaining(errs, "table name is required") {
			t.Errorf("expected 'table name is required', got %v", errMessages(errs))
		}
	})

	t.Run("table with blank mode defaults to append", func(t *testing.T) {
		table := &models.ConfigTable{Name: "users", Mode: "", RawColumns: []map[string]string{{"id": "uuid()"}}}
		cfg := &models.Config{
			Target: models.ConfigTarget{Type: "json", ToPath: "/tmp"},
			Tables: []*models.ConfigTable{table},
		}
		errs := validate(cfg)
		if len(errs) > 0 {
			t.Errorf("unexpected errors: %v", errMessages(errs))
		}
	})

	t.Run("table with invalid mode returns error", func(t *testing.T) {
		cfg := &models.Config{
			Target: models.ConfigTarget{Type: "json", ToPath: "/tmp"},
			Tables: []*models.ConfigTable{{Name: "t", Mode: "truncate"}},
		}
		errs := validate(cfg)
		if !hasErrContaining(errs, "invalid table mode") {
			t.Errorf("expected 'invalid table mode', got %v", errMessages(errs))
		}
	})

	t.Run("column with invalid expression returns error", func(t *testing.T) {
		cfg := &models.Config{
			Target: models.ConfigTarget{Type: "json", ToPath: "/tmp"},
			Tables: []*models.ConfigTable{{
				Name:       "users",
				Mode:       "append",
				RawColumns: []map[string]string{{"id": "nofunc"}},
			}},
		}
		errs := validate(cfg)
		if !hasErrContaining(errs, "failed to parse table.column") {
			t.Errorf("expected column parse error, got %v", errMessages(errs))
		}
	})

	t.Run("valid db config with all required fields", func(t *testing.T) {
		cfg := &models.Config{
			Target: models.ConfigTarget{Type: "db", DbType: "sqlite", DbConnStr: "file:test.db"},
			Tables: []*models.ConfigTable{{
				Name:       "users",
				Mode:       "append",
				RawColumns: []map[string]string{{"id": "uuid()"}},
			}},
		}
		errs := validate(cfg)
		if len(errs) > 0 {
			t.Errorf("unexpected errors: %v", errMessages(errs))
		}
	})

	t.Run("multiple errors collected", func(t *testing.T) {
		cfg := &models.Config{
			Target: models.ConfigTarget{Type: "bogus"},
			Tables: []*models.ConfigTable{{Name: "", Mode: "bad"}},
		}
		errs := validate(cfg)
		if len(errs) < 2 {
			t.Errorf("expected multiple errors, got %d: %v", len(errs), errMessages(errs))
		}
	})
}

// ---- parseValueExpr ----

func TestParseValueExpr(t *testing.T) {
	t.Run("missing opening paren returns error", func(t *testing.T) {
		_, err := parseValueExpr("col", "noparen")
		if err == nil || !strings.Contains(err.Error(), "function expression required") {
			t.Errorf("expected function expression error, got %v", err)
		}
	})

	t.Run("missing closing paren returns error", func(t *testing.T) {
		_, err := parseValueExpr("col", "uuid(")
		if err == nil || !strings.Contains(err.Error(), "missing ')'") {
			t.Errorf("expected missing ')' error, got %v", err)
		}
	})

	t.Run("blank function name returns error", func(t *testing.T) {
		_, err := parseValueExpr("col", "()")
		if err == nil || !strings.Contains(err.Error(), "function name required") {
			t.Errorf("expected 'function name required', got %v", err)
		}
	})

	t.Run("invalid function name returns error", func(t *testing.T) {
		_, err := parseValueExpr("col", "notafunc()")
		if err == nil || !strings.Contains(err.Error(), "invalid function name") {
			t.Errorf("expected 'invalid function name', got %v", err)
		}
	})

	t.Run("valid fn with no params", func(t *testing.T) {
		col, err := parseValueExpr("id", "uuid()")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if col.FnName != "uuid" {
			t.Errorf("FnName = %q, want 'uuid'", col.FnName)
		}
	})

	t.Run("value fn with positional value (no =)", func(t *testing.T) {
		col, err := parseValueExpr("flag", "boolean(true)")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if col.Value != "true" {
			t.Errorf("Value = %v, want 'true'", col.Value)
		}
	})

	t.Run("fn requiring value with no value param returns error", func(t *testing.T) {
		_, err := parseValueExpr("col", "string()")
		if err == nil || !strings.Contains(err.Error(), "value param required") {
			t.Errorf("expected 'value param required', got %v", err)
		}
	})

	t.Run("invalid param key returns error", func(t *testing.T) {
		_, err := parseValueExpr("col", "random_string(badparam=10)")
		if err == nil || !strings.Contains(err.Error(), "invalid param key") {
			t.Errorf("expected 'invalid param key', got %v", err)
		}
	})

	t.Run("valid fn with valid param key", func(t *testing.T) {
		col, err := parseValueExpr("col", "random_string(len=10)")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if col.Len == nil || *col.Len != 10 {
			t.Errorf("Len = %v, want 10", col.Len)
		}
	})

	t.Run("valid fn with multiple params", func(t *testing.T) {
		col, err := parseValueExpr("col", "random_string(min=5,max=20)")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if col.Min == nil || *col.Min != 5 {
			t.Errorf("Min = %v, want 5", col.Min)
		}
		if col.Max == nil || *col.Max != 20 {
			t.Errorf("Max = %v, want 20", col.Max)
		}
	})

	t.Run("params with spaces trimmed", func(t *testing.T) {
		col, err := parseValueExpr("col", "random_string( len = 8 )")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if col.Len == nil || *col.Len != 8 {
			t.Errorf("Len = %v, want 8", col.Len)
		}
	})

	t.Run("integer fn with value", func(t *testing.T) {
		col, err := parseValueExpr("version", "integer(42)")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if col.Value != "42" {
			t.Errorf("Value = %v, want '42'", col.Value)
		}
	})

	t.Run("string fn with value", func(t *testing.T) {
		col, err := parseValueExpr("src", "string(generated)")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if col.Value != "generated" {
			t.Errorf("Value = %v, want 'generated'", col.Value)
		}
	})

	t.Run("random_date with format param", func(t *testing.T) {
		col, err := parseValueExpr("created", "random_date(format=2006-01-02)")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if col.Format == nil || *col.Format != "2006-01-02" {
			t.Errorf("Format = %v, want '2006-01-02'", col.Format)
		}
	})

	t.Run("random_string with case param", func(t *testing.T) {
		col, err := parseValueExpr("code", "random_string(case=upper)")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if col.Case == nil || *col.Case != "upper" {
			t.Errorf("Case = %v, want 'upper'", col.Case)
		}
	})
}

// ---- buildColumn ----

func TestBuildColumn(t *testing.T) {
	t.Run("value param sets column.Value", func(t *testing.T) {
		col := buildColumn("col", "string", map[string]string{"value": "hello"})
		if col.Value != "hello" {
			t.Errorf("Value = %v, want 'hello'", col.Value)
		}
	})

	t.Run("len param sets column.Len", func(t *testing.T) {
		col := buildColumn("col", "random_string", map[string]string{"len": "15"})
		if col.Len == nil || *col.Len != 15 {
			t.Errorf("Len = %v, want 15", col.Len)
		}
	})

	t.Run("min and max params set column.Min and Max", func(t *testing.T) {
		col := buildColumn("col", "random_number", map[string]string{"min": "1", "max": "99"})
		if col.Min == nil || *col.Min != 1 {
			t.Errorf("Min = %v, want 1", col.Min)
		}
		if col.Max == nil || *col.Max != 99 {
			t.Errorf("Max = %v, want 99", col.Max)
		}
	})

	t.Run("format param sets column.Format", func(t *testing.T) {
		col := buildColumn("col", "random_date", map[string]string{"format": "2006-01-02"})
		if col.Format == nil || *col.Format != "2006-01-02" {
			t.Errorf("Format = %v, want '2006-01-02'", col.Format)
		}
	})

	t.Run("case param sets column.Case", func(t *testing.T) {
		col := buildColumn("col", "random_string", map[string]string{"case": "upper"})
		if col.Case == nil || *col.Case != "upper" {
			t.Errorf("Case = %v, want 'upper'", col.Case)
		}
	})

	t.Run("numpairs param sets column.NumPairs", func(t *testing.T) {
		col := buildColumn("col", "random_format", map[string]string{"numpairs": "3"})
		if col.NumPairs == nil || *col.NumPairs != 3 {
			t.Errorf("NumPairs = %v, want 3", col.NumPairs)
		}
	})

	t.Run("separator param sets column.Separator", func(t *testing.T) {
		col := buildColumn("col", "random_format", map[string]string{"separator": "-"})
		if col.Separator == nil || *col.Separator != "-" {
			t.Errorf("Separator = %v, want '-'", col.Separator)
		}
	})

	t.Run("unknown param key is silently ignored", func(t *testing.T) {
		col := buildColumn("col", "uuid", map[string]string{"unknown": "x"})
		if col.Value != nil || col.Len != nil || col.Min != nil {
			t.Errorf("expected all fields nil for unknown param, got %+v", col)
		}
	})

	t.Run("fn name and column name set correctly", func(t *testing.T) {
		col := buildColumn("my_col", "uuid", map[string]string{})
		if col.Name != "my_col" {
			t.Errorf("Name = %q, want 'my_col'", col.Name)
		}
		if col.FnName != "uuid" {
			t.Errorf("FnName = %q, want 'uuid'", col.FnName)
		}
	})

	t.Run("empty params leaves all optional fields nil", func(t *testing.T) {
		col := buildColumn("col", "uuid", map[string]string{})
		if col.Value != nil || col.Len != nil || col.Min != nil || col.Max != nil ||
			col.Format != nil || col.Case != nil || col.NumPairs != nil || col.Separator != nil {
			t.Errorf("expected all optional fields nil, got %+v", col)
		}
	})
}

// ---- Load ----

func TestLoad(t *testing.T) {
	t.Run("nonexistent file returns error", func(t *testing.T) {
		_, errs := Load(filepath.Join(t.TempDir(), "does-not-exist.yml"))
		if len(errs) == 0 {
			t.Error("expected error for missing file, got none")
		}
	})

	t.Run("invalid YAML returns error", func(t *testing.T) {
		path := writeTempYAML(t, ":::invalid yaml:::")
		_, errs := Load(path)
		if len(errs) == 0 {
			t.Error("expected error for invalid YAML, got none")
		}
	})

	t.Run("valid YAML with validation error returns error", func(t *testing.T) {
		yaml := `
target:
  type: json
tables: []
`
		path := writeTempYAML(t, yaml)
		_, errs := Load(path)
		if len(errs) == 0 {
			t.Error("expected validation errors, got none")
		}
	})

	t.Run("valid config returns parsed config", func(t *testing.T) {
		yaml := `
target:
  type: json
  to_path: /tmp/out
tables:
  - name: users
    mode: append
    row_count: 5
    columns:
      - id: uuid()
      - name: random_full_name()
`
		path := writeTempYAML(t, yaml)
		cfg, errs := Load(path)
		if len(errs) > 0 {
			t.Fatalf("unexpected errors: %v", errMessages(errs))
		}
		if cfg == nil {
			t.Fatal("expected config, got nil")
		}
		if cfg.Target.Type != "json" {
			t.Errorf("Target.Type = %q, want 'json'", cfg.Target.Type)
		}
		if len(cfg.Tables) != 1 || cfg.Tables[0].Name != "users" {
			t.Errorf("unexpected tables: %+v", cfg.Tables)
		}
		if len(cfg.Tables[0].Columns) != 2 {
			t.Errorf("expected 2 parsed columns, got %d", len(cfg.Tables[0].Columns))
		}
	})

	t.Run("valid db config returns parsed config", func(t *testing.T) {
		yaml := `
target:
  type: db
  db_type: sqlite
  db_conn_str: file:test.db
tables:
  - name: orders
    mode: merge
    row_count: 10
    columns:
      - id: uuid()
`
		path := writeTempYAML(t, yaml)
		cfg, errs := Load(path)
		if len(errs) > 0 {
			t.Fatalf("unexpected errors: %v", errMessages(errs))
		}
		if cfg.Target.DbType != "sqlite" {
			t.Errorf("DbType = %q, want 'sqlite'", cfg.Target.DbType)
		}
	})

	t.Run("env var in db_conn_str is substituted", func(t *testing.T) {
		t.Setenv("TEST_PASS", "s3cret")
		yaml := `
target:
  type: db
  db_type: sqlite
  db_conn_str: "file:test.db?password=%%TEST_PASS%%"
tables:
  - name: t
    mode: append
    columns:
      - id: uuid()
`
		path := writeTempYAML(t, yaml)
		cfg, errs := Load(path)
		if len(errs) > 0 {
			t.Fatalf("unexpected errors: %v", errMessages(errs))
		}
		if !strings.Contains(cfg.Target.DbConnStr, "s3cret") {
			t.Errorf("env var not substituted: %q", cfg.Target.DbConnStr)
		}
	})
}
