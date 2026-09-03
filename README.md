[![Test](https://github.com/steinfletcher/apitest-jsonpath/actions/workflows/ci.yml/badge.svg)](https://github.com/steinfletcher/apitest-jsonpath/actions/workflows/ci.yml)

# apitest-jsonpath

This library provides jsonpath assertions for [apitest](https://github.com/steinfletcher/apitest).

# Installation

```bash
go get github.com/steinfletcher/apitest-jsonpath
```

## Examples

### Equal

`Equal` checks for value equality when the json path expression returns a single result. Given the response is `{"id": 12345}`

```go
apitest.Handler(handler).
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
	Assert(jsonpath.Equal(`$`, map[string]any{"message": "hello", "id": float64(12345)})).
	End()
```

### NotEqual

`NotEqual` checks that the json path expression value is not equal to given value

```go
apitest.Handler(handler).
	Get("/hello").
	Expect(t).
	Assert(jsonpath.NotEqual(`$.a`, float64(56789))).
	End()
```

we can also provide more complex expected values

```go
apitest.New().
	Handler(handler).
	Get("/hello").
	Expect(t).
	Assert(jsonpath.NotEqual(`$`, map[string]any{"a": "hello", "b": float64(56789)})).
	End()
```

given the response is `{"a": "hello", "b": 12345}`

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

Use `Present` and `NotPresent` to check the presence of a field in the response without evaluating its value. A field
is present when the expression resolves to any value other than `null`, so `false`, `0` and `""` all count as present.
An empty array or object counts as not present, since that is also what a wildcard or filter expression yields when
nothing matches.

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

### Chain

`Chain` is used to provide several assertions at once. `Equal`, `NotEqual`, `Contains`, `Len`, `GreaterThan`,
`LessThan`, `Present`, `NotPresent` and `Matches` are all available on the chain.

```go
Assert(
	jsonpath.Chain().
		Equal("a", "1").
		NotEqual("b", "2").
		Present("c").
		Len("d", 3).
		End(),
).
```

### Root

`Root` is used to avoid duplicated paths in body expectations. For example, instead of writing:

```go
Assert(jsonpath.Equal("a.b.c.d", "a").
Assert(jsonpath.Equal("a.b.c.e", "b").
Assert(jsonpath.Equal("a.b.c.f", "c").
```

it is possible to define a root path like so

```go
Assert(
	jsonpath.Root("$.a.b.c").
		Equal("d", "a").
		Equal("e", "b").
		Equal("f", "c").
		End(),
).
```

A dot is inserted between the root and each sub-expression unless the sub-expression starts with a bracket, so array
elements under the root can be addressed directly:

```go
Assert(
	jsonpath.Root("$.items").
		Equal("[0].id", float64(1)).
		Equal("[1].id", float64(2)).
		End(),
).
```

### Matching mock request bodies

The `mocks` package provides the same assertions as `apitest.Matcher` functions, for matching the JSON body of a
request made to an [apitest mock](https://github.com/steinfletcher/apitest#mocking-external-http-calls).

```go
import "github.com/steinfletcher/apitest-jsonpath/mocks"

var createUser = apitest.NewMock().
	Post("/user-api").
	AddMatcher(mocks.Equal(`$.name`, "jon")).
	AddMatcher(mocks.Contains(`$.roles`, "admin")).
	RespondWith().
	Status(http.StatusCreated).
	End()
```

`mocks.Equal`, `mocks.NotEqual`, `mocks.Contains`, `mocks.Len` and `mocks.GreaterThan` are available.