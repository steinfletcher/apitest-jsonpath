# api-test-jsonpath

This library provides jsonpath assertions for [api-test](https://github.com/steinfletcher/api-test).

# Installation

```bash
go get -u github.com/steinfletcher/api-test-jsonpath
```

## Examples

`Equals` checks for value equality when the json path expression returns a single result. Given the response is `{"a": 12345}`, the result can be asserted as follows

```go
	apitest.New(handler).
		Get("/hello").
		Expect(t).
		Assert(jsonpath.Equal(`$.a`, float64(12345))).
		End()
```

When the jsonpath expression returns an array, use `jsonpath.Contains` to assert the expected value is contained in the result. Given the response is `{"a": 12345, "b": [{"key": "c", "value": "result"}]}`, we can assert on the result like so

```go
	apitest.New().
		Handler(handler).
		Get("/hello").
		Expect(t).
		Assert(jsonpath.Contains(`$.b[? @.key=="c"].value`, "result")).
		End()
```

we can also provide more complex expected values

```go
	apitest.New().
		Handler(handler).
		Get("/hello").
		Expect(t).
		Assert(jsonpath.Equal(`$`, map[string]interface{}{"a": "hello", "b": float64(12345)})).
		End()
```

given the response is `{"a": "hello", "b": 12345}` 
