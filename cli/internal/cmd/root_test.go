package cmd

import (
	"strings"
	"testing"
)

func TestRootRejectsUnsupportedOutputFormat(t *testing.T) {
	cmd := NewRootCmd()
	cmd.SetArgs([]string{"--output", "yaml", "doctor"})

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), `unsupported output format "yaml"`) {
		t.Fatalf("error=%v want unsupported output format", err)
	}
}

func TestRootReadsVerboseEnvironment(t *testing.T) {
	t.Setenv("LOOMLOOM_VERBOSE", "true")
	cmd := NewRootCmd()
	flag := cmd.PersistentFlags().Lookup("verbose")
	if flag == nil || flag.DefValue != "true" {
		t.Fatalf("verbose default=%v want true from environment", flag)
	}
}
