package common

import (
	"net/http"
	"net/url"
	"testing"
)

func TestGetQueryParam(t *testing.T) {
	r := &http.Request{
		URL: &url.URL{
			RawQuery: "foo=bar&baz=qux",
		},
	}
	if v := GetQueryParam(r, "foo"); v != "bar" {
		t.Errorf("expected 'bar', got '%s'", v)
	}
	if v := GetQueryParam(r, "baz"); v != "qux" {
		t.Errorf("expected 'qux', got '%s'", v)
	}
	if v := GetQueryParam(r, "missing"); v != "" {
		t.Errorf("expected empty string for missing param, got '%s'", v)
	}
}
