package client

import "testing"

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
