package utils

import (
	"strings"
	"testing"
)

func TestGetExprTokens(t *testing.T) {
	tests := []struct {
		name       string
		expr       string
		wantFn     string
		wantParams map[string]string
		wantErr    bool
		errContain string
	}{
		// --- happy path ---
		{
			name:       "no params",
			expr:       "rand()",
			wantFn:     "rand",
			wantParams: map[string]string{},
		},
		{
			name:       "single key=value param",
			expr:       "range(min=1)",
			wantFn:     "range",
			wantParams: map[string]string{"min": "1"},
		},
		{
			name:       "multiple key=value params",
			expr:       "range(min=1,max=100)",
			wantFn:     "range",
			wantParams: map[string]string{"min": "1", "max": "100"},
		},
		{
			name:       "params with spaces trimmed",
			expr:       "range( min = 1 , max = 100 )",
			wantFn:     "range",
			wantParams: map[string]string{"min": "1", "max": "100"},
		},
		{
			name:       "leading and trailing whitespace on expr",
			expr:       "  rand()  ",
			wantFn:     "rand",
			wantParams: map[string]string{},
		},
		{
			name:       "function name with underscores",
			expr:       "first_name(case=upper)",
			wantFn:     "first_name",
			wantParams: map[string]string{"case": "upper"},
		},
		{
			name:       "string param value",
			expr:       "format(pattern=yyyy-MM-dd)",
			wantFn:     "format",
			wantParams: map[string]string{"pattern": "yyyy-MM-dd"},
		},
		{
			name:       "many params",
			expr:       "fn(a=1,b=2,c=3,d=4)",
			wantFn:     "fn",
			wantParams: map[string]string{"a": "1", "b": "2", "c": "3", "d": "4"},
		},
		{
			name:   "param without value is stored as key-only (empty value ignored)",
			expr:   "fn(key)",
			wantFn: "fn",
			// items split on "=" has len==1, so no entry is added for bare keys
			wantParams: map[string]string{},
		},
		{
			name:   "first '(' and first ')' used as boundaries — nested parens split on them",
			expr:   "outer(inner(x=1))",
			wantFn: "outer",
			// fs=5 '(' at index 5, fe=14 ')' at index 14
			// paramsExpr = "inner(x=1", split on "=" → key="inner(x", val="1"
			wantParams: map[string]string{"inner(x": "1"},
		},

		// --- error cases ---
		{
			name:       "missing opening paren",
			expr:       "noparen",
			wantErr:    true,
			errContain: "function expression required",
		},
		{
			name:       "missing closing paren",
			expr:       "fn(a=1",
			wantErr:    true,
			errContain: "missing ')'",
		},
		{
			name:       "empty string",
			expr:       "",
			wantErr:    true,
			errContain: "function expression required",
		},
		{
			name:       "whitespace only",
			expr:       "   ",
			wantErr:    true,
			errContain: "function expression required",
		},
		{
			name:       "closing paren only",
			expr:       "fn)",
			wantErr:    true,
			errContain: "function expression required",
		},
		{
			name:       "no name",
			expr:       "()",
			wantErr:    true,
			errContain: "function name required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fnName, params, err := getExprTokens(tt.expr)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.errContain)
				}
				if tt.errContain != "" && !strings.Contains(err.Error(), tt.errContain) {
					t.Errorf("error = %q, want it to contain %q", err.Error(), tt.errContain)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if fnName != tt.wantFn {
				t.Errorf("fnName = %q, want %q", fnName, tt.wantFn)
			}
			if len(params) != len(tt.wantParams) {
				t.Errorf("params len = %d, want %d: got %v", len(params), len(tt.wantParams), params)
			}
			for k, want := range tt.wantParams {
				if got, ok := params[k]; !ok {
					t.Errorf("param %q missing from result", k)
				} else if got != want {
					t.Errorf("param[%q] = %q, want %q", k, got, want)
				}
			}
		})
	}
}
