package glhttp

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	jsoniter "github.com/json-iterator/go"
)

var testJSONAPI = jsoniter.ConfigCompatibleWithStandardLibrary

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func TestSuccessResponse(t *testing.T) {
	resp := Success(map[string]string{"id": "1"})
	if resp.Code != 0 || resp.Message != "success" {
		t.Fatalf("unexpected response: %+v", resp)
	}
	if resp.Data["id"] != "1" {
		t.Fatalf("Data[id] = %q", resp.Data["id"])
	}
}

func TestFailResponse(t *testing.T) {
	resp := Fail(40001, "bad request")
	if resp.Code != 40001 || resp.Message != "bad request" {
		t.Fatalf("unexpected response: %+v", resp)
	}
	if resp.Data != nil {
		t.Fatalf("Data = %#v, want nil", resp.Data)
	}
}

func TestClientGetJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("Method = %s, want GET", r.Method)
		}
		if r.Header.Get("X-Token") != "abc" {
			t.Fatalf("X-Token = %q", r.Header.Get("X-Token"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gltools"}`))
	}))
	defer server.Close()

	var out struct {
		Name string `json:"name"`
	}
	err := NewClient(time.Second).GetJSON(context.Background(), server.URL, map[string]string{"X-Token": "abc"}, &out)
	if err != nil {
		t.Fatal(err)
	}
	if out.Name != "gltools" {
		t.Fatalf("Name = %q", out.Name)
	}
}

func TestClientPostJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("Method = %s, want POST", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Fatalf("Content-Type = %q", r.Header.Get("Content-Type"))
		}
		if r.Header.Get("X-Token") != "abc" {
			t.Fatalf("X-Token = %q", r.Header.Get("X-Token"))
		}
		var body struct {
			Name string `json:"name"`
		}
		if err := testJSONAPI.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.Name != "gltools" {
			t.Fatalf("body.Name = %q", body.Name)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	var out struct {
		OK bool `json:"ok"`
	}
	err := NewClient(time.Second).PostJSON(context.Background(), server.URL, map[string]string{"X-Token": "abc"}, map[string]string{"name": "gltools"}, &out)
	if err != nil {
		t.Fatal(err)
	}
	if !out.OK {
		t.Fatal("OK = false, want true")
	}
}

func TestClientPostForm(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("Method = %s, want POST", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
			t.Fatalf("Content-Type = %q", r.Header.Get("Content-Type"))
		}
		if r.Header.Get("X-Token") != "abc" {
			t.Fatalf("X-Token = %q", r.Header.Get("X-Token"))
		}
		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.Form.Get("name") != "gltools" {
			t.Fatalf("form name = %q", r.Form.Get("name"))
		}
		if r.Form.Get("type") != "toolkit" {
			t.Fatalf("form type = %q", r.Form.Get("type"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	var out struct {
		OK bool `json:"ok"`
	}
	form := url.Values{
		"name": []string{"gltools"},
		"type": []string{"toolkit"},
	}
	err := NewClient(time.Second).PostForm(context.Background(), server.URL, map[string]string{"X-Token": "abc"}, form, &out)
	if err != nil {
		t.Fatal(err)
	}
	if !out.OK {
		t.Fatal("OK = false, want true")
	}
}

func TestClientReturnsErrorForNon2xx(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "failed", http.StatusInternalServerError)
	}))
	defer server.Close()

	err := NewClient(time.Second).GetJSON(context.Background(), server.URL, nil, nil)
	if err == nil {
		t.Fatal("GetJSON() error = nil, want non-nil")
	}
}

func TestClientUsesInjectedHTTPClient(t *testing.T) {
	called := false
	httpClient := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			called = true
			if req.URL.String() != "https://example.com/data" {
				t.Fatalf("URL = %q, want https://example.com/data", req.URL.String())
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader(`{"ok":true}`)),
				Request:    req,
			}, nil
		}),
	}

	var out struct {
		OK bool `json:"ok"`
	}
	err := NewClientWithHTTPClient(httpClient).GetJSON(context.Background(), "https://example.com/data", nil, &out)
	if err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("injected HTTP client was not used")
	}
	if !out.OK {
		t.Fatal("OK = false, want true")
	}
}
