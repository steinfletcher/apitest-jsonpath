package jsonpath

import (
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
