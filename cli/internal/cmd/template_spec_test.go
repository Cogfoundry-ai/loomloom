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

func TestLoadTemplateSpecFile_ValidSpec(t *testing.T) {
	path := filepath.Join(t.TempDir(), "spec.json")
	content := `{
  "Meta": {"Name": "Spec Test", "Description": "desc"},
  "Steps": [{"StepID": "stp_text", "DisplayName": "Text", "ExecutionUnit": "text-generate"}],
  "InputSchema": {"Fields": [{"Key": "prompt", "Label": "Prompt", "ValueType": "string"}]},
  "FieldBindings": [{"FieldKey": "prompt", "StepID": "stp_text", "ParamKey": "prompt", "BindMode": "shared"}]
}`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write spec: %v", err)
	}

	spec, raw, err := loadTemplateSpecFile(path)
	if err != nil {
		t.Fatalf("loadTemplateSpecFile() error = %v", err)
	}
	if spec.Meta.Name != "Spec Test" {
		t.Fatalf("Meta.Name = %q, want Spec Test", spec.Meta.Name)
	}
	if len(raw) == 0 || raw[0] != '{' {
		t.Fatalf("expected compact JSON bytes, got %q", string(raw))
	}
	if !strings.Contains(string(raw), `"meta"`) {
		t.Fatalf("expected normalized lowerCamel TemplateSpec JSON, got %s", string(raw))
	}
	if strings.Contains(string(raw), `"Meta"`) {
		t.Fatalf("expected PascalCase keys to be normalized, got %s", string(raw))
	}
}

func TestLoadTemplateSpecFile_MissingName(t *testing.T) {
	path := filepath.Join(t.TempDir(), "spec.json")
	content := `{
  "Meta": {},
  "Steps": [{"StepID": "stp_text"}],
  "InputSchema": {"Fields": []},
  "FieldBindings": []
}`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write spec: %v", err)
	}

	if _, _, err := loadTemplateSpecFile(path); err == nil {
		t.Fatal("loadTemplateSpecFile() error = nil, want missing name error")
	}
}

func TestLoadTemplateSpecFile_NormalizesPascalCaseUpstreamBindings(t *testing.T) {
	path := filepath.Join(t.TempDir(), "spec.json")
	content := `{
  "Meta": {"Name": "Initial Input Binding Spec"},
  "Steps": [{
    "StepID": "stp_text",
    "DisplayName": "Text",
    "ExecutionUnit": "text-generate",
    "UpstreamBindings": [{
      "InputPort": "prompt",
      "SourceType": "initial_input",
      "SourceInputKey": "patent_input"
    }]
  }],
  "InputSchema": {"Fields": [{
    "Key": "patent_input",
    "Label": "Patent Input",
    "ValueType": "text_reference",
    "AcceptedMIMETypes": ["text/plain"]
  }]},
  "FieldBindings": []
}`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write spec: %v", err)
	}

	_, raw, err := loadTemplateSpecFile(path)
	if err != nil {
		t.Fatalf("loadTemplateSpecFile() error = %v", err)
	}
	normalized := string(raw)
	if !strings.Contains(normalized, `"upstreamBindings"`) {
		t.Fatalf("normalized spec missing upstreamBindings: %s", normalized)
	}
	if strings.Contains(normalized, `"UpstreamBindings"`) {
		t.Fatalf("normalized spec still has PascalCase UpstreamBindings: %s", normalized)
	}
	if !strings.Contains(normalized, `"acceptedMimeTypes"`) {
		t.Fatalf("normalized spec missing acceptedMimeTypes: %s", normalized)
	}
}

func TestLoadTemplateSpecFile_AllowsParamBindingOnlySpec(t *testing.T) {
	path := filepath.Join(t.TempDir(), "spec.json")
	content := `{
  "meta": {"name": "Param Binding Spec"},
  "steps": [{"stepId": "stp_text", "displayName": "Text", "executionUnit": "text-generate"}],
  "inputSchema": {"fields": [{"key": "prompt", "label": "Prompt", "valueType": "string"}]},
  "paramBindings": [{
    "stepId": "stp_text",
    "paramKey": "prompt",
    "bindMode": "shared",
    "sources": [{"kind": "field_ref", "fieldKey": "prompt"}]
  }]
}`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write spec: %v", err)
	}

	spec, raw, err := loadTemplateSpecFile(path)
	if err != nil {
		t.Fatalf("loadTemplateSpecFile() error = %v", err)
	}
	if len(spec.FieldBindings) != 0 {
		t.Fatalf("FieldBindings len = %d, want 0", len(spec.FieldBindings))
	}
	if len(spec.ParamBindings) != 1 {
		t.Fatalf("ParamBindings len = %d, want 1", len(spec.ParamBindings))
	}
	if !strings.Contains(string(raw), `"paramBindings"`) {
		t.Fatalf("normalized spec missing paramBindings: %s", string(raw))
	}
}

func TestLoadTemplateSpecFile_RejectsTextReferenceFieldBindingToPrompt(t *testing.T) {
	path := filepath.Join(t.TempDir(), "spec.json")
	content := `{
  "meta": {"name": "Invalid Text Reference Binding"},
  "steps": [{"stepId": "stp_text", "displayName": "Text", "executionUnit": "text-generate"}],
  "inputSchema": {"fields": [{
    "key": "patent_input",
    "label": "Patent Input",
    "valueType": "text_reference",
    "acceptedMimeTypes": ["text/plain"]
  }]},
  "fieldBindings": [{
    "fieldKey": "patent_input",
    "stepId": "stp_text",
    "paramKey": "prompt",
    "bindMode": "shared"
  }]
}`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write spec: %v", err)
	}

	_, _, err := loadTemplateSpecFile(path)
	if err == nil {
		t.Fatal("loadTemplateSpecFile() error = nil, want text_reference binding error")
	}
	if !strings.Contains(err.Error(), "text_reference") || !strings.Contains(err.Error(), "initial_input") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadTemplateSpecFile_AllowsThreeVisibleFieldParamBinding(t *testing.T) {
	path := filepath.Join(t.TempDir(), "spec.json")
	content := `{
  "meta": {"name": "Three Field Prompt Spec"},
  "steps": [{"stepId": "stp_text", "displayName": "Text", "executionUnit": "text-generate"}],
  "inputSchema": {"fields": [
    {"key": "body", "label": "Body content", "valueType": "string"},
    {"key": "style", "label": "Style requirements", "valueType": "string"},
    {"key": "format", "label": "Output format", "valueType": "string"}
  ]},
  "paramBindings": [{
    "stepId": "stp_text",
    "paramKey": "prompt",
    "bindMode": "shared",
    "separator": "\n\n",
    "sources": [
      {"kind": "field_ref", "fieldKey": "body"},
      {"kind": "field_ref", "fieldKey": "style"},
      {"kind": "field_ref", "fieldKey": "format"}
    ]
  }]
}`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write spec: %v", err)
	}

	spec, raw, err := loadTemplateSpecFile(path)
	if err != nil {
		t.Fatalf("loadTemplateSpecFile() error = %v", err)
	}
	if len(spec.ParamBindings) != 1 {
		t.Fatalf("ParamBindings len = %d, want 1", len(spec.ParamBindings))
	}
	if !strings.Contains(string(raw), `"fieldKey":"format"`) {
		t.Fatalf("normalized spec missing third field source: %s", string(raw))
	}
}

func TestTemplateSpecCheckCmdCountsParamBindings(t *testing.T) {
	path := filepath.Join(t.TempDir(), "spec.json")
	content := `{
  "meta": {"name": "Param Binding Spec"},
  "steps": [{"stepId": "stp_text", "displayName": "Text", "executionUnit": "text-generate"}],
  "inputSchema": {"fields": [{"key": "prompt", "label": "Prompt", "valueType": "string"}]},
  "paramBindings": [{
    "stepId": "stp_text",
    "paramKey": "prompt",
    "bindMode": "shared",
    "sources": [{"kind": "field_ref", "fieldKey": "prompt"}]
  }]
}`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write spec: %v", err)
	}

	opts := &rootOptions{output: "text"}
	cmd := newTemplateSpecCheckCmd(opts)
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{path})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("check command error = %v", err)
	}
	if !strings.Contains(out.String(), "bindings\t1") {
		t.Fatalf("output missing binding count: %s", out.String())
	}
}

func TestTemplateSpecDocsCmdListsTopics(t *testing.T) {
	opts := &rootOptions{output: "text"}
	cmd := newTemplateSpecDocsCmd(opts)
	var out bytes.Buffer
	cmd.SetOut(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("docs command error = %v", err)
	}
	output := out.String()
	for _, want := range []string{"spec", "authoring", "examples", "conversation", "loomloom template-spec docs <topic>"} {
		if !strings.Contains(output, want) {
			t.Fatalf("output missing %q: %s", want, output)
		}
	}
}

func TestTemplateSpecDocsCmdPrintsConversation(t *testing.T) {
	opts := &rootOptions{output: "text"}
	cmd := newTemplateSpecDocsCmd(opts)
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"conversation"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("docs conversation command error = %v", err)
	}
	output := out.String()
	for _, want := range []string{"# Conversational Template Authoring", "TemplatePlan", "Ask one question at a time"} {
		if !strings.Contains(output, want) {
			t.Fatalf("output missing %q: %s", want, output)
		}
	}
}

func TestTemplateSpecDocsCmdPrintsSpec(t *testing.T) {
	opts := &rootOptions{output: "text"}
	cmd := newTemplateSpecDocsCmd(opts)
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"spec"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("docs spec command error = %v", err)
	}
	output := out.String()
	for _, want := range []string{"# TemplateSpec Specification", "Step-Level Fan-In", "target port allows multiple inputs", "at most three regular visible field sources"} {
		if !strings.Contains(output, want) {
			t.Fatalf("output missing %q: %s", want, output)
		}
	}
}

func TestTemplateSpecDocsCmdSupportsJSON(t *testing.T) {
	opts := &rootOptions{output: "json"}
	cmd := newTemplateSpecDocsCmd(opts)
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"examples"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("docs examples command error = %v", err)
	}
	var payload map[string]any
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatalf("decode docs JSON: %v", err)
	}
	if payload["topic"] != "examples" {
		t.Fatalf("topic=%v want examples", payload["topic"])
	}
	content, _ := payload["content"].(string)
	for _, want := range []string{"Step-Level Fan-In Review Summary", "stp_prod01", "stp_summary"} {
		if !strings.Contains(content, want) {
			t.Fatalf("content missing %q: %s", want, content)
		}
	}
}

func TestTemplateSpecModelsCmdListsAvailableModels(t *testing.T) {
	var requestedPath string
	var requestedQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedPath = r.URL.Path
		requestedQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"models": [
				{
					"modelId": "google/gemini-2.5-flash",
					"displayName": "Gemini 2.5 Flash",
					"provider": "vertex",
					"executionAdapter": "vertex",
					"supportedStepTypes": ["text-generate"],
					"available": true,
					"isDefault": true
				}
			]
		}`))
	}))
	defer server.Close()

	opts := &rootOptions{
		server:  server.URL,
		timeout: time.Second,
		output:  "text",
	}
	cmd := newTemplateSpecModelsCmd(opts)
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"text-generate"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("models command error = %v", err)
	}
	if requestedPath != "/v1/batch/models" {
		t.Fatalf("path=%q want /v1/batch/models", requestedPath)
	}
	for _, want := range []string{"stepType=text-generate", "onlyAvailable=true"} {
		if !strings.Contains(requestedQuery, want) {
			t.Fatalf("query %q missing %q", requestedQuery, want)
		}
	}
	if strings.Contains(requestedQuery, "provider=") {
		t.Fatalf("query %q should not include provider by default", requestedQuery)
	}
	if !strings.Contains(out.String(), "google/gemini-2.5-flash") {
		t.Fatalf("output missing model id: %s", out.String())
	}
}

func TestTemplateSpecModelsCmdCanFilterProvider(t *testing.T) {
	var requestedQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"models":[]}`))
	}))
	defer server.Close()

	opts := &rootOptions{
		server:  server.URL,
		timeout: time.Second,
		output:  "json",
	}
	cmd := newTemplateSpecModelsCmd(opts)
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"text-generate", "--provider", "vertex"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("models command error = %v", err)
	}
	if !strings.Contains(requestedQuery, "provider=vertex") {
		t.Fatalf("query %q missing provider=vertex", requestedQuery)
	}
}

func TestTemplateSpecCreateVersionPostsCanonicalSpec(t *testing.T) {
	path := filepath.Join(t.TempDir(), "spec.json")
	content := `{
  "Meta": {"Name": "Spec Test", "Description": "desc"},
  "Steps": [{"StepID": "stp_text", "DisplayName": "Text", "ExecutionUnit": "text-generate"}],
  "InputSchema": {"Fields": [{"Key": "prompt", "Label": "Prompt", "ValueType": "string"}]},
  "FieldBindings": [{"FieldKey": "prompt", "StepID": "stp_text", "ParamKey": "prompt", "BindMode": "shared"}]
}`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write spec: %v", err)
	}

	var requestedPath string
	var payload map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedPath = r.URL.Path
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"versionId":"ver_123","versionNumber":2,"definitionHash":"hash_123","createdAt":"1777699967"}`))
	}))
	defer server.Close()

	opts := &rootOptions{
		server:  server.URL,
		timeout: time.Second,
		output:  "json",
	}
	cmd := newTemplateSpecCreateVersionCmd(opts)
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"tmpl_123", path, "--version-note", "fix judge template"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("create-version command error = %v", err)
	}
	if requestedPath != "/v1/user-templates/tmpl_123/versions" {
		t.Fatalf("path=%q want /v1/user-templates/tmpl_123/versions", requestedPath)
	}
	if payload["versionNote"] != "fix judge template" {
		t.Fatalf("versionNote=%q", payload["versionNote"])
	}
	spec, ok := payload["canonicalSpec"].(map[string]any)
	if !ok {
		t.Fatalf("canonicalSpec missing or wrong type: %#v", payload["canonicalSpec"])
	}
	if _, ok := spec["meta"]; !ok {
		t.Fatalf("canonicalSpec missing lowerCamel meta: %#v", spec)
	}
	if _, ok := spec["Meta"]; ok {
		t.Fatalf("canonicalSpec should not contain PascalCase Meta: %#v", spec)
	}
	if !strings.Contains(out.String(), `"templateId": "tmpl_123"`) {
		t.Fatalf("output missing template id: %s", out.String())
	}
	if !strings.Contains(out.String(), `"versionNumber": 2`) {
		t.Fatalf("output missing version number: %s", out.String())
	}
}

func TestTemplateSpecModelsCmdCanIncludeUnavailableModels(t *testing.T) {
	var requestedQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"models":[]}`))
	}))
	defer server.Close()

	opts := &rootOptions{
		server:  server.URL,
		timeout: time.Second,
		output:  "json",
	}
	cmd := newTemplateSpecModelsCmd(opts)
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"image-generate", "--include-unavailable"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("models command error = %v", err)
	}
	if !strings.Contains(requestedQuery, "onlyAvailable=false") {
		t.Fatalf("query %q missing onlyAvailable=false", requestedQuery)
	}
	if !strings.Contains(out.String(), `"models": []`) {
		t.Fatalf("json output missing models array: %s", out.String())
	}
}

func TestTemplateSpecSubmitWorkbookSendsFilename(t *testing.T) {
	workbookPath := filepath.Join(t.TempDir(), "custom-input.xlsx")
	if err := os.WriteFile(workbookPath, []byte("xlsx bytes"), 0o644); err != nil {
		t.Fatalf("write workbook: %v", err)
	}

	var submitFilename string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		switch r.URL.Path {
		case "/v1/batch/user-template-workbook:validate":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"valid":true}`))
		case "/v1/batch/user-template-workbook:submit":
			submitFilename, _ = payload["filename"].(string)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"runId":"run_123","status":"pending","acceptedAt":"1777699967"}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	opts := &rootOptions{
		server:  server.URL,
		timeout: time.Second,
		output:  "json",
	}
	cmd := newTemplateSpecSubmitWorkbookCmd(opts)
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"tmpl_123", "ver_123", workbookPath})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("submit-workbook command error = %v", err)
	}
	if submitFilename != "custom-input.xlsx" {
		t.Fatalf("filename=%q want custom-input.xlsx", submitFilename)
	}
	if !strings.Contains(out.String(), `"runId": "run_123"`) {
		t.Fatalf("output missing run id: %s", out.String())
	}
}
