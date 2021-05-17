package jsonpath

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	regex "regexp"
	"strings"

	"github.com/PaesslerAG/jsonpath"
)

// Contains is a convenience function to assert that a jsonpath expression extracts a value in an array
func Contains(expression string, expected interface{}) func(*http.Response, *http.Request) error {
	return func(res *http.Response, req *http.Request) error {
		value, err := jsonPath(res.Body, expression)
		if err != nil {
			return err
		}

		ok, found := includesElement(value, expected)
		if !ok {
			return errors.New(fmt.Sprintf("\"%s\" could not be applied builtin len()", expected))
		}
		if !found {
			return errors.New(fmt.Sprintf("\"%s\" does not contain \"%s\"", value, expected))
		}
		return nil
	}
}

// Equal is a convenience function to assert that a jsonpath expression extracts a value
func Equal(expression string, expected interface{}) func(*http.Response, *http.Request) error {
	return func(res *http.Response, req *http.Request) error {
		value, err := jsonPath(res.Body, expression)
		if err != nil {
			return err
		}

		if !objectsAreEqual(value, expected) {
			return errors.New(fmt.Sprintf("\"%s\" not equal to \"%s\"", value, expected))
		}
		return nil
	}
}

// NotEqual is a function to check json path expression value is not equal to given value
func NotEqual(expression string, expected interface{}) func(*http.Response, *http.Request) error {
	return func(res *http.Response, req *http.Request) error {
		value, err := jsonPath(res.Body, expression)
		if err != nil {
			return err
		}

		if objectsAreEqual(value, expected) {
			return errors.New(fmt.Sprintf("\"%s\" value is equal to \"%s\"", expression, expected))
		}
		return nil
	}
}

// Len asserts that value is the expected length, determined by reflect.Len
func Len(expression string, expectedLength int) func(*http.Response, *http.Request) error {
	return func(res *http.Response, req *http.Request) error {
		value, err := jsonPath(res.Body, expression)
		if err != nil {
			return err
		}

		v := reflect.ValueOf(value)
		if v.Len() != expectedLength {
			return errors.New(fmt.Sprintf("\"%d\" not equal to \"%d\"", v.Len(), expectedLength))
		}
		return nil
	}
}

// GreaterThan asserts that value is greater than the given length, determined by reflect.Len
func GreaterThan(expression string, minimumLength int) func(*http.Response, *http.Request) error {
	return func(res *http.Response, req *http.Request) error {
		value, err := jsonPath(res.Body, expression)
		if err != nil {
			return err
		}

		v := reflect.ValueOf(value)
		if v.Len() < minimumLength {
			return errors.New(fmt.Sprintf("\"%d\" is greater than \"%d\"", v.Len(), minimumLength))
		}
		return nil
	}
}

// LessThan asserts that value is less than the given length, determined by reflect.Len
func LessThan(expression string, maximumLength int) func(*http.Response, *http.Request) error {
	return func(res *http.Response, req *http.Request) error {
		value, err := jsonPath(res.Body, expression)
		if err != nil {
			return err
		}

		v := reflect.ValueOf(value)
		if v.Len() > maximumLength {
			return errors.New(fmt.Sprintf("\"%d\" is less than \"%d\"", v.Len(), maximumLength))
		}
		return nil
	}
}

// Present asserts that value returned by the expression is present
func Present(expression string) func(*http.Response, *http.Request) error {
	return func(res *http.Response, req *http.Request) error {
		value, _ := jsonPath(res.Body, expression)
		if isEmpty(value) {
			return errors.New(fmt.Sprintf("value not present for expression: '%s'", expression))
		}
		return nil
	}
}

// NotPresent asserts that value returned by the expression is not present
func NotPresent(expression string) func(*http.Response, *http.Request) error {
	return func(res *http.Response, req *http.Request) error {
		value, _ := jsonPath(res.Body, expression)
		if !isEmpty(value) {
			return errors.New(fmt.Sprintf("value present for expression: '%s'", expression))
		}
		return nil
	}
}

// Matches asserts that the value matches the given regular expression
func Matches(expression string, regexp string) func(*http.Response, *http.Request) error {
	return func(res *http.Response, req *http.Request) error {
		pattern, err := regex.Compile(regexp)
		if err != nil {
			return errors.New(fmt.Sprintf("invalid pattern: '%s'", regexp))
		}
		value, _ := jsonPath(res.Body, expression)
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

func isEmpty(object interface{}) bool {
	if object == nil {
		return true
	}

	objValue := reflect.ValueOf(object)

	switch objValue.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		return objValue.Len() == 0
	case reflect.Ptr:
		if objValue.IsNil() {
			return true
		}
		deref := objValue.Elem().Interface()
		return isEmpty(deref)
	default:
		zero := reflect.Zero(objValue.Type())
		return reflect.DeepEqual(object, zero.Interface())
	}
}

func jsonPath(reader io.Reader, expression string) (interface{}, error) {
	v := interface{}(nil)
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &v)
	if err != nil {
		return nil, err
	}

	value, err := jsonpath.Get(expression, v)
	if err != nil {
		return nil, fmt.Errorf("evaluating '%s' resulted in error: '%s'", expression, err)
	}
	return value, nil
}

// courtesy of github.com/stretchr/testify
func includesElement(list interface{}, element interface{}) (ok, found bool) {
	listValue := reflect.ValueOf(list)
	elementValue := reflect.ValueOf(element)
	defer func() {
		if e := recover(); e != nil {
			ok = false
			found = false
		}
	}()

	if reflect.TypeOf(list).Kind() == reflect.String {
		return true, strings.Contains(listValue.String(), elementValue.String())
	}

	if reflect.TypeOf(list).Kind() == reflect.Map {
		mapKeys := listValue.MapKeys()
		for i := 0; i < len(mapKeys); i++ {
			if objectsAreEqual(mapKeys[i].Interface(), element) {
				return true, true
			}
		}
		return true, false
	}

	for i := 0; i < listValue.Len(); i++ {
		if objectsAreEqual(listValue.Index(i).Interface(), element) {
			return true, true
		}
	}
	return true, false
}

func objectsAreEqual(expected, actual interface{}) bool {
	if expected == nil || actual == nil {
		return expected == actual
	}

	exp, ok := expected.([]byte)
	if !ok {
		return reflect.DeepEqual(expected, actual)
	}

	act, ok := actual.([]byte)
	if !ok {
		return false
	}
	if exp == nil || act == nil {
		return exp == nil && act == nil
	}
	return bytes.Equal(exp, act)
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
