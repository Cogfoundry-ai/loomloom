package client

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestNormalizeBaseURLRequiresExplicitServer(t *testing.T) {
	_, err := normalizeBaseURL("")
	if err == nil {
		t.Fatal("expected empty server URL to fail")
	}
}

func TestNormalizeBaseURLAddsScheme(t *testing.T) {
	got, err := normalizeBaseURL("api.cogfoundry.example/loom/v1")
	if err != nil {
		t.Fatalf("normalizeBaseURL failed: %v", err)
	}
	want := "https://api.cogfoundry.example/loom/v1"
	if got != want {
		t.Fatalf("normalizeBaseURL = %q, want %q", got, want)
	}
}

func TestEndpointUsesConfiguredBaseURL(t *testing.T) {
	c, err := New(Config{BaseURL: "https://api-test.cogfoundry.example/loom/v1"})
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	got := c.endpoint("/templates")
	want := "https://api-test.cogfoundry.example/loom/v1/templates"
	if got != want {
		t.Fatalf("endpoint = %q, want %q", got, want)
	}
}

func TestEndpointNormalizesMissingLeadingSlash(t *testing.T) {
	c, err := New(Config{BaseURL: "https://api-test.cogfoundry.example/loom/v1"})
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	got := c.endpoint("marketListings")
	want := "https://api-test.cogfoundry.example/loom/v1/marketListings"
	if got != want {
		t.Fatalf("endpoint = %q, want %q", got, want)
	}
}

func TestEndpointDoesNotRewriteAbsolutePath(t *testing.T) {
	c, err := New(Config{BaseURL: "https://api-test.cogfoundry.example/loom/v1"})
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	got := c.endpoint("/officialTemplates")
	want := "https://api-test.cogfoundry.example/loom/v1/officialTemplates"
	if got != want {
		t.Fatalf("endpoint = %q, want %q", got, want)
	}
}

func TestVerboseHTTPLogsExcludeTokenQueryAndBody(t *testing.T) {
	const token = "secret-token-value"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer "+token {
			t.Fatalf("Authorization=%q", got)
		}
		w.Header().Set("X-Request-ID", "request-123")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	var logs bytes.Buffer
	c, err := New(Config{
		BaseURL:   server.URL,
		Token:     token,
		Verbose:   true,
		LogWriter: &logs,
	})
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	var response map[string]any
	err = c.PostProductJSONWithQuery(
		context.Background(),
		"/templates/template-1:run",
		url.Values{"sensitive": []string{"query-secret"}},
		map[string]string{"prompt": "body-secret"},
		&response,
	)
	if err != nil {
		t.Fatalf("PostProductJSONWithQuery failed: %v", err)
	}

	got := logs.String()
	for _, want := range []string{
		"[debug] POST /templates/template-1:run",
		"response status=200",
		"request_id=request-123",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("logs=%q want %q", got, want)
		}
	}
	for _, secret := range []string{token, "query-secret", "body-secret", "Authorization"} {
		if strings.Contains(got, secret) {
			t.Fatalf("logs contain sensitive value %q: %s", secret, got)
		}
	}
}

func TestHTTPLogsAreDisabledByDefault(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	var logs bytes.Buffer
	c, err := New(Config{BaseURL: server.URL, LogWriter: &logs})
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	var response map[string]any
	if err := c.GetJSON(context.Background(), "/health", &response); err != nil {
		t.Fatalf("GetJSON failed: %v", err)
	}
	if logs.Len() != 0 {
		t.Fatalf("logs=%q want no default logs", logs.String())
	}
}
