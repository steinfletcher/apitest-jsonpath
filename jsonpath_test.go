package jsonpath_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/assert"

	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

func TestApiTest_Contains(t *testing.T) {
	handler := http.NewServeMux()
	handler.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"a": 12345, "b": [{"key": "c", "value": "result"}]}`))
		if err != nil {
			panic(err)
		}
	})

	apitest.New().
		Handler(handler).
		Get("/hello").
		Expect(t).
		Assert(jsonpath.Contains(`$.b[? @.key=="c"].value`, "result")).
		End()
}

func TestApiTest_Equal_Numeric(t *testing.T) {
	handler := http.NewServeMux()
	handler.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"a": 12345, "b": [{"key": "c", "value": "result"}]}`))
		if err != nil {
			panic(err)
		}
	})

	apitest.New().
		Handler(handler).
		Get("/hello").
		Expect(t).
		Assert(jsonpath.Equal(`$.a`, float64(12345))).
		End()
}

func TestApiTest_Equal_String(t *testing.T) {
	handler := http.NewServeMux()
	handler.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"a": "12345", "b": [{"key": "c", "value": "result"}]}`))
		if err != nil {
			panic(err)
		}
	})

	apitest.New().
		Handler(handler).
		Get("/hello").
		Expect(t).
		Assert(jsonpath.Equal(`$.a`, "12345")).
		End()
}

func TestApiTest_Equal_Map(t *testing.T) {
	handler := http.NewServeMux()
	handler.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"a": "hello", "b": 12345}`))
		if err != nil {
			panic(err)
		}
	})

	apitest.New().
		Handler(handler).
		Get("/hello").
		Expect(t).
		Assert(jsonpath.Equal(`$`, map[string]interface{}{"a": "hello", "b": float64(12345)})).
		End()
}

func TestApiTest_Len(t *testing.T) {
	handler := http.NewServeMux()
	handler.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"a": [1, 2, 3], "b": "c"}`))
		if err != nil {
			panic(err)
		}
	})

	apitest.New().
		Handler(handler).
		Get("/hello").
		Expect(t).
		Assert(jsonpath.Len(`$.a`, 3)).
		Assert(jsonpath.Len(`$.b`, 1)).
		End()
}

func TestApiTest_GreaterThan(t *testing.T) {
	handler := http.NewServeMux()
	handler.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"a": [1, 2, 3], "b": "c"}`))
		if err != nil {
			panic(err)
		}
	})

	apitest.New().
		Handler(handler).
		Get("/hello").
		Expect(t).
		Assert(jsonpath.GreaterThan(`$.a`, 2)).
		Assert(jsonpath.GreaterThan(`$.b`, 0)).
		End()
}

func TestApiTest_LessThan(t *testing.T) {
	handler := http.NewServeMux()
	handler.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"a": [1, 2, 3], "b": "c"}`))
		if err != nil {
			panic(err)
		}
	})

	apitest.New().
		Handler(handler).
		Get("/hello").
		Expect(t).
		Assert(jsonpath.LessThan(`$.a`, 4)).
		Assert(jsonpath.LessThan(`$.b`, 2)).
		End()
}

func TestApiTest_Present(t *testing.T) {
	handler := http.NewServeMux()
	handler.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"a": 22}`))
		if err != nil {
			panic(err)
		}
	})

	apitest.New().
		Handler(handler).
		Get("/hello").
		Expect(t).
		Assert(jsonpath.Present(`$.a`)).
		Assert(jsonpath.NotPresent(`$.password`)).
		End()
}

func TestApiTest_Matches(t *testing.T) {
	testCases := [][]string{
		{`$.aString`, `^[mot]{3}<3[AB][re]{3}$`},
		{`$.aNumber`, `^\d$`},
		{`$.anObject.aNumber`, `^\d\.\d{3}$`},
		{`$.aNumberSlice[1]`, `^[80]$`},
		{`$.anObject.aBool`, `^true$`},
	}

	handler := http.NewServeMux()
	handler.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"anObject":{"aString":"tom<3Beer","aNumber":7.212,"aBool":true},"aString":"tom<3Beer","aNumber":7,"aNumberSlice":[7,8,9],"aStringSlice":["7","8","9"]}`))
		if err != nil {
			panic(err)
		}
	})

	for testNumber, testCase := range testCases {
		t.Run(fmt.Sprintf("match test %d", testNumber), func(t *testing.T) {
			apitest.New().
				Handler(handler).
				Get("/hello").
				Expect(t).
				Assert(jsonpath.Matches(testCase[0], testCase[1])).
				End()
		})
	}
}

func TestApiTest_Matches_FailCompile(t *testing.T) {
	willFailToCompile := jsonpath.Matches(`$.b[? @.key=="c"].value`, `\`)
	err := willFailToCompile(nil, nil)

	assert.EqualError(t, err, `invalid pattern: '\'`)
}

func TestApiTest_Matches_FailForObject(t *testing.T) {
	matcher := jsonpath.Matches(`$.anObject`, `.+`)

	err := matcher(&http.Response{
		Body: ioutil.NopCloser(bytes.NewBuffer([]byte(`{"anObject":{"aString":"lol"}}`))),
	}, nil)

	assert.EqualError(t, err, "unable to match using type: map")
}

func TestApiTest_Matches_FailForArray(t *testing.T) {
	matcher := jsonpath.Matches(`$.aSlice`, `.+`)

	err := matcher(&http.Response{
		Body: ioutil.NopCloser(bytes.NewBuffer([]byte(`{"aSlice":[1,2,3]}`))),
	}, nil)

	assert.EqualError(t, err, "unable to match using type: slice")
}

func TestApiTest_Matches_FailForNilValue(t *testing.T) {
	matcher := jsonpath.Matches(`$.nothingHere`, `.+`)

	err := matcher(&http.Response{
		Body: ioutil.NopCloser(bytes.NewBuffer([]byte(`{"aSlice":[1,2,3]}`))),
	}, nil)

	assert.EqualError(t, err, "no match for pattern: '$.nothingHere'")
}
