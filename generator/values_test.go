package generator

import (
	"strings"
	"testing"

	"github.com/Pallinder/go-randomdata"
	"github.com/google/uuid"
	"kumarvv.com/mockdata/constants/functiontypes"
	"kumarvv.com/mockdata/models"
	"kumarvv.com/mockdata/utils"
)

// ---- helpers ----

func col(fnName string) *models.Column {
	return &models.Column{FnName: fnName}
}

func colWithValue(fnName string, value interface{}) *models.Column {
	return &models.Column{FnName: fnName, Value: value}
}

func colWithCase(fnName, c string) *models.Column {
	return &models.Column{FnName: fnName, Case: utils.StrPtr(c)}
}

func colWithLen(fnName string, l int) *models.Column {
	return &models.Column{FnName: fnName, Len: utils.IntPtr(l)}
}

func colWithMinMax(fnName string, min, max int) *models.Column {
	return &models.Column{FnName: fnName, Min: utils.IntPtr(min), Max: utils.IntPtr(max)}
}

func simpleTable() *models.ConfigTable {
	return &models.ConfigTable{Name: "t", RowCount: 10}
}

// ---- withCase ----

func TestWithCase(t *testing.T) {
	tests := []struct {
		name   string
		column *models.Column
		value  string
		want   string
	}{
		{"nil case returns original", col("string"), "Hello World", "Hello World"},
		{"lower case", colWithCase("string", "lower"), "Hello World", "hello world"},
		{"upper case", colWithCase("string", "upper"), "Hello World", "HELLO WORLD"},
		{"unknown case returns original", colWithCase("string", "title"), "Hello World", "Hello World"},
		{"empty string with upper", colWithCase("string", "upper"), "", ""},
		{"empty string with lower", colWithCase("string", "lower"), "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := withCase(tt.column, tt.value)
			if got != tt.want {
				t.Errorf("withCase() = %q, want %q", got, tt.want)
			}
		})
	}
}

// ---- withLen ----

func TestWithLen(t *testing.T) {
	t.Run("nil len/min/max returns value unchanged", func(t *testing.T) {
		got := withLen(col("random_string"), "hello")
		if got != "hello" {
			t.Errorf("got %q, want %q", got, "hello")
		}
	})

	t.Run("len truncates longer string", func(t *testing.T) {
		c := colWithLen("random_string", 3)
		got := withLen(c, "hello")
		if len(got) != 3 {
			t.Errorf("len = %d, want 3", len(got))
		}
	})

	t.Run("len pads shorter string to exact length", func(t *testing.T) {
		c := colWithLen("random_string", 50)
		got := withLen(c, "hi")
		if len(got) != 50 {
			t.Errorf("len = %d, want 50", len(got))
		}
	})

	t.Run("len with exact match returns same length", func(t *testing.T) {
		c := colWithLen("random_string", 5)
		got := withLen(c, "hello")
		if len(got) != 5 {
			t.Errorf("len = %d, want 5", len(got))
		}
	})

	t.Run("min pads string shorter than min", func(t *testing.T) {
		c := &models.Column{FnName: "random_string", Min: utils.IntPtr(20)}
		got := withLen(c, "hi")
		if len(got) < 20 {
			t.Errorf("len = %d, want >= 20", len(got))
		}
	})

	t.Run("min does not change string already at min", func(t *testing.T) {
		c := &models.Column{FnName: "random_string", Min: utils.IntPtr(5)}
		got := withLen(c, "hello")
		if len(got) < 5 {
			t.Errorf("len = %d, want >= 5", len(got))
		}
	})

	t.Run("max truncates string longer than max", func(t *testing.T) {
		c := &models.Column{FnName: "random_string", Max: utils.IntPtr(3)}
		got := withLen(c, "hello")
		if len(got) > 3 {
			t.Errorf("len = %d, want <= 3", len(got))
		}
	})

	t.Run("max does not change string shorter than max", func(t *testing.T) {
		c := &models.Column{FnName: "random_string", Max: utils.IntPtr(20)}
		got := withLen(c, "hello")
		if got != "hello" {
			t.Errorf("got %q, want %q", got, "hello")
		}
	})

	t.Run("min and max both applied", func(t *testing.T) {
		c := colWithMinMax("random_string", 10, 15)
		got := withLen(c, "hi")
		if len(got) < 10 || len(got) > 15 {
			t.Errorf("len = %d, want 10-15", len(got))
		}
	})
}

// ---- getValue / Value() ----

func TestGetValueAndValue(t *testing.T) {
	t.Run("non-string fn returns value as-is", func(t *testing.T) {
		v := getValue(col(functiontypes.Integer), int64(42))
		got, err := v.Value()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != int64(42) {
			t.Errorf("got %v, want 42", got)
		}
	})

	t.Run("non-string fn returns bool as-is", func(t *testing.T) {
		v := getValue(col(functiontypes.RandomBoolean), true)
		got, err := v.Value()
		if err != nil || got != true {
			t.Errorf("got %v, err %v", got, err)
		}
	})

	t.Run("string fn converts value to string", func(t *testing.T) {
		v := getValue(col(functiontypes.String), "hello")
		got, err := v.Value()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "hello" {
			t.Errorf("got %v, want 'hello'", got)
		}
	})

	t.Run("string fn applies case", func(t *testing.T) {
		v := getValue(colWithCase(functiontypes.RandomFirstName, "upper"), "alice")
		got, err := v.Value()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "ALICE" {
			t.Errorf("got %v, want 'ALICE'", got)
		}
	})

	t.Run("string fn applies len", func(t *testing.T) {
		v := getValue(colWithLen(functiontypes.RandomString, 4), "hello world")
		got, err := v.Value()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if s, ok := got.(string); !ok || len(s) != 4 {
			t.Errorf("got %v (len %d), want string of len 4", got, len(s))
		}
	})
}

// ---- generateValue ----

func TestGenerateValue_FixedTypes(t *testing.T) {
	table := simpleTable()

	t.Run("string returns string value", func(t *testing.T) {
		c := colWithValue(functiontypes.String, "hello")
		got, err := generateValue(table, c, randomdata.Male, 0)
		if err != nil || got != "hello" {
			t.Errorf("got %v, err %v", got, err)
		}
	})

	t.Run("integer returns int64 value", func(t *testing.T) {
		c := colWithValue(functiontypes.Integer, "42")
		got, err := generateValue(table, c, randomdata.Male, 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != int64(42) {
			t.Errorf("got %v (%T), want int64(42)", got, got)
		}
	})

	t.Run("float returns float64 value", func(t *testing.T) {
		c := colWithValue(functiontypes.Float, "3.14")
		got, err := generateValue(table, c, randomdata.Male, 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != 3.14 {
			t.Errorf("got %v, want 3.14", got)
		}
	})

	t.Run("boolean true", func(t *testing.T) {
		c := colWithValue(functiontypes.Boolean, "true")
		got, err := generateValue(table, c, randomdata.Male, 0)
		if err != nil || got != true {
			t.Errorf("got %v, err %v", got, err)
		}
	})

	t.Run("boolean false", func(t *testing.T) {
		c := colWithValue(functiontypes.Boolean, "false")
		got, err := generateValue(table, c, randomdata.Male, 0)
		if err != nil || got != false {
			t.Errorf("got %v, err %v", got, err)
		}
	})

	t.Run("date returns time.Time from string", func(t *testing.T) {
		c := colWithValue(functiontypes.Date, "2024-06-15")
		got, err := generateValue(table, c, randomdata.Male, 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got == nil {
			t.Error("expected time.Time, got nil")
		}
	})

	t.Run("date with custom format", func(t *testing.T) {
		c := &models.Column{FnName: functiontypes.Date, Value: "15/06/2024", Format: utils.StrPtr("02/01/2006")}
		got, err := generateValue(table, c, randomdata.Male, 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got == nil {
			t.Error("expected time.Time, got nil")
		}
	})

	t.Run("datetime returns time.Time", func(t *testing.T) {
		c := colWithValue(functiontypes.DateTime, "2024-06-15")
		got, err := generateValue(table, c, randomdata.Male, 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got == nil {
			t.Error("expected time.Time, got nil")
		}
	})
}

func TestGenerateValue_Serial(t *testing.T) {
	table := simpleTable()

	t.Run("serial starts at 1 by default", func(t *testing.T) {
		c := col(functiontypes.Serial)
		got, err := generateValue(table, c, randomdata.Male, 0)
		if err != nil || got != 1 {
			t.Errorf("got %v, want 1", got)
		}
	})

	t.Run("serial increments with ix", func(t *testing.T) {
		c := col(functiontypes.Serial)
		got, err := generateValue(table, c, randomdata.Male, 4)
		if err != nil || got != 5 {
			t.Errorf("got %v, want 5", got)
		}
	})

	t.Run("serial with custom min", func(t *testing.T) {
		c := &models.Column{FnName: functiontypes.Serial, Min: utils.IntPtr(100)}
		got, err := generateValue(table, c, randomdata.Male, 0)
		if err != nil || got != 100 {
			t.Errorf("got %v, want 100", got)
		}
	})

	t.Run("serial with custom min and ix", func(t *testing.T) {
		c := &models.Column{FnName: functiontypes.Serial, Min: utils.IntPtr(10)}
		got, err := generateValue(table, c, randomdata.Male, 5)
		if err != nil || got != 15 {
			t.Errorf("got %v, want 15", got)
		}
	})
}

func TestGenerateValue_UUID(t *testing.T) {
	table := simpleTable()

	t.Run("uuid returns valid uuid string", func(t *testing.T) {
		c := col(functiontypes.UUID)
		got, err := generateValue(table, c, randomdata.Male, 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		s, ok := got.(string)
		if !ok {
			t.Fatalf("expected string, got %T", got)
		}
		if _, err := uuid.Parse(s); err != nil {
			t.Errorf("not a valid UUID: %q", s)
		}
	})

	t.Run("uuid produces unique values", func(t *testing.T) {
		c := col(functiontypes.UUID)
		seen := map[string]bool{}
		for i := 0; i < 10; i++ {
			got, _ := generateValue(table, c, randomdata.Male, i)
			s := got.(string)
			if seen[s] {
				t.Errorf("duplicate UUID: %q", s)
			}
			seen[s] = true
		}
	})
}

func TestGenerateValue_RandomString(t *testing.T) {
	table := simpleTable()

	t.Run("random_string returns non-empty string", func(t *testing.T) {
		c := col(functiontypes.RandomString)
		got, err := generateValue(table, c, randomdata.Male, 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if s, ok := got.(string); !ok || s == "" {
			t.Errorf("expected non-empty string, got %v", got)
		}
	})

	t.Run("random_string with len constraint", func(t *testing.T) {
		c := colWithLen(functiontypes.RandomString, 10)
		got, err := generateValue(table, c, randomdata.Male, 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if s, ok := got.(string); !ok || len(s) != 10 {
			t.Errorf("expected string len=10, got %q (len=%d)", s, len(s))
		}
	})

	t.Run("random_string with case=upper", func(t *testing.T) {
		c := colWithCase(functiontypes.RandomString, "upper")
		got, err := generateValue(table, c, randomdata.Male, 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		s := got.(string)
		if s != strings.ToUpper(s) {
			t.Errorf("expected uppercase, got %q", s)
		}
	})

	t.Run("random_string with case=lower", func(t *testing.T) {
		c := colWithCase(functiontypes.RandomString, "lower")
		got, err := generateValue(table, c, randomdata.Male, 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		s := got.(string)
		if s != strings.ToLower(s) {
			t.Errorf("expected lowercase, got %q", s)
		}
	})
}

func TestGenerateValue_RandomNames(t *testing.T) {
	table := simpleTable()

	nameTests := []struct {
		fnName string
		gender int
	}{
		{functiontypes.RandomFirstName, randomdata.Male},
		{functiontypes.RandomFirstName, randomdata.Female},
		{functiontypes.RandomLastName, randomdata.Male},
		{functiontypes.RandomFullName, randomdata.Male},
		{functiontypes.RandomFullName, randomdata.Female},
		{functiontypes.RandomTitle, randomdata.Male},
		{functiontypes.RandomTitle, randomdata.Female},
	}
	for _, tt := range nameTests {
		t.Run(tt.fnName, func(t *testing.T) {
			c := col(tt.fnName)
			got, err := generateValue(table, c, tt.gender, 0)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if s, ok := got.(string); !ok || s == "" {
				t.Errorf("expected non-empty string, got %v", got)
			}
		})
	}
}

func TestGenerateValue_RandomGender(t *testing.T) {
	table := simpleTable()

	t.Run("male gender", func(t *testing.T) {
		c := col(functiontypes.RandomGender)
		got, err := generateValue(table, c, randomdata.Male, 0)
		if err != nil || got != "male" {
			t.Errorf("got %v, want 'male'", got)
		}
	})

	t.Run("female gender", func(t *testing.T) {
		c := col(functiontypes.RandomGender)
		got, err := generateValue(table, c, randomdata.Female, 0)
		if err != nil || got != "female" {
			t.Errorf("got %v, want 'female'", got)
		}
	})
}

func TestGenerateValue_RandomContact(t *testing.T) {
	table := simpleTable()

	contactFns := []string{
		functiontypes.RandomEmail,
		functiontypes.RandomPhone,
	}
	for _, fnName := range contactFns {
		t.Run(fnName+" returns non-empty string", func(t *testing.T) {
			c := col(fnName)
			got, err := generateValue(table, c, randomdata.Male, 0)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if s, ok := got.(string); !ok || s == "" {
				t.Errorf("expected non-empty string, got %v", got)
			}
		})
	}
}

func TestGenerateValue_RandomLocation(t *testing.T) {
	table := simpleTable()

	locationFns := []string{
		functiontypes.RandomAddress,
		functiontypes.RandomStreet,
		functiontypes.RandomCity,
		functiontypes.RandomState,
		functiontypes.RandomState2,
		functiontypes.RandomCountry,
		functiontypes.RandomCountry2,
		functiontypes.RandomCountry3,
		functiontypes.RandomCurrency,
	}
	for _, fnName := range locationFns {
		t.Run(fnName+" returns non-empty string", func(t *testing.T) {
			c := col(fnName)
			got, err := generateValue(table, c, randomdata.Male, 0)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if s, ok := got.(string); !ok || s == "" {
				t.Errorf("expected non-empty string, got %v", got)
			}
		})
	}
}

func TestGenerateValue_RandomNumber(t *testing.T) {
	table := simpleTable()

	t.Run("random_number no constraints panics (randomdata.Number() with no args is unsupported)", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected panic from randomdata.Number() with no args, but did not panic")
			}
		}()
		c := col(functiontypes.RandomNumber)
		_, _ = generateValue(table, c, randomdata.Male, 0) //nolint
	})

	t.Run("random_number with min and max", func(t *testing.T) {
		c := colWithMinMax(functiontypes.RandomNumber, 10, 20)
		for i := 0; i < 20; i++ {
			got, err := generateValue(table, c, randomdata.Male, i)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			n := got.(int)
			if n < 10 || n > 20 {
				t.Errorf("got %d, want 10-20", n)
			}
		}
	})

	t.Run("random_number with min only", func(t *testing.T) {
		c := &models.Column{FnName: functiontypes.RandomNumber, Min: utils.IntPtr(5)}
		got, err := generateValue(table, c, randomdata.Male, 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := got.(int); !ok {
			t.Errorf("expected int, got %T", got)
		}
	})

	t.Run("random_number with max only", func(t *testing.T) {
		c := &models.Column{FnName: functiontypes.RandomNumber, Max: utils.IntPtr(100)}
		for i := 0; i < 20; i++ {
			got, err := generateValue(table, c, randomdata.Male, i)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			n := got.(int)
			if n > 100 {
				t.Errorf("got %d, want <= 100", n)
			}
		}
	})
}

func TestGenerateValue_RandomDecimal(t *testing.T) {
	table := simpleTable()

	t.Run("random_decimal no constraints panics (randomdata.Decimal() with no args is unsupported)", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected panic from randomdata.Decimal() with no args, but did not panic")
			}
		}()
		c := col(functiontypes.RandomDecimal)
		_, _ = generateValue(table, c, randomdata.Male, 0) //nolint
	})

	t.Run("random_decimal with min and max", func(t *testing.T) {
		c := colWithMinMax(functiontypes.RandomDecimal, 1, 10)
		for i := 0; i < 20; i++ {
			got, err := generateValue(table, c, randomdata.Male, i)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			f := got.(float64)
			if f < 1 || f > 10 {
				t.Errorf("got %v, want 1-10", f)
			}
		}
	})
}

func TestGenerateValue_RandomBoolean(t *testing.T) {
	table := simpleTable()

	t.Run("random_boolean returns bool", func(t *testing.T) {
		c := col(functiontypes.RandomBoolean)
		got, err := generateValue(table, c, randomdata.Male, 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := got.(bool); !ok {
			t.Errorf("expected bool, got %T", got)
		}
	})

	t.Run("random_boolean produces both true and false", func(t *testing.T) {
		c := col(functiontypes.RandomBoolean)
		seenTrue, seenFalse := false, false
		for i := 0; i < 100; i++ {
			got, _ := generateValue(table, c, randomdata.Male, i)
			if got.(bool) {
				seenTrue = true
			} else {
				seenFalse = true
			}
			if seenTrue && seenFalse {
				break
			}
		}
		if !seenTrue || !seenFalse {
			t.Error("random_boolean never produced both true and false in 100 iterations")
		}
	})
}

func TestGenerateValue_RandomDate(t *testing.T) {
	table := simpleTable()

	t.Run("random_date returns non-empty string", func(t *testing.T) {
		c := col(functiontypes.RandomDate)
		got, err := generateValue(table, c, randomdata.Male, 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if s, ok := got.(string); !ok || s == "" {
			t.Errorf("expected non-empty string, got %v", got)
		}
	})

	t.Run("random_date with custom format", func(t *testing.T) {
		c := &models.Column{FnName: functiontypes.RandomDate, Format: utils.StrPtr("2006/01/02")}
		got, err := generateValue(table, c, randomdata.Male, 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		s, ok := got.(string)
		if !ok || s == "" {
			t.Fatalf("expected non-empty string, got %v", got)
		}
		// formatted date should contain / separators
		if !strings.Contains(s, "/") {
			t.Errorf("expected / in formatted date, got %q", s)
		}
	})

	t.Run("random_day returns non-empty string", func(t *testing.T) {
		c := col(functiontypes.RandomDay)
		got, err := generateValue(table, c, randomdata.Male, 0)
		if err != nil || got == nil {
			t.Errorf("got %v, err %v", got, err)
		}
	})

	t.Run("random_month returns non-empty string", func(t *testing.T) {
		c := col(functiontypes.RandomMonth)
		got, err := generateValue(table, c, randomdata.Male, 0)
		if err != nil || got == nil {
			t.Errorf("got %v, err %v", got, err)
		}
	})

	t.Run("random_year returns int in valid range", func(t *testing.T) {
		c := col(functiontypes.RandomYear)
		for i := 0; i < 20; i++ {
			got, err := generateValue(table, c, randomdata.Male, i)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			n := got.(int)
			if n < 1900 || n > 2999 {
				t.Errorf("random_year = %d, want 1900-2999", n)
			}
		}
	})
}

func TestGenerateValue_RandomIn(t *testing.T) {
	table := simpleTable()

	t.Run("random_in_string picks from comma-separated values", func(t *testing.T) {
		c := colWithValue(functiontypes.RandomInString, "alpha,beta,gamma")
		allowed := []string{"alpha", "beta", "gamma"}
		for i := 0; i < 30; i++ {
			got, err := generateValue(table, c, randomdata.Male, i)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			s := got.(string)
			if !utils.Includes(allowed, s) {
				t.Errorf("got %q, not in allowed set %v", s, allowed)
			}
		}
	})

	t.Run("random_in_integer picks from comma-separated integers", func(t *testing.T) {
		c := colWithValue(functiontypes.RandomInInteger, "10,20,30")
		allowed := []int64{10, 20, 30}
		for i := 0; i < 30; i++ {
			got, err := generateValue(table, c, randomdata.Male, i)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			n := got.(int64)
			found := false
			for _, a := range allowed {
				if a == n {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("got %v, not in allowed set %v", n, allowed)
			}
		}
	})

	t.Run("random_in_float picks from comma-separated floats", func(t *testing.T) {
		c := colWithValue(functiontypes.RandomInFloat, "1.1,2.2,3.3")
		allowed := []float64{1.1, 2.2, 3.3}
		for i := 0; i < 30; i++ {
			got, err := generateValue(table, c, randomdata.Male, i)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			f := got.(float64)
			found := false
			for _, a := range allowed {
				if a == f {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("got %v, not in allowed set %v", f, allowed)
			}
		}
	})

	t.Run("random_in_string returns all options eventually", func(t *testing.T) {
		c := colWithValue(functiontypes.RandomInString, "x,y,z")
		seen := map[string]bool{}
		for i := 0; i < 200; i++ {
			got, _ := generateValue(table, c, randomdata.Male, i)
			seen[got.(string)] = true
		}
		for _, v := range []string{"x", "y", "z"} {
			if !seen[v] {
				t.Errorf("value %q never returned in 200 iterations", v)
			}
		}
	})
}

func TestGenerateValue_RandomRange(t *testing.T) {
	table := simpleTable()

	t.Run("random_range with min and max stays in range", func(t *testing.T) {
		c := colWithMinMax(functiontypes.RandomRange, 5, 15)
		for i := 0; i < 30; i++ {
			got, err := generateValue(table, c, randomdata.Male, i)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			n := got.(int)
			if n < 5 || n > 15 {
				t.Errorf("got %d, want 5-15", n)
			}
		}
	})

	t.Run("random_range with min only uses table.RowCount as upper bound", func(t *testing.T) {
		c := &models.Column{FnName: functiontypes.RandomRange, Min: utils.IntPtr(1)}
		got, err := generateValue(table, c, randomdata.Male, 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := got.(int); !ok {
			t.Errorf("expected int, got %T", got)
		}
	})

	t.Run("random_range with max only starts from 1", func(t *testing.T) {
		c := &models.Column{FnName: functiontypes.RandomRange, Max: utils.IntPtr(50)}
		for i := 0; i < 20; i++ {
			got, err := generateValue(table, c, randomdata.Male, i)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			n := got.(int)
			if n < 1 || n > 50 {
				t.Errorf("got %d, want 1-50", n)
			}
		}
	})

	t.Run("random_range with no constraints panics (randomdata.Number() with no args is unsupported)", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected panic from randomdata.Number() with no args, but did not panic")
			}
		}()
		c := col(functiontypes.RandomRange)
		_, _ = generateValue(table, c, randomdata.Male, 0) //nolint
	})
}

func TestGenerateValue_RandomFormat(t *testing.T) {
	table := simpleTable()

	t.Run("random_format with numPairs and separator returns formatted string", func(t *testing.T) {
		c := &models.Column{
			FnName:    functiontypes.RandomFormat,
			NumPairs:  utils.IntPtr(3),
			Separator: utils.StrPtr("-"),
		}
		got, err := generateValue(table, c, randomdata.Male, 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		s, ok := got.(string)
		if !ok || s == "" {
			t.Errorf("expected non-empty string, got %v", got)
		}
		// StringNumber(3, "-") produces pairs separated by "-"
		if strings.Count(s, "-") < 2 {
			t.Errorf("expected at least 2 dashes in %q", s)
		}
	})

	t.Run("random_paragraph returns non-empty string", func(t *testing.T) {
		c := col(functiontypes.RandomParagraph)
		got, err := generateValue(table, c, randomdata.Male, 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if s, ok := got.(string); !ok || s == "" {
			t.Errorf("expected non-empty string, got %v", got)
		}
	})
}
