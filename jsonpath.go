package jsonpath

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/steinfletcher/apitest-jsonpath/jsonpath"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	regex "regexp"
)

// Contains is a convenience function to assert that a jsonpath expression extracts a value in an array
func Contains(expression string, expected interface{}) func(*http.Response, *http.Request) error {
	return func(res *http.Response, req *http.Request) error {
		return jsonpath.Contains(expression, expected, res.Body)
	}
}

// Equal is a convenience function to assert that a jsonpath expression extracts a value
func Equal(expression string, expected interface{}) func(*http.Response, *http.Request) error {
	return func(res *http.Response, req *http.Request) error {
		return jsonpath.Equal(expression, expected, res.Body)
	}
}

// NotEqual is a function to check json path expression value is not equal to given value
func NotEqual(expression string, expected interface{}) func(*http.Response, *http.Request) error {
	return func(res *http.Response, req *http.Request) error {
		return jsonpath.NotEqual(expression, expected, res.Body)
	}
}

// Len asserts that value is the expected length, determined by reflect.Len
func Len(expression string, expectedLength int) func(*http.Response, *http.Request) error {
	return func(res *http.Response, req *http.Request) error {
		return jsonpath.Length(expression, expectedLength, res.Body)
	}
}

// GreaterThan asserts that value is greater than the given length, determined by reflect.Len
func GreaterThan(expression string, minimumLength int) func(*http.Response, *http.Request) error {
	return func(res *http.Response, req *http.Request) error {
		return jsonpath.GreaterThan(expression, minimumLength, res.Body)
	}
}

// LessThan asserts that value is less than the given length, determined by reflect.Len
func LessThan(expression string, maximumLength int) func(*http.Response, *http.Request) error {
	return func(res *http.Response, req *http.Request) error {
		return jsonpath.LessThan(expression, maximumLength, res.Body)
	}
}

// Present asserts that value returned by the expression is present
func Present(expression string) func(*http.Response, *http.Request) error {
	return func(res *http.Response, req *http.Request) error {
		return jsonpath.Present(expression, res.Body)
	}
}

// NotPresent asserts that value returned by the expression is not present
func NotPresent(expression string) func(*http.Response, *http.Request) error {
	return func(res *http.Response, req *http.Request) error {
		return jsonpath.NotPresent(expression, res.Body)
	}
}

// Matches asserts that the value matches the given regular expression
func Matches(expression string, regexp string) func(*http.Response, *http.Request) error {
	return func(res *http.Response, req *http.Request) error {
		pattern, err := regex.Compile(regexp)
		if err != nil {
			return errors.New(fmt.Sprintf("invalid pattern: '%s'", regexp))
		}
		value, _ := jsonpath.JsonPath(res.Body, expression)
		if value == nil {
			return errors.New(fmt.Sprintf("no match for pattern: '%s'", expression))
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
			if !pattern.Match([]byte(fmt.Sprintf("%v", value))) {
				return errors.New(fmt.Sprintf("value '%v' does not match pattern '%v'", value, regexp))
			}
			return nil
		default:
			return errors.New(fmt.Sprintf("unable to match using type: %s", kind.String()))
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
func (r *AssertionChain) Equal(expression string, expected interface{}) *AssertionChain {
	r.assertions = append(r.assertions, Equal(r.rootExpression+expression, expected))
	return r
}

// NotEqual adds an NotEqual assertion to the chain
func (r *AssertionChain) NotEqual(expression string, expected interface{}) *AssertionChain {
	r.assertions = append(r.assertions, NotEqual(r.rootExpression+expression, expected))
	return r
}

// Contains adds an Contains assertion to the chain
func (r *AssertionChain) Contains(expression string, expected interface{}) *AssertionChain {
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
			if err := assertion(copyHttpResponse(res), copyHttpRequest(req)); err != nil {
				return err
			}
		}
		return nil
	}
}

func copyHttpResponse(response *http.Response) *http.Response {
	if response == nil {
		return nil
	}

	var resBodyBytes []byte
	if response.Body != nil {
		resBodyBytes, _ = ioutil.ReadAll(response.Body)
		response.Body = ioutil.NopCloser(bytes.NewBuffer(resBodyBytes))
	}

	resCopy := &http.Response{
		Header:        map[string][]string{},
		StatusCode:    response.StatusCode,
		Status:        response.Status,
		Body:          ioutil.NopCloser(bytes.NewBuffer(resBodyBytes)),
		Proto:         response.Proto,
		ProtoMinor:    response.ProtoMinor,
		ProtoMajor:    response.ProtoMajor,
		ContentLength: response.ContentLength,
	}

	for name, values := range response.Header {
		resCopy.Header[name] = values
	}

	return resCopy
}

func copyHttpRequest(request *http.Request) *http.Request {
	resCopy := &http.Request{
		Method:        request.Method,
		Host:          request.Host,
		Proto:         request.Proto,
		ProtoMinor:    request.ProtoMinor,
		ProtoMajor:    request.ProtoMajor,
		ContentLength: request.ContentLength,
		RemoteAddr:    request.RemoteAddr,
	}

	if request.Body != nil {
		bodyBytes, _ := ioutil.ReadAll(request.Body)
		resCopy.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	if request.URL != nil {
		r2URL := new(url.URL)
		*r2URL = *request.URL
		resCopy.URL = r2URL
	}

	headers := make(http.Header)
	for k, values := range request.Header {
		for _, hValue := range values {
			headers.Add(k, hValue)
		}
	}
	resCopy.Header = headers

	return resCopy
}
