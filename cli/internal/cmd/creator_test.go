package cmd

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestCreatorEarningsUsesTokenPrincipal(t *testing.T) {
	var requestedPath string
	var requestedQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedPath = r.URL.Path
		requestedQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"items":[]}`))
	}))
	defer server.Close()

	opts := &rootOptions{
		server:  server.URL + "/loom/v1",
		timeout: time.Second,
	}
	cmd := newCreatorEarningsCmd(opts)
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--limit", "25"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("creator earnings command error = %v", err)
	}
	if requestedPath != "/loom/v1/creators/me/earnings" {
		t.Fatalf("path=%q want creator earnings endpoint", requestedPath)
	}
	if strings.Contains(requestedQuery, "creator_user_id=") {
		t.Fatalf("query %q should not include creator_user_id", requestedQuery)
	}
	for _, want := range []string{"pageSize=25"} {
		if !strings.Contains(requestedQuery, want) {
			t.Fatalf("query %q missing %q", requestedQuery, want)
		}
	}
	if !strings.Contains(out.String(), `"items": []`) {
		t.Fatalf("unexpected output: %s", out.String())
	}
}
