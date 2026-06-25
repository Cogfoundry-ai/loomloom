package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestMarketPublishBuildsRequestWithoutGeneratedFields(t *testing.T) {
	var requestedPath string
	var body map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedPath = r.URL.Path
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":"listing-1"}`))
	}))
	defer server.Close()

	opts := &rootOptions{server: server.URL + "/loom/v1", timeout: time.Second}
	cmd := newListingPublishCmd(opts)
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{
		"template-1",
		"--listing-id", "listing-1",
		"--template-version-id", "version-1",
		"--display-name", "PRD Review Bot",
		"--description", "Review PRD docs",
		"--task-fixed-fee-t", "0",
	})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("market publish command error = %v", err)
	}
	if requestedPath != "/loom/v1/marketListings" {
		t.Fatalf("path=%q want publish endpoint", requestedPath)
	}
	if _, ok := body["creator_user_id"]; ok {
		t.Fatalf("publish body should not include creator_user_id: %#v", body)
	}
	if body["taskFixedFeeT"] != float64(0) {
		t.Fatalf("taskFixedFeeT=%v want 0", body["taskFixedFeeT"])
	}
	if body["displayName"] != "PRD Review Bot" {
		t.Fatalf("displayName=%v want PRD Review Bot", body["displayName"])
	}
	if body["templateId"] != "template-1" {
		t.Fatalf("templateId=%v want template-1", body["templateId"])
	}
	if body["listingId"] != "listing-1" {
		t.Fatalf("listingId=%v want listing-1", body["listingId"])
	}
	if body["templateVersionId"] != "version-1" {
		t.Fatalf("templateVersionId=%v want version-1", body["templateVersionId"])
	}
	if _, ok := body["workflow_definition"]; ok {
		t.Fatalf("publish body should not include workflow_definition: %#v", body)
	}
	if _, ok := body["definition_hash"]; ok {
		t.Fatalf("publish body should not include definition_hash: %#v", body)
	}
	if !strings.Contains(out.String(), `"id": "listing-1"`) {
		t.Fatalf("unexpected output: %s", out.String())
	}
}

func TestMarketListSendsProductAPIQueryParams(t *testing.T) {
	var requestedPath string
	var requestedQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedPath = r.URL.Path
		requestedQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"items":[]}`))
	}))
	defer server.Close()

	opts := &rootOptions{server: server.URL + "/loom/v1", timeout: time.Second}
	cmd := newMarketListCmd(opts)
	cmd.SetArgs([]string{
		"--keyword", "prd",
		"--page-size", "5",
		"--page-token", "next",
		"--order-by", "createdAt desc",
	})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("market list command error = %v", err)
	}
	if requestedPath != "/loom/v1/marketListings" {
		t.Fatalf("path=%q want market listings endpoint", requestedPath)
	}
	for _, want := range []string{"keyword=prd", "pageSize=5", "pageToken=next", "orderBy=createdAt+desc"} {
		if !strings.Contains(requestedQuery, want) {
			t.Fatalf("query=%q missing %s", requestedQuery, want)
		}
	}
}

func TestMarketListRejectsLimitAndPageSizeTogether(t *testing.T) {
	opts := &rootOptions{server: "https://example.test/loom/v1", timeout: time.Second}
	cmd := newMarketListCmd(opts)
	cmd.SetArgs([]string{"--limit", "5", "--page-size", "10"})

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "--limit and --page-size cannot be used together") {
		t.Fatalf("error=%v want limit/page-size conflict", err)
	}
}

func TestMarketShowUsesDetailEndpoint(t *testing.T) {
	var requestedPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"listing-1"}`))
	}))
	defer server.Close()

	opts := &rootOptions{server: server.URL + "/loom/v1", timeout: time.Second}
	cmd := newMarketShowCmd(opts)
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"listing-1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("market show command error = %v", err)
	}
	if requestedPath != "/loom/v1/marketListings/listing-1" {
		t.Fatalf("path=%q want detail endpoint", requestedPath)
	}
	if !strings.Contains(out.String(), `"id": "listing-1"`) {
		t.Fatalf("unexpected output: %s", out.String())
	}
}

func TestMarketPublishRequiresExplicitTaskFixedFee(t *testing.T) {
	opts := &rootOptions{server: "https://example.test", timeout: time.Second}
	cmd := newListingPublishCmd(opts)
	cmd.SetArgs([]string{
		"template-1",
		"--template-version-id", "version-1",
		"--display-name", "PRD Review Bot",
	})

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "task-fixed-fee-t") {
		t.Fatalf("error=%v want task fee required", err)
	}
}

func TestDeprecatedMarketCommandsRemainAvailable(t *testing.T) {
	cmd := newMarketCmd(&rootOptions{})
	if found, _, err := cmd.Find([]string{"publish"}); err != nil || found.Name() != "publish" {
		t.Fatalf("market publish compatibility command missing: found=%v err=%v", found, err)
	}
	if found, _, err := cmd.Find([]string{"relist"}); err != nil || found.Name() != "relist" {
		t.Fatalf("market relist compatibility command missing: found=%v err=%v", found, err)
	}
}

func TestMarketQuoteRemovesIdentityFieldsAndPreservesPayload(t *testing.T) {
	var body map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/loom/v1/marketListings/listing-1:quote" {
			t.Fatalf("path=%q want quote endpoint", r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"quote_id":"quote-1"}`))
	}))
	defer server.Close()

	inputPath := writeMarketInputFile(t, `{
		"listingVersionId": "lv-1",
		"user_id": 1,
		"buyer_user_id": 13005,
		"creator_user_id": 13004,
		"taskInputs": [{"row": 1}],
		"workflow_definition": {"name": "wf"}
	}`)
	opts := &rootOptions{server: server.URL + "/loom/v1", timeout: time.Second}
	cmd := newMarketQuoteCmd(opts)
	cmd.SetArgs([]string{"listing-1", "--input-file", inputPath})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("market quote command error = %v", err)
	}
	for _, field := range []string{"user_id", "buyer_user_id", "creator_user_id"} {
		if _, ok := body[field]; ok {
			t.Fatalf("quote body should not include %s: %#v", field, body)
		}
	}
	if body["listingVersionId"] != "lv-1" {
		t.Fatalf("listingVersionId=%v want lv-1", body["listingVersionId"])
	}
	if _, ok := body["taskInputs"]; !ok {
		t.Fatalf("taskInputs was not preserved: %#v", body)
	}
	if _, ok := body["workflow_definition"]; !ok {
		t.Fatalf("workflow_definition was not preserved: %#v", body)
	}
}

func TestMarketRunRequiresConfirmBeforeRequest(t *testing.T) {
	called := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer server.Close()

	inputPath := writeMarketInputFile(t, `{"taskInputs":[]}`)
	opts := &rootOptions{server: server.URL + "/loom/v1", timeout: time.Second}
	cmd := newMarketRunCmd(opts)
	cmd.SetArgs([]string{"listing-1", "--input-file", inputPath})

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "--confirm is required") {
		t.Fatalf("error=%v want confirm required", err)
	}
	if called {
		t.Fatal("server should not be called without --confirm")
	}
}

func TestMarketRunRemovesIdentityFieldsAndAddsConfirm(t *testing.T) {
	var body map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/loom/v1/marketListings/listing-1:execute" {
			t.Fatalf("path=%q want run endpoint", r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"runId":"run-1","runTransactionId":"transaction-1"}`))
	}))
	defer server.Close()

	inputPath := writeMarketInputFile(t, `{
		"clientRequestId":"req-1",
		"user_id": 1,
		"userId": 2,
		"buyer_user_id": 13005,
		"buyerUserId": 13006,
		"creator_user_id": 13004,
		"creatorUserId": 13007,
		"taskInputs":[]
	}`)
	var logs bytes.Buffer
	opts := &rootOptions{
		server:    server.URL + "/loom/v1",
		timeout:   time.Second,
		verbose:   true,
		logWriter: &logs,
	}
	cmd := newMarketRunCmd(opts)
	cmd.SetArgs([]string{"listing-1", "--input-file", inputPath, "--confirm"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("market run command error = %v", err)
	}
	for _, field := range []string{"user_id", "userId", "buyer_user_id", "buyerUserId", "creator_user_id", "creatorUserId"} {
		if _, ok := body[field]; ok {
			t.Fatalf("run body should not include %s: %#v", field, body)
		}
	}
	if body["confirm"] != true {
		t.Fatalf("confirm=%v want true", body["confirm"])
	}
	if body["clientRequestId"] != "req-1" {
		t.Fatalf("clientRequestId=%v want req-1", body["clientRequestId"])
	}
	if !strings.Contains(logs.String(), "market run: submitted listing_id=listing-1 run_id=run-1 transaction_id=transaction-1") {
		t.Fatalf("logs=%q want market run submission identifiers", logs.String())
	}
}

func TestMarketRunPrintsGeneratedClientRequestIDBeforeRequestFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"error":"temporary failure"}`, http.StatusServiceUnavailable)
	}))
	defer server.Close()

	inputPath := writeMarketInputFile(t, `{"taskInputs":[]}`)
	opts := &rootOptions{server: server.URL + "/loom/v1", timeout: time.Second}
	cmd := newMarketRunCmd(opts)
	var stderr bytes.Buffer
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{"listing-1", "--input-file", inputPath, "--confirm"})

	if err := cmd.Execute(); err == nil {
		t.Fatal("market run error = nil, want request failure")
	}
	if !strings.Contains(stderr.String(), "clientRequestId: loomloom-cli-") {
		t.Fatalf("stderr=%q want generated clientRequestId before request failure", stderr.String())
	}
}

func TestMarketRunGeneratesClientRequestIDWhenMissing(t *testing.T) {
	var body map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"run_id":"run-1"}`))
	}))
	defer server.Close()

	inputPath := writeMarketInputFile(t, `{"taskInputs":[]}`)
	opts := &rootOptions{server: server.URL + "/loom/v1", timeout: time.Second}
	cmd := newMarketRunCmd(opts)
	var stderr bytes.Buffer
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{"listing-1", "--input-file", inputPath, "--confirm"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("market run command error = %v", err)
	}
	clientRequestID, ok := body["clientRequestId"].(string)
	if !ok || !strings.HasPrefix(clientRequestID, "loomloom-cli-") {
		t.Fatalf("clientRequestId=%#v want generated loomloom-cli-*", body["clientRequestId"])
	}
	if body["confirm"] != true {
		t.Fatalf("confirm=%v want true", body["confirm"])
	}
	if !strings.Contains(stderr.String(), "clientRequestId: "+clientRequestID) {
		t.Fatalf("stderr=%q want generated clientRequestId", stderr.String())
	}
}

func TestMarketRunClientRequestIDFlagOverridesInputFile(t *testing.T) {
	var body map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"run_id":"run-1"}`))
	}))
	defer server.Close()

	inputPath := writeMarketInputFile(t, `{"clientRequestId":"from-file","taskInputs":[]}`)
	opts := &rootOptions{server: server.URL + "/loom/v1", timeout: time.Second}
	cmd := newMarketRunCmd(opts)
	cmd.SetArgs([]string{
		"listing-1",
		"--input-file", inputPath,
		"--client-request-id", "from-flag",
		"--confirm",
	})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("market run command error = %v", err)
	}
	if body["clientRequestId"] != "from-flag" {
		t.Fatalf("clientRequestId=%v want from-flag", body["clientRequestId"])
	}
}

func TestMarketRelistUsesListEndpointWithoutCreatorUserID(t *testing.T) {
	var requestedPath string
	var requestedQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedPath = r.URL.Path
		requestedQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"listing-1"}`))
	}))
	defer server.Close()

	opts := &rootOptions{server: server.URL + "/loom/v1", timeout: time.Second}
	cmd := newListingRelistCmd(opts)
	cmd.SetArgs([]string{"listing-1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("listing relist command error = %v", err)
	}
	if requestedPath != "/loom/v1/marketListings/listing-1:list" {
		t.Fatalf("path=%q want relist endpoint", requestedPath)
	}
	if requestedQuery != "" {
		t.Fatalf("query=%q want no identity query", requestedQuery)
	}
}

func writeMarketInputFile(t *testing.T, content string) string {
	t.Helper()
	filePath := filepath.Join(t.TempDir(), "request.json")
	if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
		t.Fatalf("write input file: %v", err)
	}
	return filePath
}
