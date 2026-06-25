package cmd

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestDoctorProbesProductAPIExecutables(t *testing.T) {
	requests := map[string]string{}
	var authorization string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests[r.URL.Path] = r.URL.RawQuery
		switch r.URL.Path {
		case "/version":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"service":"loomloom-product-api"}`))
		case "/loom/v1/marketListings":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"items":[]}`))
		case "/loom/v1/users/me/executables":
			authorization = r.Header.Get("Authorization")
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"items":[]}`))
		case "/release":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"tag_name":"v0.0.1"}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()
	t.Setenv("LOOMLOOM_CLI_RELEASE_API", server.URL+"/release")

	opts := &rootOptions{
		server:  server.URL + "/loom/v1",
		token:   "token-1",
		timeout: time.Second,
		output:  "json",
	}
	cmd := newDoctorCmd(opts)
	var out bytes.Buffer
	cmd.SetOut(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("doctor command error = %v", err)
	}
	if requests["/version"] != "" {
		t.Fatalf("version query=%q want empty", requests["/version"])
	}
	if requests["/loom/v1/marketListings"] != "pageSize=1" {
		t.Fatalf("market query=%q want pageSize=1", requests["/loom/v1/marketListings"])
	}
	if requests["/loom/v1/users/me/executables"] != "pageSize=1" {
		t.Fatalf("executables query=%q want pageSize=1", requests["/loom/v1/users/me/executables"])
	}
	if authorization != "Bearer token-1" {
		t.Fatalf("authorization=%q want Bearer token-1", authorization)
	}
	if !strings.Contains(out.String(), `"healthy": true`) {
		t.Fatalf("output=%s want healthy true", out.String())
	}
	if !strings.Contains(out.String(), `"message": "ok"`) {
		t.Fatalf("output=%s want message ok", out.String())
	}
	if !strings.Contains(out.String(), `"service": "loomloom-product-api"`) {
		t.Fatalf("output=%s want Product API version", out.String())
	}
}

func TestDoctorSuppressesReleaseCheckErrorsInJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/version":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"service":"loomloom-product-api"}`))
		case "/loom/v1/marketListings":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"items":[]}`))
		case "/loom/v1/users/me/executables":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"items":[]}`))
		case "/release":
			http.NotFound(w, r)
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()
	t.Setenv("LOOMLOOM_CLI_RELEASE_API", server.URL+"/release")

	opts := &rootOptions{
		server:  server.URL + "/loom/v1",
		token:   "token-1",
		timeout: time.Second,
		output:  "json",
	}
	cmd := newDoctorCmd(opts)
	var out bytes.Buffer
	cmd.SetOut(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("doctor command error = %v", err)
	}
	if strings.Contains(out.String(), "version_check_error") {
		t.Fatalf("output=%s should not include release check error", out.String())
	}
	if !strings.Contains(out.String(), `"healthy": true`) {
		t.Fatalf("output=%s want healthy true", out.String())
	}
}

func TestDoctorAcceptsNonJSONVersionResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/version":
			w.Header().Set("Content-Type", "text/html")
			_, _ = w.Write([]byte(`<!doctype html>`))
		case "/loom/v1/marketListings":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"items":[]}`))
		case "/loom/v1/users/me/executables":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"items":[]}`))
		case "/release":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"tag_name":"v0.0.1"}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()
	t.Setenv("LOOMLOOM_CLI_RELEASE_API", server.URL+"/release")

	opts := &rootOptions{
		server:  server.URL + "/loom/v1",
		token:   "token-1",
		timeout: time.Second,
		output:  "json",
	}
	cmd := newDoctorCmd(opts)
	var out bytes.Buffer
	cmd.SetOut(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("doctor command error = %v", err)
	}
	if !strings.Contains(out.String(), `"reachable": true`) {
		t.Fatalf("output=%s want reachable true", out.String())
	}
}

func TestProductAPISystemBaseURL(t *testing.T) {
	got, err := productAPISystemBaseURL("https://example.test/gateway/loom/v1/")
	if err != nil {
		t.Fatalf("productAPISystemBaseURL error = %v", err)
	}
	if got != "https://example.test/gateway" {
		t.Fatalf("base URL=%q want https://example.test/gateway", got)
	}
}
