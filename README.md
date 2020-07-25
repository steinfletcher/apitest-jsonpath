[![Build Status](https://travis-ci.org/steinfletcher/apitest-jsonpath.svg?branch=master)](https://travis-ci.org/steinfletcher/apitest-jsonpath)

# apitest-jsonpath

This library provides jsonpath assertions for [apitest](https://github.com/steinfletcher/apitest).

# Installation

```bash
go get -u github.com/steinfletcher/apitest-jsonpath
```

## Examples

### Equals

`Equals` checks for value equality when the json path expression returns a single result. Given the response is `{"a": 12345}`, the result can be asserted as follows

```go
apitest.New(handler).
	Get("/hello").
	Expect(t).
	Assert(jsonpath.Equal(`$.a`, float64(12345))).
	End()
```

we can also provide more complex expected values for the response of `{"a": "hello", "b": 12345}`:

```go
apitest.New().
	Handler(handler).
	Get("/hello").
	Expect(t).
	Assert(jsonpath.Equal(`$`, map[string]interface{}{"a": "hello", "b": float64(12345)})).
	End()
```

given the response is `{"a": "hello", "b": 12345}`

### Contains

When the jsonpath expression returns an array, `Contains` should be used to assert the expected value is contained in the result. <br/>
For simple array response of `[{"key": "ka", "value": "va"},{"key": "kb", "value": "vb"}]`, to assert if response contains value of "vb":
```go
apitest.New().
	Handler(handler).
	Get("/hello").
	Expect(t).
	Assert(jsonpath.Contains("please advise how to do it").
	End()
```

For response `{"a": 12345, "b": [{"key": "c", "value": "result"}]}`, we can assert on the result like so:

```go
apitest.New().
	Handler(handler).
	Get("/hello").
	Expect(t).
	Assert(jsonpath.Contains(`$.b[? @.key=="c"].value`, "result")).
	End()
```

### Present / NotPresent

Use `Present` and `NotPresent` to check the presence of a field in the response without evaluating its value (this is the difference between present and contains).

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

Use `Len` to check to the length of the returned value. Example below for given response `{please advise how response should look like}`

```go
apitest.New().
	Handler(handler).
	Get("/articles?category=golang").
	Expect(t).
	Assert(jsonpath.Len(`$.items`, 3).
	End()
```

For above simple array to assert if the number of returned keys is two the assert would look like:
```go
please advive how it should look like.
```

### GreaterThan

Use `GreaterThan` to enforce a minimum length on the returned value.

```go
apitest.New().
	Handler(handler).
	Get("/articles?category=golang").
	Expect(t).
	Assert(jsonpath.GreaterThan(`$.items`, 2).
	End()
```

### LessThan

Use `LessThan` to enforce a maximum length on the returned value.

```go
apitest.New().
	Handler(handler).
	Get("/articles?category=golang").
	Expect(t).
	Assert(jsonpath.LessThan(`$.items`, 4).
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