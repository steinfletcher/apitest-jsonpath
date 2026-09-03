package http

import (
	"bytes"
	"io"
	nethttp "net/http"
	"net/http/httptest"
	"testing"
)

func TestCopyResponse(t *testing.T) {
	original := &nethttp.Response{
		StatusCode: nethttp.StatusOK,
		Header:     nethttp.Header{"X-Custom": {"a"}},
		Body:       io.NopCloser(bytes.NewBufferString("body")),
	}

	copied := CopyResponse(original)
	copied.Header["X-Custom"][0] = "changed"
	copied.Header.Add("X-Custom", "b")

	if got := original.Header["X-Custom"]; len(got) != 1 || got[0] != "a" {
		t.Fatalf("expected the original headers to be untouched, got %v", got)
	}
	for name, body := range map[string]*nethttp.Response{"copy": copied, "original": original} {
		b, _ := io.ReadAll(body.Body)
		if string(b) != "body" {
			t.Fatalf("expected %s body to be readable, got %q", name, b)
		}
	}
	if CopyResponse(nil) != nil {
		t.Fatal("expected a nil response to copy as nil")
	}
}

func TestCopyRequest(t *testing.T) {
	original := httptest.NewRequest(nethttp.MethodPost, "/path?a=1", bytes.NewBufferString("body"))
	original.Header.Set("X-Custom", "a")

	copied := CopyRequest(original)
	copied.Header.Set("X-Custom", "changed")
	copied.URL.Path = "/other"

	if original.Header.Get("X-Custom") != "a" || original.URL.Path != "/path" {
		t.Fatalf("expected the original request to be untouched, got %v %s", original.Header, original.URL)
	}
	for name, req := range map[string]*nethttp.Request{"copy": copied, "original": original} {
		b, _ := io.ReadAll(req.Body)
		if string(b) != "body" {
			t.Fatalf("expected %s body to be readable, got %q", name, b)
		}
	}
	if CopyRequest(nil) != nil {
		t.Fatal("expected a nil request to copy as nil")
	}
}
