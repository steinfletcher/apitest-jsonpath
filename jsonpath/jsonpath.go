package jsonpath

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/PaesslerAG/jsonpath"
)

func Contains(expression string, expected any, data io.Reader) error {
	value, err := JsonPath(data, expression)
	if err != nil {
		return err
	}
	ok, found := IncludesElement(value, expected)
	if !ok {
		return fmt.Errorf("\"%v\" could not be applied builtin len()", expected)
	}
	if !found {
		return fmt.Errorf("\"%v\" does not contain \"%v\"", value, expected)
	}
	return nil
}

func Equal(expression string, expected any, data io.Reader) error {
	value, err := JsonPath(data, expression)
	if err != nil {
		return err
	}
	if !ObjectsAreEqual(value, expected) {
		return fmt.Errorf("\"%v\" not equal to \"%v\"", value, expected)
	}
	return nil
}

func NotEqual(expression string, expected any, data io.Reader) error {
	value, err := JsonPath(data, expression)
	if err != nil {
		return err
	}

	if ObjectsAreEqual(value, expected) {
		return fmt.Errorf("\"%s\" value is equal to \"%v\"", expression, expected)
	}
	return nil
}

func Length(expression string, expectedLength int, data io.Reader) error {
	length, err := lengthOf(expression, data)
	if err != nil {
		return err
	}

	if length != expectedLength {
		return fmt.Errorf("\"%d\" not equal to \"%d\"", length, expectedLength)
	}
	return nil
}

func GreaterThan(expression string, minimumLength int, data io.Reader) error {
	length, err := lengthOf(expression, data)
	if err != nil {
		return err
	}

	if length < minimumLength {
		return fmt.Errorf("\"%d\" is less than \"%d\"", length, minimumLength)
	}
	return nil
}

func LessThan(expression string, maximumLength int, data io.Reader) error {
	length, err := lengthOf(expression, data)
	if err != nil {
		return err
	}

	if length > maximumLength {
		return fmt.Errorf("\"%d\" is greater than \"%d\"", length, maximumLength)
	}
	return nil
}

// lengthOf evaluates the expression and returns the length of the result. Only arrays, slices,
// maps and strings have a length; a null result or a result of any other type is an error rather
// than a panic, so that a failed assertion is reported normally.
func lengthOf(expression string, data io.Reader) (int, error) {
	value, err := JsonPath(data, expression)
	if err != nil {
		return 0, err
	}

	if value == nil {
		return 0, errors.New("value is null")
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		return v.Len(), nil
	default:
		return 0, fmt.Errorf("value of type %s has no length", v.Kind())
	}
}

func Present(expression string, data io.Reader) error {
	value, err := JsonPath(data, expression)
	if err != nil || !isPresent(value) {
		return fmt.Errorf("value not present for expression: '%s'", expression)
	}
	return nil
}

func NotPresent(expression string, data io.Reader) error {
	value, err := JsonPath(data, expression)
	if err == nil && isPresent(value) {
		return fmt.Errorf("value present for expression: '%s'", expression)
	}
	return nil
}

func JsonPath(reader io.Reader, expression string) (any, error) {
	v := any(nil)
	b, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &v)
	if err != nil {
		return nil, err
	}

	value, err := jsonpath.Get(expression, v)
	if err != nil {
		return nil, fmt.Errorf("evaluating '%s' resulted in error: '%w'", expression, err)
	}
	return value, nil
}

// courtesy of github.com/stretchr/testify
func IncludesElement(list any, element any) (ok, found bool) {
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
		for i := range mapKeys {
			if ObjectsAreEqual(mapKeys[i].Interface(), element) {
				return true, true
			}
		}
		return true, false
	}

	for i := 0; i < listValue.Len(); i++ {
		if ObjectsAreEqual(listValue.Index(i).Interface(), element) {
			return true, true
		}
	}
	return true, false
}

func ObjectsAreEqual(expected, actual any) bool {
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

// isPresent reports whether the expression produced a value. A JSON null is treated as absent,
// as is an empty array or object, which is what the evaluator returns when a wildcard or filter
// expression matches nothing. Every other value is present, including false, 0 and "".
func isPresent(value any) bool {
	if value == nil {
		return false
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice:
		return v.Len() > 0
	default:
		return true
	}
}
