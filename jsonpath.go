package jsonpath

import (
	"fmt"
	"net/http"
	"reflect"
	regex "regexp"

	httputil "github.com/steinfletcher/apitest-jsonpath/http"
	"github.com/steinfletcher/apitest-jsonpath/jsonpath"
)

// Contains is a convenience function to assert that a jsonpath expression extracts a value in an array
func Contains(expression string, expected any) func(*http.Response, *http.Request) error {
	return func(res *http.Response, req *http.Request) error {
		return jsonpath.Contains(expression, expected, res.Body)
	}
}

// Equal is a convenience function to assert that a jsonpath expression extracts a value
func Equal(expression string, expected any) func(*http.Response, *http.Request) error {
	return func(res *http.Response, req *http.Request) error {
		var responseErr error
		defer func() { responseErr = res.Body.Close() }()
		responseErr = jsonpath.Equal(expression, expected, res.Body)

		return responseErr
	}
}

// NotEqual is a function to check json path expression value is not equal to given value
func NotEqual(expression string, expected any) func(*http.Response, *http.Request) error {
	return func(res *http.Response, req *http.Request) error {
		var responseErr error
		defer func() { responseErr = res.Body.Close() }()
		responseErr = jsonpath.NotEqual(expression, expected, res.Body)

		return responseErr
	}
}

// Len asserts that value is the expected length, determined by reflect.Len
func Len(expression string, expectedLength int) func(*http.Response, *http.Request) error {
	return func(res *http.Response, req *http.Request) error {
		var responseErr error
		defer func() { responseErr = res.Body.Close() }()
		responseErr = jsonpath.Length(expression, expectedLength, res.Body)

		return responseErr
	}
}

// GreaterThan asserts that value is greater than the given length, determined by reflect.Len
func GreaterThan(expression string, minimumLength int) func(*http.Response, *http.Request) error {
	return func(res *http.Response, req *http.Request) error {
		var responseErr error
		defer func() { responseErr = res.Body.Close() }()
		responseErr = jsonpath.GreaterThan(expression, minimumLength, res.Body)

		return responseErr
	}
}

// LessThan asserts that value is less than the given length, determined by reflect.Len
func LessThan(expression string, maximumLength int) func(*http.Response, *http.Request) error {
	return func(res *http.Response, req *http.Request) error {
		var responseErr error
		defer func() { responseErr = res.Body.Close() }()
		responseErr = jsonpath.LessThan(expression, maximumLength, res.Body)

		return responseErr
	}
}

// Present asserts that value returned by the expression is present
func Present(expression string) func(*http.Response, *http.Request) error {
	return func(res *http.Response, req *http.Request) error {
		var responseErr error
		defer func() { responseErr = res.Body.Close() }()
		responseErr = jsonpath.Present(expression, res.Body)

		return responseErr
	}
}

// NotPresent asserts that value returned by the expression is not present
func NotPresent(expression string) func(*http.Response, *http.Request) error {
	return func(res *http.Response, req *http.Request) error {
		var responseErr error
		defer func() { responseErr = res.Body.Close() }()
		responseErr = jsonpath.NotPresent(expression, res.Body)

		return responseErr
	}
}

// Matches asserts that the value matches the given regular expression
func Matches(expression, regexp string) func(*http.Response, *http.Request) error {
	return func(res *http.Response, req *http.Request) error {
		var responseErr error
		defer func() { responseErr = res.Body.Close() }()

		pattern, err := regex.Compile(regexp)
		if err != nil {
			return fmt.Errorf("invalid pattern: '%s'", regexp)
		}
		value, _ := jsonpath.JsonPath(res.Body, expression)
		if value == nil {
			return fmt.Errorf("no match for pattern: '%s'", expression)
		}
		kind := reflect.ValueOf(value).Kind()
		switch kind {
		case reflect.Bool,
			reflect.Int,
			reflect.Int8,
			reflect.Int16,
			reflect.Int32,
			reflect.Int64,
			reflect.Uint,
			reflect.Uint8,
			reflect.Uint16,
			reflect.Uint32,
			reflect.Uint64,
			reflect.Uintptr,
			reflect.Float32,
			reflect.Float64,
			reflect.String:
			if !pattern.MatchString(fmt.Sprintf("%v", value)) {
				return fmt.Errorf("value '%v' does not match pattern '%v'", value, regexp)
			}
			return responseErr
		default:
			return fmt.Errorf("unable to match using type: %s", kind.String())
		}
	}
}

// Chain creates a new assertion chain
func Chain() *AssertionChain {
	return &AssertionChain{rootExpression: ""}
}

// Root creates a new assertion chain prefixed with the given expression
func Root(expression string) *AssertionChain {
	return &AssertionChain{rootExpression: expression + "."}
}

// AssertionChain supports chaining assertions and root expressions
type AssertionChain struct {
	rootExpression string
	assertions     []func(*http.Response, *http.Request) error
}

// Equal adds an Equal assertion to the chain
func (r *AssertionChain) Equal(expression string, expected any) *AssertionChain {
	r.assertions = append(r.assertions, Equal(r.rootExpression+expression, expected))
	return r
}

// NotEqual adds an NotEqual assertion to the chain
func (r *AssertionChain) NotEqual(expression string, expected any) *AssertionChain {
	r.assertions = append(r.assertions, NotEqual(r.rootExpression+expression, expected))
	return r
}

// Contains adds an Contains assertion to the chain
func (r *AssertionChain) Contains(expression string, expected any) *AssertionChain {
	r.assertions = append(r.assertions, Contains(r.rootExpression+expression, expected))
	return r
}

// Present adds an Present assertion to the chain
func (r *AssertionChain) Present(expression string) *AssertionChain {
	r.assertions = append(r.assertions, Present(r.rootExpression+expression))
	return r
}

// NotPresent adds an NotPresent assertion to the chain
func (r *AssertionChain) NotPresent(expression string) *AssertionChain {
	r.assertions = append(r.assertions, NotPresent(r.rootExpression+expression))
	return r
}

// Matches adds an Matches assertion to the chain
func (r *AssertionChain) Matches(expression, regexp string) *AssertionChain {
	r.assertions = append(r.assertions, Matches(r.rootExpression+expression, regexp))
	return r
}

// End returns an func(*http.Response, *http.Request) error which is a combination of the registered assertions
func (r *AssertionChain) End() func(*http.Response, *http.Request) error {
	return func(res *http.Response, req *http.Request) error {
		for _, assertion := range r.assertions {
			if err := assertion(httputil.CopyResponse(res), httputil.CopyRequest(req)); err != nil {
				return err
			}
		}
		return nil
	}
}
