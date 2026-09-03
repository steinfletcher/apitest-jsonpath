package jsonpath

import (
	"strings"
	"testing"
)

func assertError(t *testing.T, err error, expected string) {
	t.Helper()
	if expected == "" {
		if err != nil {
			t.Fatalf("expected no error, got %q", err)
		}
		return
	}
	if err == nil || err.Error() != expected {
		t.Fatalf("expected error %q, got %v", expected, err)
	}
}

func TestLength(t *testing.T) {
	tests := map[string]struct {
		body       string
		expression string
		length     int
		expected   string
	}{
		"array":       {`{"items": [1, 2, 3]}`, `$.items`, 3, ""},
		"string":      {`{"name": "jan"}`, `$.name`, 3, ""},
		"map":         {`{"user": {"a": 1, "b": 2}}`, `$.user`, 2, ""},
		"wrong":       {`{"items": [1, 2, 3]}`, `$.items`, 2, `"3" not equal to "2"`},
		"null":        {`{"items": null}`, `$.items`, 0, "value is null"},
		"number":      {`{"count": 3}`, `$.count`, 3, "value of type float64 has no length"},
		"bool":        {`{"ok": true}`, `$.ok`, 1, "value of type bool has no length"},
		"missing key": {`{"a": 1}`, `$.items`, 0, "evaluating '$.items' resulted in error: 'unknown key items'"},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assertError(t, Length(test.expression, test.length, strings.NewReader(test.body)), test.expected)
		})
	}
}

func TestGreaterThan(t *testing.T) {
	tests := map[string]struct {
		body     string
		minimum  int
		expected string
	}{
		"longer":  {`{"items": [1, 2, 3]}`, 2, ""},
		"equal":   {`{"items": [1, 2]}`, 2, ""},
		"shorter": {`{"items": [1]}`, 2, `"1" is less than "2"`},
		"null":    {`{"items": null}`, 0, "value is null"},
		"number":  {`{"items": 5}`, 0, "value of type float64 has no length"},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assertError(t, GreaterThan(`$.items`, test.minimum, strings.NewReader(test.body)), test.expected)
		})
	}
}

func TestLessThan(t *testing.T) {
	tests := map[string]struct {
		body     string
		maximum  int
		expected string
	}{
		"shorter": {`{"items": [1]}`, 2, ""},
		"equal":   {`{"items": [1, 2]}`, 2, ""},
		"longer":  {`{"items": [1, 2, 3]}`, 2, `"3" is greater than "2"`},
		"null":    {`{"items": null}`, 0, "value is null"},
		"bool":    {`{"items": false}`, 0, "value of type bool has no length"},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assertError(t, LessThan(`$.items`, test.maximum, strings.NewReader(test.body)), test.expected)
		})
	}
}

func TestContains(t *testing.T) {
	tests := map[string]struct {
		body     string
		expected any
		err      string
	}{
		"number in array":      {`{"items": [1, 2]}`, float64(2), ""},
		"string in array":      {`{"items": ["a", "b"]}`, "b", ""},
		"substring":            {`{"items": "abc"}`, "b", ""},
		"key in map":           {`{"items": {"a": 1}}`, "a", ""},
		"number not in array":  {`{"items": [1, 2]}`, float64(5), `"[1 2]" does not contain "5"`},
		"value without length": {`{"items": 5}`, float64(5), `"5" could not be applied builtin len()`},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assertError(t, Contains(`$.items`, test.expected, strings.NewReader(test.body)), test.err)
		})
	}
}

func TestEqualAndNotEqual(t *testing.T) {
	body := `{"id": 12345, "name": "jan"}`

	assertError(t, Equal(`$.id`, float64(12345), strings.NewReader(body)), "")
	assertError(t, Equal(`$.id`, float64(1), strings.NewReader(body)), `"12345" not equal to "1"`)
	assertError(t, Equal(`$.missing`, "x", strings.NewReader(body)), "evaluating '$.missing' resulted in error: 'unknown key missing'")
	assertError(t, NotEqual(`$.name`, "jon", strings.NewReader(body)), "")
	assertError(t, NotEqual(`$.name`, "jan", strings.NewReader(body)), `"$.name" value is equal to "jan"`)
}
