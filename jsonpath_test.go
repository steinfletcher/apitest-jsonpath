package jsonpath

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
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
		Assert(Contains(`$.b[? @.key=="c"].value`, "result")).
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
		Assert(Equal(`$.a`, float64(12345))).
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
		Assert(Equal(`$.a`, "12345")).
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
		Assert(Equal(`$`, map[string]interface{}{"a": "hello", "b": float64(12345)})).
		End()
}

func Test_IncludesElement(t *testing.T) {
	list1 := []string{"Foo", "Bar"}
	list2 := []int{1, 2}
	simpleMap := map[interface{}]interface{}{"Foo": "Bar"}

	ok, found := includesElement("Hello World", "World")
	assertTrue(t, ok)
	assertTrue(t, found)

	ok, found = includesElement(list1, "Foo")
	assertTrue(t, ok)
	assertTrue(t, found)

	ok, found = includesElement(list1, "Bar")
	assertTrue(t, ok)
	assertTrue(t, found)

	ok, found = includesElement(list2, 1)
	assertTrue(t, ok)
	assertTrue(t, found)

	ok, found = includesElement(list2, 2)
	assertTrue(t, ok)
	assertTrue(t, found)

	ok, found = includesElement(list1, "Foo!")
	assertTrue(t, ok)
	assertFalse(t, found)

	ok, found = includesElement(list2, 3)
	assertTrue(t, ok)
	assertFalse(t, found)

	ok, found = includesElement(list2, "1")
	assertTrue(t, ok)
	assertFalse(t, found)

	ok, found = includesElement(simpleMap, "Foo")
	assertTrue(t, ok)
	assertTrue(t, found)

	ok, found = includesElement(simpleMap, "Bar")
	assertTrue(t, ok)
	assertFalse(t, found)

	ok, found = includesElement(1433, "1")
	assertFalse(t, ok)
	assertFalse(t, found)
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
		Assert(Len(`$.a`, 3)).
		Assert(Len(`$.b`, 1)).
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
		Assert(Present(`$.a`)).
		Assert(NotPresent(`$.password`)).
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
				Assert(Matches(testCase[0], testCase[1])).
				End()
		})
	}
}

func TestApiTest_Matches_FailCompile(t *testing.T) {
	willFailToCompile := Matches(`$.b[? @.key=="c"].value`, `\`)
	err := willFailToCompile(nil, nil)

	assert.EqualError(t, err, `invalid pattern: '\'`)
}

func TestApiTest_Matches_FailForObject(t *testing.T) {
	matcher := Matches(`$.anObject`, `.+`)

	err := matcher(&http.Response{
		Body: ioutil.NopCloser(bytes.NewBuffer([]byte(`{"anObject":{"aString":"lol"}}`))),
	}, nil)

	assert.EqualError(t, err, "unable to match using type: map")
}

func TestApiTest_Matches_FailForArray(t *testing.T) {
	matcher := Matches(`$.aSlice`, `.+`)

	err := matcher(&http.Response{
		Body: ioutil.NopCloser(bytes.NewBuffer([]byte(`{"aSlice":[1,2,3]}`))),
	}, nil)

	assert.EqualError(t, err, "unable to match using type: slice")
}

func TestApiTest_Matches_FailForNilValue(t *testing.T) {
	matcher := Matches(`$.nothingHere`, `.+`)

	err := matcher(&http.Response{
		Body: ioutil.NopCloser(bytes.NewBuffer([]byte(`{"aSlice":[1,2,3]}`))),
	}, nil)

	assert.EqualError(t, err, "no match for pattern: '$.nothingHere'")
}

func assertTrue(t *testing.T, v bool) {
	if !v {
		t.Error("\nexpected to be true but was false")
	}
}

func assertFalse(t *testing.T, v bool) {
	if v {
		t.Error("\nexpected to be false but was true")
	}
}
