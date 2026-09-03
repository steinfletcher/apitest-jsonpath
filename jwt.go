package jsonpath

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/steinfletcher/apitest-jsonpath/jsonpath"
)

const (
	jwtHeaderIndex  = 0
	jwtPayloadIndex = 1
)

func JWTHeaderEqual(tokenSelector func(*http.Response) (string, error), expression string, expected interface{}) func(*http.Response, *http.Request) error {
	return jwtEqual(tokenSelector, expression, expected, jwtHeaderIndex)
}

func JWTPayloadEqual(tokenSelector func(*http.Response) (string, error), expression string, expected interface{}) func(*http.Response, *http.Request) error {
	return jwtEqual(tokenSelector, expression, expected, jwtPayloadIndex)
}

func jwtEqual(tokenSelector func(*http.Response) (string, error), expression string, expected interface{}, index int) func(*http.Response, *http.Request) error {
	return func(response *http.Response, request *http.Request) error {
		token, err := tokenSelector(response)
		if err != nil {
			return err
		}

		parts := strings.Split(token, ".")
		if len(parts) != 3 {
			return errors.New("invalid token: token should contain header, payload and signature")
		}

		decoded, err := base64Decode(parts[index])
		if err != nil {
			return fmt.Errorf("invalid jwt: %w", err)
		}

		value, err := jsonpath.JsonPath(bytes.NewReader(decoded), expression)
		if err != nil {
			return err
		}

		if !jsonpath.ObjectsAreEqual(value, expected) {
			return fmt.Errorf("\"%v\" not equal to \"%v\"", value, expected)
		}

		return nil
	}
}

func base64Decode(src string) ([]byte, error) {
	if l := len(src) % 4; l > 0 {
		src += strings.Repeat("=", 4-l)
	}

	decoded, err := base64.URLEncoding.DecodeString(src)
	if err != nil {
		return nil, fmt.Errorf("decoding error: %w", err)
	}
	return decoded, nil
}
