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
