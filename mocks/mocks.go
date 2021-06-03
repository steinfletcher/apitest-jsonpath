package mocks

import (
	"github.com/steinfletcher/apitest"
	"github.com/steinfletcher/apitest-jsonpath/jsonpath"
	"net/http"
)

// Contains is a convenience function to assert that a jsonpath expression extracts a value in an array
func Contains(expression string, expected interface{}) apitest.Matcher {
	return func(req *http.Request, mockReq *apitest.MockRequest) error {
		return jsonpath.Contains(expression, expected, req.Body)
	}
}

// Equal is a convenience function to assert that a jsonpath expression matches the given value
func Equal(expression string, expected interface{}) apitest.Matcher {
	return func(req *http.Request, mockReq *apitest.MockRequest) error {
		return jsonpath.Equal(expression, expected, req.Body)
	}
}

// NotEqual is a function to check json path expression value is not equal to given value
func NotEqual(expression string, expected interface{}) apitest.Matcher {
	return func(req *http.Request, mockReq *apitest.MockRequest) error {
		return jsonpath.NotEqual(expression, expected, req.Body)
	}
}

// Len asserts that value is the expected length, determined by reflect.Len
func Len(expression string, expectedLength int) apitest.Matcher {
	return func(req *http.Request, mockReq *apitest.MockRequest) error {
		return jsonpath.Length(expression, expectedLength, req.Body)
	}
}

// GreaterThan asserts that value is greater than the given length, determined by reflect.Len
func GreaterThan(expression string, minimumLength int) apitest.Matcher {
	return func(req *http.Request, mockReq *apitest.MockRequest) error {
		return jsonpath.GreaterThan(expression, minimumLength, req.Body)
	}
}