[![Build Status](https://travis-ci.org/steinfletcher/apitest-jsonpath.svg?branch=master)](https://travis-ci.org/steinfletcher/apitest-jsonpath)

# apitest-jsonpath

This library provides jsonpath assertions for [apitest](https://github.com/steinfletcher/apitest).

# Installation

```bash
go get -u github.com/steinfletcher/apitest-jsonpath
```

## Examples

### Equals

`Equals` checks for value equality when the json path expression returns a single result. Given the response is `{"id": 12345}`

```go
apitest.New(handler).
	Get("/hello").
	Expect(t).
	Assert(jsonpath.Equal(`$.id`, float64(12345))).
	End()
```

We can also provide more complex expected values. Given the response `{"message": "hello", "id": 12345}`.

```go
apitest.New().
	Handler(handler).
	Get("/hello").
	Expect(t).
	Assert(jsonpath.Equal(`$`, map[string]interface{}{"message": "hello", "id": float64(12345)})).
	End()
```

### Contains

When the jsonpath expression returns an array, use `Contains` to assert that the expected value is contained in the result. Given the response is `{"a": 12345, "b": [{"key": "c", "value": "result"}]}`, we can assert on the result like so

```go
apitest.New().
	Handler(handler).
	Get("/hello").
	Expect(t).
	Assert(jsonpath.Contains(`$.b[? @.key=="c"].value`, "result")).
	End()
```

### Present / NotPresent

Use `Present` and `NotPresent` to check the presence of a field in the response without evaluating its value.

```go
apitest.New().
	Handler(handler).
	Get("/hello").
	Expect(t).
	Assert(jsonpath.Present(`$.a`)).
	Assert(jsonpath.NotPresent(`$.password`)).
	End()
```

### Matches

Use `Matches` to check that a single path element of type string, number or bool matches a regular expression.

```go
apitest.New().
	Handler(handler).
	Get("/hello").
	Expect(t).
	Assert(jsonpath.Matches(`$.a`, `^[abc]{1,3}$`)).
	End()
```

### Len

Use `Len` to check to the length of the returned value. Given the response is `{"items": [1, 2, 3]}`, we can assert on the length of items like so

```go
apitest.New().
	Handler(handler).
	Get("/articles?category=golang").
	Expect(t).
	Assert(jsonpath.Len(`$.items`, 3)).
	End()
```

### GreaterThan

Use `GreaterThan` to enforce a minimum length on the returned value.

```go
apitest.New().
	Handler(handler).
	Get("/articles?category=golang").
	Expect(t).
	Assert(jsonpath.GreaterThan(`$.items`, 2)).
	End()
```

### LessThan

Use `LessThan` to enforce a maximum length on the returned value.

```go
apitest.New().
	Handler(handler).
	Get("/articles?category=golang").
	Expect(t).
	Assert(jsonpath.LessThan(`$.items`, 4)).
	End()
```

### JWT matchers

`JWTHeaderEqual` and `JWTPayloadEqual` can be used to assert on the contents of the JWT in the response (it does not verify a JWT).

```go
func TestX(t *testing.T) {
	apitest.New().
		HandlerFunc(myHandler).
		Post("/login").
		Expect(t).
		Assert(jsonpath.JWTPayloadEqual(fromAuthHeader, `$.sub`, "1234567890")).
		Assert(jsonpath.JWTHeaderEqual(fromAuthHeader, `$.alg`, "HS256")).
		End()
}

func fromAuthHeader(res *http.Response) (string, error) {
	return res.Header.Get("Authorization"), nil
}
```
