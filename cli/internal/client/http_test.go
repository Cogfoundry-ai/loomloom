package client

import "testing"

func TestNormalizeBaseURLRequiresExplicitServer(t *testing.T) {
	_, err := normalizeBaseURL("")
	if err == nil {
		t.Fatal("expected empty server URL to fail")
	}
}

func TestNormalizeBaseURLAddsScheme(t *testing.T) {
	got, err := normalizeBaseURL("loomloom.shengsuanyun.com/batch")
	if err != nil {
		t.Fatalf("normalizeBaseURL failed: %v", err)
	}
	want := "https://loomloom.shengsuanyun.com/batch"
	if got != want {
		t.Fatalf("normalizeBaseURL = %q, want %q", got, want)
	}
}
