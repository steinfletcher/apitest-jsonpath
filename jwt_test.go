package jsonpath_test

import (
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"

	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

const jwt = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

func TestApiTest_JWT(t *testing.T) {
	handler := http.NewServeMux()
	handler.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Authorization", jwt)
		w.WriteHeader(http.StatusOK)
	})

	apitest.New().
		Handler(handler).
		Get("/hello").
		Expect(t).
		Assert(jsonpath.JWTPayloadEqual(fromAuthHeader, `$.name`, "John Doe")).
		Assert(jsonpath.JWTPayloadEqual(fromAuthHeader, `$.sub`, "1234567890")).
		Assert(jsonpath.JWTPayloadEqual(fromAuthHeader, `$.iat`, float64(1516239022))).
		Assert(jsonpath.JWTHeaderEqual(fromAuthHeader, `$.alg`, "HS256")).
		Assert(jsonpath.JWTHeaderEqual(fromAuthHeader, `$.typ`, "JWT")).
		End()
}

func fromAuthHeader(response *http.Response) (string, error) {
	return response.Header.Get("Authorization"), nil
}

func TestApiTest_JWT_Errors(t *testing.T) {
	selector := func(token string) func(*http.Response) (string, error) {
		return func(*http.Response) (string, error) { return token, nil }
	}
	tests := map[string]struct {
		token    string
		expected string
	}{
		"wrong number of parts": {"a.b", "invalid token: token should contain header, payload and signature"},
		"payload not base64":    {"a.!!!.c", "invalid jwt: decoding error: illegal base64 data at input byte 0"},
		"payload not json":      {"a.bm90IGpzb24.c", "invalid character 'o' in literal null (expecting 'u')"},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := jsonpath.JWTPayloadEqual(selector(test.token), `$.sub`, "x")(nil, nil)
			if err == nil || err.Error() != test.expected {
				t.Fatalf("expected error %q, got %v", test.expected, err)
			}
		})
	}
}
